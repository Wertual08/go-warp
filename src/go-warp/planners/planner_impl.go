package planners

import (
	"context"
	"encoding/json"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/storage"
)

type PlannerImpl[T warp.Objective] struct {
    QueueId int32
    Options *warp.QueueOptions
    ObjectiveRepository storage.ObjectiveRepository
}

func (inst *PlannerImpl[T]) Plan(objective T, ctx context.Context) error {
    objectives := make([]T, 1)

    objectives[0] = objective

    return inst.PlanBatch(objectives, ctx)
}

func (inst *PlannerImpl[T]) PlanBatch(objectives []T, ctx context.Context) error {
    dtos := make([]storage.ObjectiveDto, len(objectives))
    now := time.Now()

    for index, objective := range objectives {
        dto := &dtos[index]
        dto.QueueId = inst.QueueId
        dto.Channel = objective.HashCode() % inst.Options.ChannelsCount
        dto.ScheduledAt = now
        dto.CreatedAt = now 
        dto.Content = objective.Serialize()
    }

    return inst.ObjectiveRepository.Create(dtos, ctx)
}

func (inst *PlannerImpl[T]) PlanParametrized(objective warp.ParametrizedObjective[T], ctx context.Context) error {
    objectives := make([]warp.ParametrizedObjective[T], 1)

    objectives[0] = objective

    return inst.PlanParametrizedBatch(objectives, ctx)
}

func (inst *PlannerImpl[T]) PlanParametrizedBatch(objectives []warp.ParametrizedObjective[T], ctx context.Context) error {
    dtos := make([]storage.ObjectiveDto, len(objectives))
    now := time.Now()

    for index, objective := range objectives {
        dto := &dtos[index]
        dto.QueueId = inst.QueueId
        dto.Channel = objective.Value.HashCode() % inst.Options.ChannelsCount
        dto.ScheduledAt = objective.ScheduledAt
        dto.CreatedAt = now 
        dto.Content = objective.Value.Serialize()
        
        metadata, err := json.Marshal(objective.Metadata)
        if err != nil {
            return err
        }

        dto.Metadata = string(metadata)
    }

    return inst.ObjectiveRepository.Create(dtos, ctx)
}
