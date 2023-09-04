package controller

import (
	"context"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/storage"
)

type Handler interface {
	Handle(objectives []storage.ObjectiveDto, ctx context.Context) error
    Succeded() int
    Failed() []storage.ObjectiveDto
}

type handlerImpl[T warp.Objective[T]] struct {
	sourceObjectives []storage.ObjectiveDto
    failedObjectives []storage.ObjectiveDto
	objectives       []T
	handler          warp.Handler[T]
}

func (inst *handlerImpl[T]) Handle(
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
            // TODO: Maybe check if already failed???
            objective.FailCount += 1
            objective.FailReason = reason

            inst.failedObjectives = append(inst.failedObjectives, objective)
        }
    }

    return nil
}

func (inst *handlerImpl[T]) Succeded() int {
    return len(inst.objectives) - len(inst.failedObjectives)
}

func (inst *handlerImpl[T]) Failed() []storage.ObjectiveDto {
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


func NewHandlerFactory[T warp.Objective[T]](
    handler warp.Handler[T],
) func() Handler {
    return func() Handler {
        return &handlerImpl[T]{
            handler: handler,
        }
    }
}
