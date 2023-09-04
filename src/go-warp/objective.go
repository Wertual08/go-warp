package warp

type Objective[T any] interface {
    Serialize() ([]byte, error) 
    Deserialize(content []byte) (T, error)
    HashCode() int32 
}
