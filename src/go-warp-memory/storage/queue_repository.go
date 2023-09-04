package storage

import (
	"context"
	"sync"
)

type QueueRepository struct {
    mtx sync.Mutex
    Queues map[string]int32
}

func (inst *QueueRepository) FindOrCreate(name string, ctx context.Context) (int32, error) {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()

    if id, ok := inst.Queues[name]; ok {
        return id, nil
    }

    id := int32(len(inst.Queues) + 1)
    inst.Queues[name] = id

    return id, nil
}
