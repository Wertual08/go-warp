package warpMemory

import (
	"context"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/memory/storage"
	warpStorage "github.com/wertual08/go-warp/storage"
)

func NewInstance(
    options *warp.InstanceOptions,
    onError func (error),
    ctx     context.Context,
) *warp.Instance {
    repositoryFactory := warp.RepositoryFactory{
        Queue: &storage.QueueRepository{
            Queues: make(map[string]int32),
        },
        Dispatcher: &storage.DispatcherRepository{},
        Objective: &storage.ObjectiveRepository{
            Objectives: make(map[int32][][]warpStorage.ObjectiveDto),
        },
    }

    return warp.NewInstance(
        options,
        repositoryFactory,
        onError,
        ctx,
    )
}
