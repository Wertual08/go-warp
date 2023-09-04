package warp

import (
	"context"
	"time"
)

type HandlerManager interface {
    Fail(index int, err error)
    FailAt(index int, err error, scheduledAt time.Time)
}

type Handler[T Objective[T]] func (
    manager    HandlerManager, 
    objectives []T,
    ctx        context.Context,
) error
