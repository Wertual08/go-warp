package storage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wertual08/go-warp/storage"
)

type DispatcherRepository struct {
    mtx         sync.Mutex
    dispatchers []dispatcherEntry
}

type dispatcherEntry struct {
    expiresAt  time.Time
    dispatcher storage.DispatcherDto
}

func (inst *DispatcherRepository) Create(
    dto storage.DispatcherDto, 
    ttl time.Duration,
    ctx context.Context,
) error {
    entry := dispatcherEntry {
        expiresAt: time.Now().Add(ttl),
        dispatcher: dto,
    }

    inst.mtx.Lock()
    defer inst.mtx.Unlock()

    inst.dispatchers = append(inst.dispatchers, entry)

    return nil
}

func (inst *DispatcherRepository) Remove(id uuid.UUID, ctx context.Context) error {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()

    for i, entry := range inst.dispatchers {
        if entry.dispatcher.Id == id {
            inst.dispatchers[i] = inst.dispatchers[len(inst.dispatchers) - 1]
            inst.dispatchers = inst.dispatchers[:len(inst.dispatchers) - 1]
        }
    }

    return nil
}

func (inst *DispatcherRepository) List(ctx context.Context) ([]storage.DispatcherDto, error) {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()
    
    result := make([]storage.DispatcherDto, len(inst.dispatchers))

    for i, entry := range inst.dispatchers {
        result[i] = entry.dispatcher
    }

    return result, nil
}
