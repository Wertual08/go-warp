package storage

import "context"

type ObjectiveRepository interface {
    Create(dtos []ObjectiveDto, ctx context.Context) error
}
