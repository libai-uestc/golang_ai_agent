package checkpointstore

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

type InMemoryStore struct {
	mem map[string][]byte
}

func (i *InMemoryStore) Set(ctx context.Context, key string, value []byte) error {
	i.mem[key] = value
	return nil
}

func (i *InMemoryStore) Get(ctx context.Context, key string) ([]byte, bool, error) {
	v, ok := i.mem[key]
	return v, ok, nil
}

func NewInMemoryStore() compose.CheckPointStore {
	return &InMemoryStore{
		mem: map[string][]byte{},
	}
}

/*
创建一个类，实现接口

	type CheckPointStore interface {
		Get(ctx context.Context, checkPointID string) ([]byte, bool, error)
		Set(ctx context.Context, checkPointID string, checkPoint []byte) error
	}
*/
