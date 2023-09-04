package warp

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const idleThreshold = 3

type Instance struct {
    options           *InstanceOptions
    repositoryFactory RepositoryFactory
    onError           func (error)
    dispatcher        dispatcher
	queues            []queue
    ctx               context.Context
}

func NewInstance(
    options           *InstanceOptions,
    repositoryFactory RepositoryFactory,
    onError           func (error),
    ctx               context.Context,
) *Instance {
    inst := &Instance {
        options: options,
        repositoryFactory: repositoryFactory,
        onError: onError,
        ctx: ctx,
    };

    inst.dispatcher = newDispatcher(
        repositoryFactory.Dispatcher,
        inst.runQueues,
    )

    return inst
}

func Register[T Objective[T]](
    instance    *Instance,
    options     *QueueOptions,
    handler     Handler[T],
) (Planner[T], error) {
    queueId, err := instance.repositoryFactory.Queue.FindOrCreate(
        options.Name, 
        instance.ctx,
    )
    if err != nil {
        return nil, err
    }

    queueController := newQueue(
        queueId,
        options,
        newHandlerFactory[T](handler),
    )

    instance.queues = append(instance.queues, queueController)

    planner := newPlanner[T](
        queueId,
        options,
        instance.repositoryFactory.Objective,
    )

    return planner, nil
}

func (inst *Instance) Run() {
    running := true
    active := false
    idleCounter := 0

    for running {
        select {
        case <- inst.ctx.Done():
            running = false
            continue
        default:
        }

        if inst.options.Enabled {
            active = true

            idle, err := inst.dispatcher.process(inst.options.IdleHeartbeatPeriod * 2, inst.ctx)

            if err != nil {
                inst.onError(err)
                time.Sleep(inst.options.FailDelay)
            } else {
                if !idle {
                    idleCounter = -1
                }

                if idleCounter < idleThreshold {
                    idleCounter += 1
                }

                if idleCounter >= idleThreshold {
                    time.Sleep(inst.options.IdleHeartbeatPeriod)
                } else {
                    time.Sleep(inst.options.ActiveHeartbeatPeriod)
                }
            }
        } else {
            if active {
                err := inst.dispatcher.finish(inst.ctx)

                if err != nil {
                    inst.onError(err)
                    time.Sleep(inst.options.FailDelay)
                } else {
                    active = false
                    time.Sleep(inst.options.ActiveHeartbeatPeriod)
                }
            } else {
                time.Sleep(inst.options.ActiveHeartbeatPeriod)
            }
        }
    }

    inst.dispatcher.finish(inst.ctx)
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
                err := inst.queues[i].handle(
                    stride,
                    offset,
                    inst.repositoryFactory.Objective,
                    inst.ctx,
                )

                if err != nil {
                    inst.onError(err)
                    time.Sleep(inst.options.FailDelay)
                }
            }
        }()
    }
}
