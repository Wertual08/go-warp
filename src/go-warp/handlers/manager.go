package handlers

import (
	"encoding/json"
	"time"

	"github.com/wertual08/go-warp/storage"
)

type managerImpl struct {
	sourceObjectives []storage.ObjectiveDto
	fails            []fail
}

type fail struct {
	index int
    retryDelay time.Duration
	err error
}

func (inst *managerImpl) Fail(index int, err error) {
    fail := fail{
        index: index,
        err: err,
    }

    inst.fails = append(inst.fails, fail)
}
    
func (inst *managerImpl) FailDelay(index int, err error, retryDelay time.Duration) {
    fail := fail{
        index: index,
        retryDelay: retryDelay,
        err: err,
    }

    inst.fails = append(inst.fails, fail)
}

func (inst *managerImpl) Metadata(index int) map[string]string {
    metadata := map[string]string{}

    json.Unmarshal([]byte(inst.sourceObjectives[index].Metadata), &metadata)

    return metadata
}
