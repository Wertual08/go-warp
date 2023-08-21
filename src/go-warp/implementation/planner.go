package implementation

import (
	"context"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/storage"
)

type planner[T warp.Objective] struct {
    QueueId int32
    Options *warp.QueueOptions
    ObjectiveRepository storage.ObjectiveRepository
}

func (inst *planner[T]) Plan(objective T, ctx context.Context) error {
    objectives := make([]T, 1)

    objectives[0] = objective

    return inst.PlanBatch(objectives, ctx)
}

func (inst *planner[T]) PlanBatch(objectives []T, ctx context.Context) error {
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

func CreatePlanner[T warp.Objective](
    queueId int32,
    options *warp.QueueOptions,
    objectiveRepository storage.ObjectiveRepository,
) warp.Planner[T] {
    return &planner[T]{
        QueueId: queueId,
        Options: options,
        ObjectiveRepository: objectiveRepository,
    }
}
