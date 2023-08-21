package controller

import (
	"github.com/wertual08/go-warp"
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
