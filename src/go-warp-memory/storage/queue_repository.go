package storage

import (
	"context"
	"sync"
)

type QueueRepository struct {
    mtx sync.Mutex
    queues map[string]int32
}

func (inst *QueueRepository) FindOrCreate(name string, ctx context.Context) (int32, error) {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()

    if id, ok := inst.queues[name]; ok {
        return id, nil
    }

    id := int32(len(inst.queues) + 1)

    inst.queues[name] = id

    return id, nil
}
