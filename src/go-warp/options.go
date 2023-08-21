package warp

import "time"

type InstanceOptions struct {
	Enabled bool
    HeartbeatPeriod time.Duration
    FailDelay time.Duration
}

type QueueOptions struct {
    Name string
    Enabled bool
    ChannelsCount int32
    ChannelsOffset int32
    BatchSize int32
    MaxFails int32
    BatchDelay time.Duration
}
