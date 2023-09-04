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
    Dispatchers []dispatcherEntry
}

type dispatcherEntry struct {
    expiresAt  time.Time
    dispatcher storage.DispatcherDto
}

func (inst *DispatcherRepository) Upsert(
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
    
    found := false

    for i, dispatcher := range inst.Dispatchers {
        if dispatcher.dispatcher.Id == dto.Id {
            inst.Dispatchers[i] = entry
            found = true
            break;
        }
    }

    if !found {
        inst.Dispatchers = append(inst.Dispatchers, entry)
    }

    return nil
}

func (inst *DispatcherRepository) Remove(id uuid.UUID, ctx context.Context) error {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()

    for i, entry := range inst.Dispatchers {
        if entry.dispatcher.Id == id {
            inst.Dispatchers[i] = inst.Dispatchers[len(inst.Dispatchers) - 1]
            inst.Dispatchers = inst.Dispatchers[:len(inst.Dispatchers) - 1]
        }
    }

    return nil
}

func (inst *DispatcherRepository) List(ctx context.Context) ([]storage.DispatcherDto, error) {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()
    
    result := make([]storage.DispatcherDto, len(inst.Dispatchers))

    for i, entry := range inst.Dispatchers {
        result[i] = entry.dispatcher
    }

    return result, nil
}
