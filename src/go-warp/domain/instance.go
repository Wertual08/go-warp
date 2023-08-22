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
    options *warp.InstanceOptions
    repositoryFactory warp.RepositoryFactory
    dispatcher controller.Dispatcher
	queues []controller.Queue
    ctx context.Context
}

func NewInstance(
    options *warp.InstanceOptions,
    repositoryFactory warp.RepositoryFactory,
    ctx context.Context,
) *Instance {
    inst := &Instance {
        options: options,
        repositoryFactory: repositoryFactory,
        ctx: ctx,
    };

    inst.dispatcher = controller.CreateDispatcher(
        repositoryFactory.Dispatcher,
        inst.runQueues,
    )

    return inst
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

func (inst *Instance) runQueues(
    running *atomic.Int32, 
    wg *sync.WaitGroup,
    stride int32,
    offset int32,
) {
    wg.Add(len(inst.queues))

    for i := range inst.queues {
        go func() {
            defer wg.Done()

            for running.Load() != 0 {
                err := inst.queues[i].Handle(
                    stride,
                    offset,
                    inst.repositoryFactory.Objective,
                    inst.ctx,
                )

                if err != nil {
                    // TODO: Log error
                    time.Sleep(inst.options.FailDelay)
                }
            }
        }()
    }
}
