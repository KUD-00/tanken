package rpc

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"tanken/backend/common/cache"
	database "tanken/backend/common/db"
	"tanken/backend/common/types"
	pb "tanken/backend/data-fetcher/rpc/pb"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// TODO: Do I need `defer pipe.Close()`?
func generateUniquePostID(ctx context.Context, rs *cache.GeoRedisCacheService) (string, error) {
	for {
		id := uuid.NewString()
		exists, err := rs.IsKeyExist(ctx, "post:"+id)

		if err != nil {
			return "", fmt.Errorf("error checking Redis: %v", err)
		}

		if !exists {
			return id, nil
		}
	}
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

	post, err := db.GetPost(ctx, postId)
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
	post, err := rs.GetPost(ctx, postId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving post from Redis: %v", err)
	}

	geoPos, err := gc.GetGeoLocation(ctx, postId)
	if err != nil {
		return nil, err
	}

	location := types.Location{
		Latitude:  geoPos.Latitude,
		Longitude: geoPos.Longitude,
	}

	post.Location = location

	return post, nil
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

	cmds, _ := pipe.Exec(ctx)

	for i, cmd := range cmds {
		if cmd.Err() != nil {
			switch i {
			case 0:
				return fmt.Errorf("error setting post details: %v", cmd.Err())
			case 1:
				return fmt.Errorf("error decrementing post likedby: %v", cmd.Err())
			}
		}
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

func uploadPictureToS3(ctx context.Context, pictureChunk []byte, s3Uploader *s3manager.Uploader, key string) (string, error) {
	bucketName := os.Getenv("POST_PICTURE_BUCKET_NAME")

	if bucketName == "" || key == "" {
		log.Fatalf("BUCKET_NAME or KEY environment variable is not set")
	}

	input := &s3manager.UploadInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(pictureChunk)),
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	_, err := s3Uploader.UploadWithContext(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload picture to S3: %w", err)
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
	return s3URL, nil
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
