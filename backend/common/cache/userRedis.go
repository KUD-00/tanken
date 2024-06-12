package cache

import (
	"context"
	"tanken/backend/common/types"
	"tanken/backend/common/utils"

	"github.com/redis/go-redis/v9"
)

type UserRedisCacheService struct {
	*RedisBase
}

var _ UserCacheService = (*UserRedisCacheService)(nil)

func NewUserRedisCacheService(client *redis.Client) *UserRedisCacheService {
	return &UserRedisCacheService{
		RedisBase: NewRedisBase(client),
	}
}

func (r *UserRedisCacheService) GetUser(ctx context.Context, userId string) (*types.User, error) {
	user, err := r.client.HGetAll(ctx, "user:"+userId).Result()
	if err != nil {
		return nil, err
	}

	return &types.User{
		UserId:     user["user_id"],
		Username:   user["username"],
		Bio:        user["bio"],
		Avatar:     user["avatar"],
		Subscribed: utils.StringToInt64(user["subscribed"], 0),
	}, nil
}

func (r *UserRedisCacheService) SetUserOptional(ctx context.Context, userId string, user *types.UserPtr) error {
	_, err := r.client.HSet(ctx, "user:"+userId, map[string]interface{}{
		"username":   *user.Username,
		"bio":        *user.Bio,
		"avatar":     *user.Avatar,
		"subscribed": *user.Subscribed,
		"changed":    *user.Changed,
	}).Result()

	return err
}

func (r *UserRedisCacheService) SetUser(ctx context.Context, userId string, user *types.User) error {
	_, err := r.client.HSet(ctx, "user:"+userId, map[string]interface{}{
		"user_id":    user.UserId,
		"username":   user.Username,
		"bio":        user.Bio,
		"avatar":     user.Avatar,
		"subscribed": user.Subscribed,
	}).Result()

	return err
}

func (r *UserRedisCacheService) RemoveUser(ctx context.Context, userId string) error {
	_, err := r.client.Del(ctx, "user:"+userId).Result()

	return err
}
