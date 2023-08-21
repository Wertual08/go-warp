package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DispatcherRepository interface {
    Create(
        dtos DispatcherDto, 
        ttl time.Duration,
        ctx context.Context,
    ) error

    Remove(id uuid.UUID, ctx context.Context) error

    List(ctx context.Context) ([]DispatcherDto, error)
}
