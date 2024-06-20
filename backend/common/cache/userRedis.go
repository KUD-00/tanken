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
		UserId:             userId,
		Username:           user["username"],
		Bio:                user["bio"],
		ProfilePictureLink: user["profilePictureLink"],
		Subscribed:         utils.StringToInt64(user["subscribed"], 0),
	}, nil
}

func (r *UserRedisCacheService) SetUserOptional(ctx context.Context, userId string, user *types.UserPtr) error {
	var changed bool

	if user.Changed == nil {
		changed = false
	} else {
		changed = *user.Changed
	}

	data := make(map[string]interface{})
	if user.Username != nil {
		data["username"] = *user.Username
	}
	if user.Bio != nil {
		data["bio"] = *user.Bio
	}
	if user.ProfilePictureLink != nil {
		data["profilePictureLink"] = *user.ProfilePictureLink
	}
	if user.Subscribed != nil {
		data["subscribed"] = *user.Subscribed
	}
	data["changed"] = changed

	_, err := r.client.HSet(ctx, "user:"+userId, data).Result()

	return err
}

func (r *UserRedisCacheService) SetUser(ctx context.Context, userId string, user *types.User) error {
	_, err := r.client.HSet(ctx, "user:"+userId, map[string]interface{}{
		"user_id":            user.UserId,
		"username":           user.Username,
		"bio":                user.Bio,
		"profilePictureLink": user.ProfilePictureLink,
		"subscribed":         user.Subscribed,
	}).Result()

	return err
}

func (r *UserRedisCacheService) RemoveUser(ctx context.Context, userId string) error {
	_, err := r.client.Del(ctx, "user:"+userId).Result()

	return err
}
