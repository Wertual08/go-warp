package controller

import (
	"bytes"
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/wertual08/go-warp/storage"
)

type Dispatcher struct {
    id uuid.UUID
    dispatcherRepository storage.DispatcherRepository

    stride int32
    offset int32

    queuesWaitGroup sync.WaitGroup
    queuesRunning atomic.Int32
    queuesCallback func (*atomic.Int32, *sync.WaitGroup, int32, int32)
}

func CreateDispatcher(
    dispatcherRepository storage.DispatcherRepository,
    queuesCallback func (*atomic.Int32, *sync.WaitGroup, int32, int32),
) Dispatcher {
    return Dispatcher{
        id: uuid.New(),
        dispatcherRepository: dispatcherRepository,
        stride: -1,
        offset: -1,
        queuesCallback: queuesCallback,
    }
}

func (inst *Dispatcher) Process(
    lifetime time.Duration,
    ctx context.Context,
) error {
    requiredStride, requiredOffset, err := inst.findRequred(ctx)
    if err != nil {
        return err
    }

    if requiredStride != inst.stride || requiredOffset != inst.offset {
        if inst.queuesRunning.Load() != 0 {
            // TODO: Maybe i should update lifetime while waiting...
            inst.queuesRunning.Store(0)
            inst.queuesWaitGroup.Wait()
        }

        inst.stride = requiredStride
        inst.offset = requiredOffset
    }

    if err := inst.upsert(lifetime, ctx); err != nil {
        return err
    }

    allValid, err := inst.checkValid(ctx)
    if err != nil {
        return err
    }

    if allValid && inst.queuesRunning.Load() == 0 {
        inst.queuesRunning.Store(1)
        inst.queuesCallback(
            &inst.queuesRunning,
            &inst.queuesWaitGroup,
            inst.stride,
            inst.offset,
        )
    }

    return nil
}

func (inst *Dispatcher) Finish(ctx context.Context) error {
    if inst.queuesRunning.Load() != 0 {
        inst.queuesRunning.Store(0)
        inst.queuesWaitGroup.Wait()
    }

    return inst.dispatcherRepository.Remove(inst.id, ctx)
}

func (inst *Dispatcher) findRequred(ctx context.Context) (int32, int32, error) {
    dispatchers, err := inst.dispatcherRepository.List(ctx)
    if err != nil {
        return 0, 0, err
    }

    stride := int32(len(dispatchers))
    offset := int32(0)

    for _, dispatcher := range dispatchers {
        if bytes.Compare(dispatcher.Id[:], inst.id[:]) < 0 {
            offset += 1
        }
    }

    return stride, offset, nil
}

func (inst *Dispatcher) upsert(
    lifetime time.Duration, 
    ctx context.Context,
) error {
    dto := storage.DispatcherDto{
        Id: inst.id,
        Stride: inst.stride,
        Offset: inst.offset,
    }

    return inst.dispatcherRepository.Create(dto, lifetime, ctx)
}

func (inst *Dispatcher) checkValid(ctx context.Context) (bool, error) {
    dispatchers, err := inst.dispatcherRepository.List(ctx)
    if err != nil {
        return false, err
    }

    sort.Slice(
        dispatchers,
        func (lhs int, rhs int) bool {
            return bytes.Compare(
                dispatchers[lhs].Id[:], 
                dispatchers[rhs].Id[:],
            ) < 0
        },
    )
    
    stride := int32(len(dispatchers))
    allValid := true
    found := false

    for i, dispatcher := range dispatchers {
        if dispatcher.Stride != stride || dispatcher.Offset != int32(i) {
            allValid = false
            break
        }

        if dispatcher.Id == inst.id {
            found = true
        }
    }
    
    return allValid && found, nil
}
