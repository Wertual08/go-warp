package controller

import (
	"context"

	"github.com/wertual08/go-warp"
	implementation "github.com/wertual08/go-warp/implementation"
	"github.com/wertual08/go-warp/storage"
)

type Handler interface {
    Reset()
    Push(dto *storage.ObjectiveDto) error
	Handle(ctx context.Context)
    Succeded() int
}

type handlerImpl[T warp.Objective] struct {
    sourceObjectives []*storage.ObjectiveDto
	objectives []T
	handler    warp.Handler[T]
	manager    implementation.HandlerManager
}

func (inst *handlerImpl[T]) Reset() {
	inst.sourceObjectives = inst.sourceObjectives[:0]
	inst.objectives = inst.objectives[:0]
}

func (inst *handlerImpl[T]) Push(dto *storage.ObjectiveDto) error {
    var deserialized T
    
    if err := deserialized.Deserialize(dto.Content); err != nil {
        return err
    }
    
    inst.sourceObjectives = append(inst.sourceObjectives, dto)
    inst.objectives = append(inst.objectives, deserialized)

    return nil
}

func (inst *handlerImpl[T]) Handle(ctx context.Context) {
    err := inst.handler.Handle(&inst.manager, inst.objectives, ctx)

    if err != nil {
        reason := err.Error()

        for _, dto := range inst.sourceObjectives {
            dto.FailCount += 1
            dto.FailReason = reason
        }
    }
}

func (inst *handlerImpl[T]) Succeded() int {
    return len(inst.objectives) - inst.manager.FailCount
}

func createHandlerFactory[T warp.Objective](handler warp.Handler[T]) func() Handler {
    return func() Handler{
        return &handlerImpl[T]{
            handler: handler,
        }
    }
}
