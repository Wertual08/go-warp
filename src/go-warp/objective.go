package warp;

type Objective interface {
    Serialize() []byte
    Deserialize(content []byte) error
    HashCode() int32
}
