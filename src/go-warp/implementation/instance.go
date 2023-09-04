package implementation

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/controller"
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

    inst.dispatcher = controller.NewDispatcher(
        repositoryFactory.Dispatcher,
        inst.runQueues,
    )

    return inst
}

func Register[T warp.Objective[T]](
    instance    *Instance,
    options     *warp.QueueOptions,
    handler     warp.Handler[T],
) (warp.Planner[T], error) {
    queueId, err := instance.repositoryFactory.Queue.FindOrCreate(
        options.Name, 
        instance.ctx,
    )
    if err != nil {
        return nil, err
    }

    queueController := controller.NewQueue(
        queueId,
        options,
        controller.NewHandlerFactory[T](handler),
    )

    instance.queues = append(instance.queues, queueController)

    planner := NewPlanner[T](
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
        
        err := inst.dispatcher.Process(inst.options.IdleHeartbeatPeriod, inst.ctx)

        if err != nil {
            // TODO: Log error
            time.Sleep(inst.options.FailDelay)
        } else {
            // TODO: Calculate period from idle and active
            time.Sleep(inst.options.IdleHeartbeatPeriod / 2)
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
