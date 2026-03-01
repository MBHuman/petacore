package serializers

import "petacore/sdk/pmem"

type BaseSerializer[GoT any, T any] interface {
	Serializer[GoT]
	Deserializer[T]
	Validator[T]
}

type Serializer[T any] interface {
	Serialize(allocator pmem.Allocator, value T) ([]byte, error)
}

type Deserializer[T any] interface {
	Deserialize([]byte) (T, error)
}

type Validator[T any] interface {
	Validate(T) error
}
