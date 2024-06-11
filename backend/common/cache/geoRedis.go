package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type GeoRedisCacheService struct {
	*RedisBase
}

var _ GeoCacheService = (*GeoRedisCacheService)(nil)

func NewGeoRedisCacheService(client *redis.Client) *GeoRedisCacheService {
	return &GeoRedisCacheService{
		RedisBase: NewRedisBase(client),
	}
}

func (r *GeoRedisCacheService) GetGeoLocation(ctx context.Context, postId string) (*redis.GeoPos, error) {
	geoPos, err := r.client.GeoPos(ctx, "post", postId).Result()
	if err != nil {
		return nil, err
	}

	return geoPos[0], nil
}

func (r *GeoRedisCacheService) AddGeoLocation(ctx context.Context, location *redis.GeoLocation, postId string) error {
	if err := r.client.GeoAdd(ctx, "post", location).Err(); err != nil {
		return err
	}

	return nil
}

func (r *GeoRedisCacheService) RemoveGeoLocation(ctx context.Context, postId string) error {
	if err := r.client.ZRem(ctx, "post", postId).Err(); err != nil {
		return err
	}

	return nil
}

func (r *GeoRedisCacheService) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
	return r.client.GeoRadius(ctx, key, longitude, latitude, query).Result()
}
