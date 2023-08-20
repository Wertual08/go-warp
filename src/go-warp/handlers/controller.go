package handlers

import (
	"context"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/storage"
)

type Controller interface {
	Handle(
        objectives []storage.ObjectiveDto,
        ctx context.Context,
    ) ([]fail, error)
}

type controllerImpl[T warp.Objective] struct {
	objectives []T
	handler    warp.Handler[T]
	manager    managerImpl
}

func CreateController[T warp.Objective](handler warp.Handler[T]) Controller {
    return &controllerImpl[T] {
        handler: handler,
    }
}

func (inst *controllerImpl[T]) Handle(
	objectives []storage.ObjectiveDto,
    ctx context.Context,
) ([]fail, error) {
	inst.objectives = inst.objectives[:len(objectives)]
    inst.manager.sourceObjectives = objectives
	inst.manager.fails = inst.manager.fails[:0]

	for index, objective := range inst.objectives {
		content := objectives[index].Content

		if err := objective.Deserialize(content); err != nil {
			return nil, err
		}
	}

	if err := inst.handler.Handle(&inst.manager, inst.objectives, ctx); err != nil {
		return nil, err
	}

	return inst.manager.fails, nil
}
