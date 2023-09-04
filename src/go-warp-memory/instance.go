package warp_memory

import (
	"context"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/implementation"
	warp_storage "github.com/wertual08/go-warp/storage"
	"github.com/wertual08/go-warp/memory/storage"
)

func NewInstance(
    options *warp.InstanceOptions,
    ctx     context.Context,
) *implementation.Instance {
    repositoryFactory := warp.RepositoryFactory{
        Queue: &storage.QueueRepository{
            Queues: make(map[string]int32),
        },
        Dispatcher: &storage.DispatcherRepository{},
        Objective: &storage.ObjectiveRepository{
            Objectives: make(map[int32][][]warp_storage.ObjectiveDto),
        },
    }

    return implementation.NewInstance(
        options,
        repositoryFactory,
        ctx,
    )
}
