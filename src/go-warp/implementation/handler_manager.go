package implementation

import (
	"time"

	"github.com/wertual08/go-warp/storage"
)

type HandlerManager struct {
    FailCount int
	SourceObjectives []*storage.ObjectiveDto
}

func (inst *HandlerManager) Fail(index int, err error) {
    inst.FailCount += 1

    dto := inst.SourceObjectives[index]
    dto.FailCount += 1
    dto.FailReason = err.Error()
}
    
func (inst *HandlerManager) FailAt(index int, err error, scheduledAt time.Time) {
    inst.FailCount += 1

    dto := inst.SourceObjectives[index]
    dto.FailCount += 1
    dto.FailReason = err.Error()
    dto.ScheduledAt = scheduledAt
}

