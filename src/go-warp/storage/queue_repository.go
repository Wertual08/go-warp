package storage

import "context"

type QueueRepository interface {
	FindOrCreate(name string, ctx context.Context) (int32, error)
}
