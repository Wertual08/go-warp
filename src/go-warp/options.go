package warp

import "time"

type InstanceOptions struct {
	Enabled bool
	Lifetime time.Duration 
    HeartbeatPeriod time.Duration
    OccupationPeriod time.Duration
    CountChannelsPeriod time.Duration
    IdleDelay time.Duration
    FailDelay time.Duration
}

type QueueOptions struct {
    Name string
    Enabled bool
    ChannelsCount int32
    BatchSize int32
    MaxFails int32
    BatchDelay time.Duration
}
