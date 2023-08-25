package storage

import (
	"context"
	"time"
)

type ObjectiveRepository interface {
    Create(dtos []ObjectiveDto, ctx context.Context) error
    Remove(dtos []ObjectiveDto, ctx context.Context) error
    List(
        queueId int32, 
        channel int32, 
        limit   int32,
        now     time.Time,
        ctx     context.Context,
    ) ([]ObjectiveDto, error)
}
