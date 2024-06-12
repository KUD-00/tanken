package cache

import (
	"context"
	dbtype "tanken/backend/common/db"
	"tanken/backend/common/types"

	"github.com/redis/go-redis/v9"
)

type ContextKey string

const PipelinerContextKey ContextKey = "redisPipeliner"

type CacheService interface {
	NewPipe(ctx context.Context) (context.Context, redis.Pipeliner)
	IsKeyExist(ctx context.Context, key string) (bool, error)
	SetHash(ctx context.Context, key string, values map[string]interface{}) error

	GetSetMembers(ctx context.Context, key string) ([]string, error)
	AddSetMember(ctx context.Context, key string, member []string) error
	RemoveSetMember(ctx context.Context, key string, member []string) error
	IsMemberInSet(ctx context.Context, key, member string) (bool, error)
}

type PostCacheService interface {
	CacheService

	// get method default without pipeliner, cmd version with pipeliner.
	// set method default with pipeliner

	GetPostDetails(ctx context.Context, postID string) (*types.PostDetailsPtr, error)
	SetPostDetails(ctx context.Context, postID string, post *types.PostDetailsPtr) error

	GetPostLikedBy(ctx context.Context, postID string) ([]string, error)
	AddPostLikedBy(ctx context.Context, postID string, userIds []string) error
	RemovePostLikedBy(ctx context.Context, postID string, userIds []string) error

	GetPostTags(ctx context.Context, postID string) ([]string, error)
	AddPostTags(ctx context.Context, postID string, tags []string) error
	RemovePostTags(ctx context.Context, postID string, tags []string) error

	GetPostPictureLinks(ctx context.Context, postID string) ([]string, error)
	AddPostPictureLinks(ctx context.Context, postID string, pictureLinks []string) error
	RemovePostPictureLinks(ctx context.Context, postID string, pictureLinks []string) error

	GetPostCommentIds(ctx context.Context, postID string) ([]string, error)
	AddPostCommentIds(ctx context.Context, postID string, commentIds []string) error
	RemovePostCommentIds(ctx context.Context, postID string, commentIds []string) error

	GetComment(ctx context.Context, commentID string) (*types.Comment, error)
	SetComment(ctx context.Context, commentID string, comment *types.Comment) error
	RemoveComment(ctx context.Context, commentID string) error

	SetUser(ctx context.Context, userID string, user *types.User) error
	GetUser(ctx context.Context, userID string) (*types.User, error)
	RemoveUser(ctx context.Context, userID string) error

	AddPostCacheScore(ctx context.Context, postID string, score int64) error
	GetNonPopularPosts(ctx context.Context, limit int64) ([]types.Post, error)
	WriteBackToDB(ctx context.Context, db dbtype.DatabaseService, postIds []string) error
}

type GeoCacheService interface {
	CacheService

	GetGeoLocation(ctx context.Context, postId string) (*redis.GeoPos, error)
	AddGeoLocation(ctx context.Context, location *redis.GeoLocation, postId string) error
	RemoveGeoLocation(ctx context.Context, postId string) error
	GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error)
}

type UserCacheService interface {
	CacheService

	GetUser(ctx context.Context, userID string) (*types.User, error)
	SetUserOptional(ctx context.Context, userId string, user *types.UserPtr) error
	SetUser(ctx context.Context, userId string, user *types.User) error
	RemoveUser(ctx context.Context, userID string) error
}

type CachedPost struct {
	types.Post
	CacheScore int64
	Changed    bool
}

type CachedUser struct {
	types.User
	CacheScore int64
	Changed    bool
}
