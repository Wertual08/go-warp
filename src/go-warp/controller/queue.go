package controller

import (
	"context"
	"errors"
	"time"

	"github.com/wertual08/go-warp"
	"github.com/wertual08/go-warp/storage"
)

type Queue struct {
    Id int32
    Options *warp.QueueOptions
    HandlerFactory func() Handler
    Handlers []Handler
}

func CreateQueue[T warp.Objective](
    queueId int32,
    options *warp.QueueOptions,
    handler warp.Handler[T],
) Queue {
    return Queue{
        Id: queueId,
        Options: options,
        HandlerFactory: createHandlerFactory(handler),
    }
}

func (inst *Queue) Handle(
    stride int32,
    offset int32,
    objectiveRepository storage.ObjectiveRepository,
    ctx context.Context,
) error {
    if !inst.Options.Enabled {
        time.Sleep(inst.Options.BatchDelay)
        return nil
    }

    now := time.Now()
    sectionsCount := (inst.Options.SectionsCount + inst.Options.SectionsOffset + stride - 1 - offset) / stride

    result := make(chan error)
    skipDelay := false

    for i := int32(0); i < sectionsCount; i += 1 {
        if int(i) == len(inst.Handlers) {
            inst.Handlers = append(inst.Handlers, inst.HandlerFactory())
        }
        
        go handleSection(
            result, 
            inst.Handlers[i], 
            inst.Id,
            offset + stride * i, 
            inst.Options.BatchSize,
            now,
            &skipDelay,
            objectiveRepository,
            ctx,
        )
    }

    var resultError error = nil

    for i := 0; i < int(sectionsCount); i += 1 {
        err := <- result
        resultError = errors.Join(resultError, err)
    }

    if resultError == nil && !skipDelay {
        time.Sleep(inst.Options.BatchDelay)
    }

    return resultError
}

func handleSection(
    result chan error, 
    handler Handler,
    queueId int32,
    section int32,
    limit int32,
    now time.Time,
    skipDelay *bool,
    objectiveRepository storage.ObjectiveRepository,
    ctx context.Context,
) {
    objectives, err := objectiveRepository.List(
        queueId,
        section,
        limit,
        now,
        ctx,
    )
    if err != nil {
        result <- err
        return
    }

    handler.Handle(objectives, ctx)

    if handler.Succeded() == int(limit) {
        *skipDelay = true
    }

    if err := objectiveRepository.Create(handler.Failed(), ctx); err != nil {
        result <- err
        return
    }
    if err := objectiveRepository.Remove(objectives, ctx); err != nil {
        result <- err
        return
    }

    result <- nil
}
