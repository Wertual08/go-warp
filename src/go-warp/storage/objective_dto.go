package storage

import (
	"time"

	"github.com/google/uuid"
)

type ObjectiveDto struct {
	QueueId     int32
	Channel     int32
	Id          uuid.UUID
	ScheduledAt time.Time
	CreatedAt   time.Time
    FailCount   int32
    FailReason  string
	Content     []byte
	Metadata    string
}
