package warp

import (
	"context"
	"time"
)

type HandlerManager interface {
    Fail(index int, err error)
    FailDelay(index int, err error, retryDelay time.Duration)
    Metadata(index int) map[string]string
}

type Handler[T Objective] interface {
	Handle(
        manager HandlerManager, 
        objectives []T,
        ctx context.Context,
    ) error
}

