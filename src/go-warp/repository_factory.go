package warp

import "github.com/wertual08/go-warp/storage"

type RepositoryFactory struct {
    Queue storage.QueueRepository
    Objective storage.ObjectiveRepository
}
