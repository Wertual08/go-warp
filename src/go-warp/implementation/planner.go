package implementation

import (
	"context"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/storage"
)

type planner[T warp.Objective[T]] struct {
    queueId             int32
    options             *warp.QueueOptions
    objectiveRepository storage.ObjectiveRepository
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

        serialized, err := objective.Serialize()
        if err != nil {
            return err
        }

        dto.QueueId = inst.queueId
        dto.Section = (objective.HashCode() & 0x7fffffff) % inst.options.SectionsCount
        dto.ScheduledAt = now
        dto.CreatedAt = now 
        dto.Content = serialized
    }

    return inst.objectiveRepository.Create(dtos, ctx)
}

func NewPlanner[T warp.Objective[T]](
    queueId             int32,
    options             *warp.QueueOptions,
    objectiveRepository storage.ObjectiveRepository,
) warp.Planner[T] {
    return &planner[T]{
        queueId: queueId,
        options: options,
        objectiveRepository: objectiveRepository,
    }
}
