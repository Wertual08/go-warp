package warp

import "time"

type InstanceOptions struct {
	Enabled bool
    ActiveHeartbeatPeriod time.Duration
    IdleHeartbeatPeriod time.Duration
    FailDelay time.Duration
}

type QueueOptions struct {
    Name string
    Enabled bool
    SectionsCount int32
    SectionsOffset int32
    BatchSize int32
    MaxFails int32
    BatchDelay time.Duration
}
