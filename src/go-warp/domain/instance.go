package domain

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/controller"
	"github.com/wertual08/go-warp/implementation"
)

type Instance struct {
    ctx context.Context
    options *warp.InstanceOptions
    repositoryFactory warp.RepositoryFactory
    dispatcher controller.Dispatcher
}

func Register[T warp.Objective](
    instance *Instance,
    options *warp.QueueOptions,
    handler warp.Handler[T],
) (warp.Planner[T], error) {
    queueId, err := instance.repositoryFactory.Queue.FindOrCreate(
        options.Name, 
        instance.ctx,
    )
    if err != nil {
        return nil, err
    }


    queueController := controller.CreateQueue[T](
        queueId,
        options,
        handler,
    )

    instance.queues = append(instance.queues, queueController)

    planner := implementation.CreatePlanner[T](
        queueId,
        options,
        instance.repositoryFactory.Objective,
    )

    return planner, nil
}

func (inst *Instance) Run() {
    running := true
    for running {
        select {
        case <- inst.ctx.Done():
            running = false
            continue
        default:
        }

        err := inst.dispatcher.Process(inst.options.HeartbeatPeriod, inst.ctx)

        if err != nil {
            // TODO: Log error
            time.Sleep(inst.options.FailDelay)
        } else {
            time.Sleep(inst.options.HeartbeatPeriod / 2)
        }
    }

    inst.dispatcher.Finish(inst.ctx)
}

func (inst *Instance) handle(
    queue *controller.Queue,
    ctx context.Context,
) error {
    if !queue.Options.Enabled {
        time.Sleep(queue.Options.BatchDelay)
        return nil
    }

    now := time.Now()

    for i, channel := range queue.Channels {
        for channel >= int32(len(queue.DispatcherOffsets)) {
            queue.DispatcherOffsets = append(queue.DispatcherOffsets, -1)
        }

        queue.DispatcherOffsets[channel] = i

        if i >= len(queue.Dispatchers) {
            queue.Dispatchers = append(queue.Dispatchers, queue.DispatcherFactory())
        }
    }

    objectives, err := inst.repositoryFactory.Objective.List(
        queue.Id,
        queue.Channels,
        queue.Options.BatchSize,
        now,
        ctx,
    )
    if err != nil {
        return err
    }

    for index := range objectives {
        objective := &objectives[index]
        dispatcher := queue.DispatcherOffsets[objective.Channel]

        if err := queue.Dispatchers[dispatcher].Push(objective); err != nil {
            return err
        }
    }

    batchDelay := false

    wg := sync.WaitGroup{}
    for _, channel := range queue.Channels {
        wg.Add(1)
        go func() {
            defer wg.Done()

            queue.Dispatchers[channel].Handle(ctx)
        }()
    }
    wg.Done()
    
    return nil
}
