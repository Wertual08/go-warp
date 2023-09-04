package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/implementation"
	warpMemory "github.com/wertual08/go-warp/memory"
)

type Objective struct {
    Id   int64
    Name string
}

func (inst Objective) Serialize() ([]byte, error) {
    return json.Marshal(&inst)
}

func (Objective) Deserialize(content []byte) (Objective, error) {
    var deserialized Objective
    return deserialized, json.Unmarshal(content, &deserialized);
}

func (inst Objective) HashCode() int32 {
    return int32(((inst.Id >> 32) & 0xffffffff) ^ (inst.Id & 0xffffffff))
}

func main() {
    COUNT := 1000000
    bus := make(chan Objective)

    instanceOptions := &warp.InstanceOptions{
        Enabled: true,
        ActiveHeartbeatPeriod: time.Second,
        IdleHeartbeatPeriod: time.Second,
        FailDelay: 4 * time.Second,
    }
    queueOptions := &warp.QueueOptions{
        Name: "fucker",
        Enabled: true,
        SectionsCount: 16,
        SectionsOffset: 1,
        BatchSize: 64,
        MaxFails: 3,
        BatchDelay: time.Second,
    }
    ctx := context.TODO()

	instance := warpMemory.NewInstance(
        instanceOptions,
        ctx,
    )

    planner, err := implementation.Register[Objective](
        instance,
        queueOptions,
        func (
            manager    warp.HandlerManager, 
            objectives []Objective,
            ctx        context.Context,
        ) error {
            for _, objective := range objectives {
                bus <- objective
            }
            return nil
        },
    )
    if err != nil {
        fmt.Printf("Planner failed %s\n", err)
        return
    }

    go func() {
        instance.Run()
    }()

    go func() {
        for i := 0; i < COUNT; i += 1 {
            objective := Objective{ 
                Id: int64(i), 
                Name: "FUCK YOU",
            }
            planner.Plan(objective, ctx)
        }
    }()

    result := make(map[int64]struct{})

    start := time.Now()

    for dto := range bus {
        result[dto.Id] = struct{}{}

        if len(result) == COUNT {
            break
        }
    }

    fmt.Printf("%s\n", time.Now().Sub(start))
}
