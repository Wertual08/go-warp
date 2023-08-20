package warp

import "time"

type ParametrizedObjective[T Objective] struct {
    Value T
    Metadata map[string]string
    ScheduledAt time.Time
}
