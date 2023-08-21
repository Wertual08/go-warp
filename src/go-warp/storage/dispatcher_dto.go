package storage

import "github.com/google/uuid"

type DispatcherDto struct {
    Id uuid.UUID
    Stride int32
    Offset int32
}
/*
Select all
If me invalid {
    Stop processing
    Wait until finished
    Update myself
}
Insert myself
If all are valid {
    Start processing
}
*/
