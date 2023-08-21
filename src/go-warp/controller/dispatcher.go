package controller

import (
	"bytes"
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/storage"
)

type Dispatcher struct {
    id uuid.UUID
    dispatcherRepository storage.DispatcherRepository

    Stride int32
    Offset int32

    queuesWaitGroup sync.WaitGroup
    queuesRunning atomic.Int32
	queues []Queue
}

func CreateDispatcher(
    options *warp.InstanceOptions,
    dispatcherRepository storage.DispatcherRepository,
) Dispatcher {
    return Dispatcher{
        id: uuid.New(),
        dispatcherRepository: dispatcherRepository,
        Stride: -1,
        Offset: -1,
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

    if requiredStride != inst.Stride || requiredOffset != inst.Offset {
        // TODO: Maybe i should update lifetime while waiting...
        inst.queuesRunning.Store(0)
        inst.queuesWaitGroup.Wait()

        inst.Stride = requiredStride
        inst.Offset = requiredOffset
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
        inst.runQueues()
    }

    return nil
}

func (inst *Dispatcher) Finish(ctx context.Context) error {
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
        Stride: inst.Stride,
        Offset: inst.Offset,
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

func (inst *Dispatcher) runQueues() {
    inst.queuesWaitGroup.Add(len(inst.queues))

    for _ = range inst.queues {
        go func() {
            defer inst.queuesWaitGroup.Done()

            for inst.queuesRunning.Load() != 0 {
                //if err := inst.handle(&inst.queues[i], inst.ctx); err != nil {
                //    // TODO: Log error
                //    time.Sleep(inst.options.FailDelay)
                //}
            }
        }()
    }
}