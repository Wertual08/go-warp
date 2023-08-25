package storage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wertual08/go-warp/storage"
)

type ObjectiveRepository struct {
    mtx sync.RWMutex
    // TODO: Replace with priority queue mb...
    objectives map[int32][][]storage.ObjectiveDto
}

func (inst *ObjectiveRepository) Create(
    dtos []storage.ObjectiveDto, 
    ctx  context.Context,
) error {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()

    for _, objective := range dtos {
        var sections [][]storage.ObjectiveDto

        if currentSections, ok := inst.objectives[objective.QueueId]; ok {
            appended := false
            for len(currentSections) <= int(objective.Section) {
                appended = true
                currentSections = append(currentSections , make([]storage.ObjectiveDto, 0))
            }

            if appended {
                inst.objectives[objective.QueueId] = currentSections 
            }

            sections = currentSections
        } else {
            sections = make([][]storage.ObjectiveDto, objective.Section + 1)
            inst.objectives[objective.QueueId] = sections
        }
        
        objective.Id = uuid.New()
        sections[objective.Section] = append(sections[objective.Section], objective)
    }

    return nil
}

func (inst *ObjectiveRepository) Remove(
    dtos []storage.ObjectiveDto, 
    ctx  context.Context,
) error {
    inst.mtx.Lock()
    defer inst.mtx.Unlock()

    for _, objective := range dtos {
        sections := inst.objectives[objective.QueueId]
        section := sections[objective.Section]

        for i, currentObjective := range section {
            if currentObjective.Id == objective.Id {
                section[i] = section[len(section) - 1]
                sections[objective.Section] = section[:len(section) - 1]
                break
            }
        }
    }

    return nil
}

func (inst *ObjectiveRepository) List(
    queueId int32, 
    channel int32, 
    limit   int32,
    now     time.Time,
    ctx     context.Context,
) ([]storage.ObjectiveDto, error) {
    inst.mtx.RLock()
    defer inst.mtx.RUnlock()

    result := make([]storage.ObjectiveDto, 0, limit)

    if sections, ok := inst.objectives[queueId]; ok {
        if int(channel) < len(sections) {
            for _, objective := range sections[channel] {
                if len(result) >= int(limit) {
                    break
                }

                if objective.ScheduledAt.Before(now) {
                    result = append(result, objective)
                }
            }
        }
    }

    return result, nil
}
