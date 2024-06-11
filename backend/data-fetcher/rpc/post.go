package rpc

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"tanken/backend/common/cache"
	database "tanken/backend/common/db"
	"tanken/backend/common/types"
	utils "tanken/backend/common/utils"
	pb "tanken/backend/data-fetcher/rpc/pb"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// TODO: Do I need `defer pipe.Close()`?
func generateUniquePostID(ctx context.Context, rs *cache.GeoRedisCacheService) (string, error) {
	for {
		id := uuid.NewString()[:8]
		exists, err := rs.IsKeyExist(ctx, "post:"+id)

		if err != nil {
			return "", fmt.Errorf("error checking Redis: %v", err)
		}

		if !exists {
			return id, nil
		}
	}
}

func cacheNewPost(ctx context.Context, postID string, content string, userId string, tags []string, rs *cache.PostRedisCacheService) error {
	ctx, pipe := rs.NewPipe(ctx)

	timestamp := utils.Int64Ptr(time.Now().Unix())

	details := &types.PostDetailsPtr{
		CreatedAt:  timestamp,
		UpdatedAt:  timestamp,
		UserId:     utils.StringPtr(userId),
		Content:    utils.StringPtr(content),
		Likes:      utils.Int64Ptr(0),
		Bookmarks:  utils.Int64Ptr(0),
		CacheScore: utils.Int64Ptr(1),
	}

	rs.SetPostDetails(ctx, postID, details)
	rs.AddPostTags(ctx, postID, tags)
	rs.AddPostPictureLinks(ctx, postID, []string{})

	_, err := pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("error caching new post: %v", err)
	}

	return nil
}

func getPost(ctx context.Context, postId string, rs cache.PostCacheService, gc cache.GeoCacheService, db database.DatabaseService) (*types.Post, error) {
	exists, err := rs.IsKeyExist(ctx, "post:"+postId)

	if err != nil {
		return nil, fmt.Errorf("error checking Redis: %v", err)
	}

	if exists {
		post, err := getPostFromCache(ctx, postId, rs, gc)
		if err != nil {
			return nil, err
		}
		return post, nil
	}

	post, err := getPostFromDB(ctx, postId, rs, db)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func getPosts(ctx context.Context, postIds []string, rs cache.PostCacheService, gc cache.GeoCacheService, db database.DatabaseService) ([]*types.Post, error) {
	var posts []*types.Post
	var mu sync.Mutex
	var wg sync.WaitGroup
	errorsChan := make(chan error, len(postIds))

	for _, postId := range postIds {
		wg.Add(1)
		go func(postId string) {
			defer wg.Done()
			post, err := getPost(ctx, postId, rs, gc, db)
			if err != nil {
				errorsChan <- err
				return
			}
			mu.Lock()
			posts = append(posts, post)
			mu.Unlock()
		}(postId)
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		return nil, <-errorsChan
	}

	return posts, nil
}

func getCachedPosts(ctx context.Context, postIds []string, rs cache.PostCacheService, gc cache.GeoCacheService) ([]*types.Post, []string, error) {
	var posts []*types.Post
	var cacheMissedPostIds []string

	for _, postId := range postIds {
		post, err := getPostFromCache(ctx, postId, rs, gc)
		if err != nil {
			cacheMissedPostIds = append(cacheMissedPostIds, postId)
			fmt.Errorf("error retrieving post from Redis: %v", err)
		}

		if post == nil {
			cacheMissedPostIds = append(cacheMissedPostIds, postId)
		} else {
			posts = append(posts, post)
		}
	}

	return posts, cacheMissedPostIds, nil
}

func getPostFromCache(ctx context.Context, postId string, rs cache.PostCacheService, gc cache.GeoCacheService) (*types.Post, error) {
	exists, err := rs.IsKeyExist(ctx, "post:"+postId)
	if err != nil {
		return nil, fmt.Errorf("error checking key existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("post not found in cache")
	}

	prcs, ok := rs.(*cache.PostRedisCacheService)
	if !ok {
		return nil, fmt.Errorf("cache service does not support GetPostDetailsCmd")
	}

	ctx, pipe := rs.NewPipe(ctx)
	postDetailsCmd, _ := prcs.GetPostDetailsCmd(ctx, postId)
	tagsCmd, _ := prcs.GetPostTagsCmd(ctx, postId)
	pictureLinksCmd, _ := prcs.GetPostPictureLinksCmd(ctx, postId)
	commentsCmd, _ := prcs.GetPostCommentIdsCmd(ctx, postId)
	likedByCmd, _ := prcs.GetPostLikedByCmd(ctx, postId)

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	post := postDetailsCmd.Val()
	tags := tagsCmd.Val()
	pictureLinks := pictureLinksCmd.Val()
	comments := commentsCmd.Val()
	likedBy := likedByCmd.Val()

	geoPos, err := gc.GetGeoLocation(ctx, postId)
	if err != nil {
		return nil, err
	}

	location := types.Location{
		Latitude:  geoPos.Latitude,
		Longitude: geoPos.Longitude,
	}

	return &types.Post{
		//TODO: make this as function
		PostDetails: types.PostDetails{
			PostId:    postId,
			Location:  location,
			CreatedAt: utils.StringToInt64(post["CreatedAt"], 0),
			UpdatedAt: utils.StringToInt64(post["UpdatedAt"], 0),
			UserId:    post["UserId"],
			Content:   post["Content"],
			Likes:     utils.StringToInt64(post["Likes"], 0),
			Bookmarks: utils.StringToInt64(post["Bookmarks"], 0),
		},
		PostSets: types.PostSets{
			Tags:         tags,
			PictureLinks: pictureLinks,
			CommentIds:   comments,
			LikedBy:      likedBy,
		},
	}, nil
}

func getPostFromDB(ctx context.Context, postId string, rs cache.PostCacheService, db database.DatabaseService) (*types.Post, error) {
	var post types.Post
	var location types.Location

	db.GetPostDetails(ctx, postId)

	details := types.PostDetailsPtr{
		CreatedAt: &post.CreatedAt,
		UpdatedAt: &post.UpdatedAt,
		UserId:    &post.UserId,
		Content:   &post.Content,
		Likes:     &post.Likes,
		Bookmarks: &post.Bookmarks,
	}

	ctx, pipe := rs.NewPipe(ctx)

	rs.SetPostDetails(ctx, postId, &details)
	rs.AddPostLikedBy(ctx, postId, post.LikedBy)
	rs.AddPostCommentIds(ctx, postId, post.CommentIds)
	rs.AddPostTags(ctx, postId, post.Tags)
	rs.AddPostPictureLinks(ctx, postId, post.PictureLinks)

	_, err := pipe.Exec(ctx)

	if err != nil {
		return &types.Post{}, err
	}

	//TODO: should i only save location in geo-redis?
	//TODO: maybe location geo-redis need to backup one day.

	return &types.Post{
		PostDetails: types.PostDetails{
			PostId:    postId,
			Location:  location,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			UserId:    post.UserId,
			Content:   post.Content,
			Likes:     post.Likes,
			Bookmarks: post.Bookmarks,
		},
		PostSets: types.PostSets{
			Tags:         post.Tags,
			PictureLinks: post.PictureLinks,
			CommentIds:   post.CommentIds,
		},
	}, nil
}

func removePostFromCache(ctx context.Context, postId string, rs cache.PostCacheService) error {
	ctx, pipe := rs.NewPipe(ctx)

	pipe.Del(ctx, "post:"+postId)
	pipe.Del(ctx, "post:"+postId+":likedBy")
	pipe.Del(ctx, "post:"+postId+":tags")
	pipe.Del(ctx, "post:"+postId+":pictureLinks")
	pipe.Del(ctx, "post:"+postId+":comments") //TODO: need to delete comments in cache

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error removing post from Redis: %v", err)
	}

	return nil
}

func removePostFromDB(ctx context.Context, postId string, postgres *sql.DB) error {
	_, err := postgres.ExecContext(ctx, "DELETE FROM posts WHERE post_id = $1", postId)
	if err != nil {
		return fmt.Errorf("error deleting post from PostgreSQL: %v", err)
	}

	return nil
}

// Others
func getPostGeoData(ctx context.Context, postID string, gc cache.GeoCacheService) (types.Location, error) {
	geoPos, err := gc.GetGeoLocation(ctx, postID)
	if err != nil {
		return types.Location{}, fmt.Errorf("error retrieving geo data from Redis: %v", err)
	}

	return types.Location{
		Latitude:  geoPos.Latitude,
		Longitude: geoPos.Longitude,
	}, nil
}

// TODO: check if is liked
func incrementLikes(ctx context.Context, postId string, userId string, rs cache.PostCacheService, gc cache.GeoCacheService, db database.DatabaseService) error {
	post, err := getPost(ctx, postId, rs, gc, db)

	if err != nil {
		return fmt.Errorf("error retrieving post: %v", err)
	}

	newLikes := post.Likes + 1

	details := types.PostDetailsPtr{
		Likes: &newLikes,
	}

	ctx, pipe := rs.NewPipe(ctx)
	rs.SetPostDetails(ctx, postId, &details)
	rs.AddPostLikedBy(ctx, postId, []string{userId})

	_, err = pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("error incrementing likes in Redis: %v", err)
	}

	return nil
}

// TODO: check if not liked
func decrementLikes(ctx context.Context, postId string, userId string, rs cache.PostCacheService, gc cache.GeoCacheService, db database.DatabaseService) error {
	post, err := getPost(ctx, postId, rs, gc, db)
	if err != nil {
		return fmt.Errorf("error retrieving post details from Redis: %v", err)
	}

	newLikes := post.Likes - 1

	details := types.PostDetailsPtr{
		Likes: &newLikes,
	}

	ctx, pipe := rs.NewPipe(ctx)
	rs.SetPostDetails(ctx, postId, &details)
	rs.RemovePostLikedBy(ctx, postId, []string{userId})

	_, err = pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("error decrementing likes in Redis: %v", err)
	}

	return nil
}

func registerPostToGeoRedis(ctx context.Context, postID string, longitude, latitude float64, rs *cache.GeoRedisCacheService) error {
	geoLocation := &redis.GeoLocation{
		Name:      postID,
		Longitude: float64(longitude),
		Latitude:  float64(latitude),
	}

	err := rs.AddGeoLocation(ctx, geoLocation, postID)

	if err != nil {
		return fmt.Errorf("error registering post to geo-redis: %v", err)
	}

	return nil
}

func removePostIDinGeoRedis(ctx context.Context, postID string, rs *cache.GeoRedisCacheService) error {
	if err := rs.RemoveGeoLocation(ctx, postID); err != nil {
		fmt.Errorf("error removing geo data: %v", err)
		return err
	}
	return nil
}

/*
	func (s *server) getCachedPbPostsAndMissedPostIds(ctx context.Context, postIds []string) ([]*pb.Post, []string, error) {
		var posts []*pb.Post
		var cacheMissedPostIds []string

		for _, postId := range postIds {
			post, err := s.post_cache_rdb.HGetAll(ctx, "post:"+postId).Result()
			if err != nil {
				return nil, nil, fmt.Errorf("Error retrieving post from Redis: %v", err)
			}
			if len(post) == 0 {
				cacheMissedPostIds = append(cacheMissedPostIds, postId)
			} else {
				transformedPost := s.transformCachedPostIdToPbPost(ctx, postId)
				posts = append(posts, transformedPost)
			}
		}

		return posts, cacheMissedPostIds, nil
	}

	func (s *server) getGeoPostIDs(ctx context.Context, req *connect.Request[pb.GetPostsByLocationRequest]) ([]redis.GeoLocation, error) {
		return s.geo_postid_rdb.GeoRadius(ctx, "geo-postid", float64(req.Msg.Location.Longitude), float64(req.Msg.Location.Latitude), &redis.GeoRadiusQuery{
			Radius:    float64(req.Msg.Radius),
			Unit:      "km",
			WithDist:  true,
			WithCoord: true,
			Count:     100,
			Sort:      "ASC",
		}).Result()
	}
*/
func newUploadErrorResponse(message string) *connect.Response[pb.AddPostResponse] {
	return connect.NewResponse(&pb.AddPostResponse{
		Ok:  0,
		Msg: message,
	})
}
