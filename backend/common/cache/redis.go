package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisBase struct {
	client *redis.Client
}

func (rb *RedisBase) NewPipe(ctx context.Context) (context.Context, redis.Pipeliner) {
	pipeliner := rb.client.Pipeline()
	ctx = context.WithValue(ctx, PipelinerContextKey, pipeliner)
	return ctx, pipeliner
}

func (rb *RedisBase) GetPipe(ctx context.Context) redis.Pipeliner {
	pipeliner, ok := ctx.Value(PipelinerContextKey).(redis.Pipeliner)
	if !ok {
		return nil
	}
	return pipeliner
}

func (rb *RedisBase) execPipeIfNeeded(ctx context.Context, pipeliner redis.Pipeliner) error {
	if _, ok := ctx.Value(PipelinerContextKey).(redis.Pipeliner); !ok {
		_, err := pipeliner.Exec(ctx)
		return err
	}
	return nil
}

func (rb *RedisBase) IsKeyExist(ctx context.Context, key string) (bool, error) {
	exists, err := rb.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (rb *RedisBase) SetHash(ctx context.Context, key string, values map[string]interface{}) error {
	ctx, pipe := rb.NewPipe(ctx)
	pipe.HSet(ctx, key, values)
	return rb.execPipeIfNeeded(ctx, pipe)
}

func (rb *RedisBase) GetSetMembers(ctx context.Context, key string) ([]string, error) {
	ctx, pipe := rb.NewPipe(ctx)
	cmd := pipe.SMembers(ctx, key)
	err := rb.execPipeIfNeeded(ctx, pipe)
	if err != nil {
		return nil, err
	}
	return cmd.Result()
}

func (rb *RedisBase) AddSetMember(ctx context.Context, key string, member []string) error {
	ctx, pipe := rb.NewPipe(ctx)
	pipe.SAdd(ctx, key, member)
	return rb.execPipeIfNeeded(ctx, pipe)
}

func (rb *RedisBase) RemoveSetMember(ctx context.Context, key string, members []string) error {
	ctx, pipe := rb.NewPipe(ctx)
	pipe.SRem(ctx, key, members)
	return rb.execPipeIfNeeded(ctx, pipe)
}

func (rb *RedisBase) IsMemberInSet(ctx context.Context, key, member string) (bool, error) {
	ctx, pipe := rb.NewPipe(ctx)

	cmd := pipe.SIsMember(ctx, key, member)
	err := rb.execPipeIfNeeded(ctx, pipe)
	if err != nil {
		return false, err
	}

	return cmd.Val(), nil
}

func NewRedisBase(client *redis.Client) *RedisBase {
	return &RedisBase{client: client}
}
