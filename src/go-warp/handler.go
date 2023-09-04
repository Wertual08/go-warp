package warp

import (
	"context"
	"time"

	"github.com/wertual08/go-warp/storage"
)

type HandlerManager interface {
    Fail(index int, err error)
    FailAt(index int, err error, scheduledAt time.Time)
}

type Handler[T Objective[T]] func (
    manager    HandlerManager, 
    objectives []T,
    ctx        context.Context,
) error

type handler interface {
	handle(objectives []storage.ObjectiveDto, ctx context.Context) error
    succeded() int
    failed() []storage.ObjectiveDto
}

type handlerImpl[T Objective[T]] struct {
	sourceObjectives []storage.ObjectiveDto
    failedObjectives []storage.ObjectiveDto
	objectives       []T
	handler          Handler[T]
}

func (inst *handlerImpl[T]) handle(
    objectives []storage.ObjectiveDto,
    ctx        context.Context,
) error {
    inst.sourceObjectives = objectives
    inst.failedObjectives = inst.failedObjectives[:0]
	inst.objectives = inst.objectives[:0]
    for _, objective := range objectives {
        var deserialized T

        deserialized, err := deserialized.Deserialize(objective.Content)

        if err != nil {
            return err
        }
        
        inst.objectives = append(inst.objectives, deserialized)
    }

    err := inst.handler(inst, inst.objectives, ctx)

    if err != nil {
        reason := err.Error()
        inst.failedObjectives = inst.failedObjectives[:0]

        for _, objective := range objectives {
            objective.FailCount += 1
            objective.FailReason = reason

            inst.failedObjectives = append(inst.failedObjectives, objective)
        }
    }

    return nil
}

func (inst *handlerImpl[T]) succeded() int {
    return len(inst.objectives) - len(inst.failedObjectives)
}

func (inst *handlerImpl[T]) failed() []storage.ObjectiveDto {
    return inst.failedObjectives
}


func (inst *handlerImpl[T]) Fail(index int, err error) {
    objective := inst.sourceObjectives[index]

    objective.FailCount += 1 
    objective.FailReason = err.Error()
    inst.failedObjectives = append(inst.failedObjectives, objective)
}
    
func (inst *handlerImpl[T]) FailAt(index int, err error, scheduledAt time.Time) {
    objective := inst.sourceObjectives[index]

    objective.FailCount += 1
    objective.FailReason = err.Error()
    objective.ScheduledAt = scheduledAt

    inst.failedObjectives = append(inst.failedObjectives, objective)
}


func newHandlerFactory[T Objective[T]](
    h Handler[T],
) func() handler {
    return func() handler {
        return &handlerImpl[T]{
            handler: h,
        }
    }
}
