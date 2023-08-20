package instances

import (
	"context"
	"sync"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/handlers"
	"github.com/wertual08/go-warp/planners"
)

type Instance struct {
    ctx context.Context
    options *warp.InstanceOptions
    repositoryFactory warp.RepositoryFactory
	handlers map[int32]handlers.Controller
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

    controller := handlers.CreateController(handler)

    instance.handlers[queueId] = controller

    planner := &planners.PlannerImpl[T] {
        QueueId: queueId,
        Options: options,
        ObjectiveRepository: instance.repositoryFactory.Objective,
    }

    return planner, nil
}

func (inst *Instance) Run() error {
    wg := sync.WaitGroup{}
    wg.Add(len(inst.handlers) + 1)

    go func() {
        defer wg.Done()

        running := true
        for running {
            select {
            case <- inst.ctx.Done():
                running = false
                continue
            default:
                break
            }

            if err := inst.hearbeat(inst.ctx); err != nil {
                // TODO: Log error
                time.Sleep(inst.options.FailDelay)
            } else {
                time.Sleep(inst.options.IdleDelay)
            }
        }
    }()

    for queueId, handler := range inst.handlers {
        go func() {
            defer wg.Done()

            running := true
            for running {
                select {
                case <- inst.ctx.Done():
                    running = false
                    continue
                default:
                    break
                }

                if err := inst.handle(queueId, handler, inst.ctx); err != nil {
                    // TODO: Log error
                    time.Sleep(inst.options.FailDelay)
                } else {
                    time.Sleep(inst.options.IdleDelay)
                }
            }
        }()
    }

    wg.Done()

    return nil
}

func (inst *Instance) hearbeat(ctx context.Context) error {
}

func (inst *Instance) handle(
    queueId int32, 
    handler handlers.Controller,
    ctx context.Context,
) error {
}
