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

type handlerImpl[T warp.Objective] struct {
	SourceObjectives []storage.ObjectiveDto
    FailedObjectives []storage.ObjectiveDto
	objectives []T
	handler    warp.Handler[T]
}

func (inst *handlerImpl[T]) Handle(
    objectives []storage.ObjectiveDto,
    ctx context.Context,
) error {
    inst.SourceObjectives = objectives
    inst.FailedObjectives = inst.FailedObjectives[:0]
	inst.objectives = inst.objectives[:0]
    for _, objective := range objectives {
        var deserialized T
        
        if err := deserialized.Deserialize(objective.Content); err != nil {
            return err
        }
        
        inst.objectives = append(inst.objectives, deserialized)
    }

    err := inst.handler.Handle(inst, inst.objectives, ctx)

    if err != nil {
        reason := err.Error()
        inst.FailedObjectives = inst.FailedObjectives[:0]

        for _, objective := range objectives {
            objective.FailCount += 1
            objective.FailReason = reason

            inst.FailedObjectives = append(inst.FailedObjectives, objective)
        }
    }

    return nil
}

func (inst *handlerImpl[T]) Succeded() int {
    return len(inst.objectives) - len(inst.FailedObjectives)
}

func (inst *handlerImpl[T]) Failed() []storage.ObjectiveDto {
    return inst.FailedObjectives
}


func (inst *handlerImpl[T]) Fail(index int, err error) {
    objective := inst.SourceObjectives[index]

    objective.FailCount += 1
    objective.FailReason = err.Error()

    inst.FailedObjectives = append(inst.FailedObjectives, objective)
}
    
func (inst *handlerImpl[T]) FailAt(index int, err error, scheduledAt time.Time) {
    objective := inst.SourceObjectives[index]

    objective.FailCount += 1
    objective.FailReason = err.Error()
    objective.ScheduledAt = scheduledAt

    inst.FailedObjectives = append(inst.FailedObjectives, objective)
}


func createHandlerFactory[T warp.Objective](handler warp.Handler[T]) func() Handler {
    return func() Handler {
        return &handlerImpl[T]{
            handler: handler,
        }
    }
}
