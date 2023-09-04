package warp

import "time"

type InstanceOptions struct {
	Enabled bool
    ActiveHeartbeatPeriod time.Duration
    IdleHeartbeatPeriod time.Duration
    FailDelay time.Duration
}

type QueueOptions struct {
    Name           string
    SectionsCount  int32
    SectionsOffset int32
    BatchSize      int32
    BatchDelay     time.Duration
    MaxFails       int32
}
