package pubsub

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type RedisDispatcher struct {
	nodeId string
	ctx context.Context
	rdb *redis.Client
}

func (d *RedisDispatcher) Dispatch(topic string, message Message) error {
	bytes, err := msgpack.Marshal(message)

	if err != nil {
		return err
	}

	return d.rdb.Publish(d.ctx, topic, bytes).Err()
}

func (d *RedisDispatcher) DispatchTo(target MessageTarget, message Message) error {
	return d.Dispatch(target.Build(), message)
}

func NewRedisDispatcher(ctx context.Context, nodeId string, rdb *redis.Client) *RedisDispatcher {
	return &RedisDispatcher{
		ctx:             ctx,
		nodeId:          nodeId,
		rdb:             rdb,
	}
}