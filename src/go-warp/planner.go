package warp

import "context"

type Planner[T Objective] interface {
	Plan(objective T, ctx context.Context) error
	PlanBatch(objectives []T, ctx context.Context) error
}
