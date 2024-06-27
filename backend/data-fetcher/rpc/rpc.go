// 将数据访问服务进一步拆分，细化到不同的业务逻辑和数据类型，减少单个服务的负担。
// 对于不需要实时处理的请求，考虑使用异步处理和消息队列，减少同步调用的压力。
package rpc

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"tanken/backend/common/cache"
	postgres "tanken/backend/common/db/postgres"
	types "tanken/backend/common/types"
	commonUtils "tanken/backend/common/utils"
	"tanken/backend/data-fetcher/rpc/connectrpc/pbconnect"
	pb "tanken/backend/data-fetcher/rpc/pb"
	utils "tanken/backend/data-fetcher/utils"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type server struct {
	db         *postgres.PostgresDatabaseService
	postCache  *cache.PostRedisCacheService
	geoCache   *cache.GeoRedisCacheService
	userCache  *cache.UserRedisCacheService
	uploaderS3 *s3manager.Uploader
}

func (s *server) GetPostsByLocation(ctx context.Context, req *connect.Request[pb.GetPostsByLocationRequest]) (*connect.Response[pb.GetPostsByLocationResponse], error) {
	// Get the posts from the geo redis
	postGeoDatas, err := s.geoCache.GeoRadius(ctx, "post", req.Msg.Location.Longitude, req.Msg.Location.Latitude, &redis.GeoRadiusQuery{
		Radius:      float64(req.Msg.Radius),
		Unit:        "km",
		WithDist:    true,
		WithCoord:   true,
		WithGeoHash: false,
		Count:       100,
		Sort:        "ASC",
	})

	if err != nil {
		return nil, fmt.Errorf("error querying geo redis: %v", err)
	}

	//TODO: 聚类查询+算法得到应该返回的帖子ID

	var postIds []string
	for _, geoLoc := range postGeoDatas {
		if geoLoc.Name != "" {
			postIds = append(postIds, geoLoc.Name)
		}
	}

	posts, err := getPosts(ctx, postIds, s.postCache, s.geoCache, s.db)

	if err != nil {
		return nil, fmt.Errorf("error getting posts: %v", err)
	}

	return connect.NewResponse(&pb.GetPostsByLocationResponse{Ok: 1, Posts: utils.CommonPostsToPbPosts(posts)}), nil
}

func (s *server) GetPostsByUser(ctx context.Context, req *connect.Request[pb.GetPostsByUserIdRequest]) (*connect.Response[pb.GetPostsByUserIdResponse], error) {
	/*   索引优化：
	     在 userId 和 timestamp 上创建复合索引。这样可以快速筛选出特定用户的帖子，并且根据时间戳排序。
	     索引应该是 (userId, timestamp DESC)，这样可以直接按照时间戳降序排列，方便查询最新的帖子。 */

	/* 	query := `
	    SELECT post_id, created_at, updated_at, author, content, likes, bookmarks, picture_links, location, comments
	    FROM posts
	    WHERE user_id = $1 AND timestamp <= $2
	    ORDER BY timestamp DESC
	    LIMIT 10
	    `
		rows, err := s.db.QueryContext(ctx, query, req.Msg.UserId, req.Msg.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("error querying posts: %v", err)
		}
		defer rows.Close()

		var posts []*pb.Post

		for rows.Next() {
			var post pb.Post
			if err := rows.Scan(&post.PostId, &post.CreatedAt, &post.UpdatedAt, &post.Author, &post.Content, &post.Likes, &post.Bookmarks, &post.Location, &post.Comments); err != nil {
				return nil, fmt.Errorf("error scanning post: %v", err)
			}
			posts = append(posts, &post)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error reading rows: %v", err)
		}

		return connect.NewResponse(&pb.GetPostsByUserIdResponse{Posts: posts}), nil */
	return connect.NewResponse(&pb.GetPostsByUserIdResponse{Posts: nil}), nil
}

// Integrated tested
func (s *server) GetPostsByPostIds(ctx context.Context, req *connect.Request[pb.GetPostsByPostIdsRequest]) (*connect.Response[pb.GetPostsByPostIdsResponse], error) {
	posts, err := getPosts(ctx, req.Msg.PostIds, s.postCache, s.geoCache, s.db)

	if err != nil {
		return connect.NewResponse(&pb.GetPostsByPostIdsResponse{
			Ok:  0,
			Msg: "error getting posts: " + err.Error(),
		}), nil
	}

	return connect.NewResponse(&pb.GetPostsByPostIdsResponse{
		Ok:    1,
		Posts: utils.CommonPostsToPbPosts(posts),
	}), nil
}

// Integrated tested
func (s *server) AddPost(ctx context.Context, req *connect.Request[pb.AddPostRequest]) (*connect.Response[pb.AddPostResponse], error) {
	longitude := req.Msg.Location.Longitude
	latitude := req.Msg.Location.Latitude

	postId, err := generateUniquePostID(ctx, s.geoCache)
	if err != nil {
		return newUploadErrorResponse("Failed to get or register new postId: " + err.Error()), nil
	}

	pictureLink, err := uploadPictureToS3(ctx, req.Msg.PictureChunk, s.uploaderS3, postId)
	if err != nil {
		return newUploadErrorResponse("Failed to upload picture to S3: " + err.Error()), nil
	}

	err = registerPostToGeoRedis(ctx, postId, longitude, latitude, s.geoCache)
	if err != nil {
		return newUploadErrorResponse("Failed to register new postId to geo-redis: " + err.Error()), nil
	}

	ctx, pipe := s.postCache.NewPipe(ctx)

	timestamp := commonUtils.Int64Ptr(time.Now().Unix())

	details := &types.PostDetailsPtr{
		CreatedAt:  timestamp,
		UpdatedAt:  timestamp,
		UserId:     commonUtils.StringPtr(req.Msg.UserId),
		Content:    commonUtils.StringPtr(req.Msg.Content),
		Likes:      commonUtils.Int64Ptr(0),
		Bookmarks:  commonUtils.Int64Ptr(0),
		CacheScore: commonUtils.Int64Ptr(1),
	}

	s.postCache.SetPostDetails(ctx, postId, details)
	s.postCache.AddPostTags(ctx, postId, req.Msg.Tags)
	s.postCache.AddPostPictureLinks(ctx, postId, []string{pictureLink})

	_, err = pipe.Exec(ctx)

	// Attempt to rollback if caching fails
	if err != nil {
		if err = removePostIDinGeoRedis(ctx, postId, s.geoCache); err != nil {
			return newUploadErrorResponse("Failed to cache post details and failed to rollback geoRedis: " + err.Error()), nil
		}
		return newUploadErrorResponse("Failed to cache post details, rollback geoRedis is success: " + err.Error()), nil
	}

	return connect.NewResponse(&pb.AddPostResponse{Ok: 1, PostId: postId}), nil
}

func (s *server) HardDeletePost(ctx context.Context, req *connect.Request[pb.HardDeletePostRequest]) (*connect.Response[pb.HardDeletePostResponse], error) {
	err := s.postCache.RemovePost(ctx, req.Msg.PostId)
	if err != nil {
		return connect.NewResponse(&pb.HardDeletePostResponse{Ok: 0, Msg: "Error in hard delete post in cache:" + err.Error()}), nil
	}

	err = s.db.HardDeletePost(ctx, req.Msg.PostId)
	if err != nil {
		return connect.NewResponse(&pb.HardDeletePostResponse{Ok: 0, Msg: "Error in hard delete post in db:" + err.Error()}), nil
	}

	return connect.NewResponse(&pb.HardDeletePostResponse{Ok: 1}), nil
}

func (s *server) SoftDeletePost(ctx context.Context, req *connect.Request[pb.SoftDeletePostRequest]) (*connect.Response[pb.SoftDeletePostResponse], error) {
	err := s.postCache.RemovePost(ctx, req.Msg.PostId)
	if err != nil {
		return connect.NewResponse(&pb.SoftDeletePostResponse{Ok: 0, Msg: "Error in soft delete post in cache:" + err.Error()}), nil
	}

	err = s.db.SoftDeletePost(ctx, req.Msg.PostId)
	if err != nil {
		return connect.NewResponse(&pb.SoftDeletePostResponse{Ok: 0, Msg: "Error in soft delete post in db:" + err.Error()}), nil
	}

	return connect.NewResponse(&pb.SoftDeletePostResponse{Ok: 1}), nil
}

// Integrated tested
func (s *server) AddLike(ctx context.Context, req *connect.Request[pb.AddLikeRequest]) (*connect.Response[pb.AddLikeResponse], error) {
	if err := incrementLikes(ctx, req.Msg.PostId, req.Msg.UserId, s.postCache, s.geoCache, s.db); err != nil {
		return connect.NewResponse(&pb.AddLikeResponse{Ok: 0, Msg: err.Error()}), nil
	}
	return connect.NewResponse(&pb.AddLikeResponse{Ok: 1}), nil
}

// Integrated tested
func (s *server) RemoveLike(ctx context.Context, req *connect.Request[pb.RemoveLikeRequest]) (*connect.Response[pb.RemoveLikeResponse], error) {
	if err := decrementLikes(ctx, req.Msg.PostId, req.Msg.UserId, s.postCache, s.geoCache, s.db); err != nil {
		return connect.NewResponse(&pb.RemoveLikeResponse{Ok: 0, Msg: err.Error()}), nil
	}
	return connect.NewResponse(&pb.RemoveLikeResponse{Ok: 1}), nil
}

// Integrated tested
func (s *server) AddComment(ctx context.Context, req *connect.Request[pb.AddCommentRequest]) (*connect.Response[pb.AddCommentResponse], error) {
	commentId, err := generateUniqueCommentID(ctx, req.Msg.PostId, s.postCache, s.db)

	if err != nil {
		return connect.NewResponse(&pb.AddCommentResponse{Ok: 0, Msg: "err generating unique comment ID" + err.Error()}), nil
	}

	comment := types.Comment{
		CommentId: commentId,
		PostId:    req.Msg.PostId,
		UserId:    req.Msg.UserId,
		Content:   req.Msg.Content,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Likes:     0,
		Status:    int64(1),
	}

	err = s.postCache.SetComment(ctx, commentId, &comment)
	if err != nil {
		return connect.NewResponse(&pb.AddCommentResponse{Ok: 0, Msg: err.Error()}), nil
	}

	return connect.NewResponse(&pb.AddCommentResponse{Ok: 1}), nil
}

func (s *server) HardDeleteComment(ctx context.Context, req *connect.Request[pb.HardDeleteCommentRequest]) (*connect.Response[pb.HardDeleteCommentResponse], error) {
	ctx, pipe := s.postCache.NewPipe(ctx)
	s.postCache.RemoveComments(ctx, []string{req.Msg.CommentId})
	s.postCache.RemovePostCommentIds(ctx, req.Msg.PostId, []string{req.Msg.CommentId})

	_, err := pipe.Exec(ctx)
	if err != nil {
		return connect.NewResponse(&pb.HardDeleteCommentResponse{Ok: 0, Msg: "Error in hard delete comment in cache:" + err.Error()}), nil
	}

	err = s.db.HardDeleteCommentById(ctx, req.Msg.CommentId)
	if err != nil {
		return connect.NewResponse(&pb.HardDeleteCommentResponse{Ok: 0, Msg: "Error in hard delete comment in db:" + err.Error()}), nil
	}

	return connect.NewResponse(&pb.HardDeleteCommentResponse{Ok: 1}), nil
}

func (s *server) SoftDeleteComment(ctx context.Context, req *connect.Request[pb.SoftDeleteCommentRequest]) (*connect.Response[pb.SoftDeleteCommentResponse], error) {
	ctx, pipe := s.postCache.NewPipe(ctx)
	s.postCache.RemoveComments(ctx, []string{req.Msg.CommentId})
	s.postCache.RemovePostCommentIds(ctx, req.Msg.PostId, []string{req.Msg.CommentId})

	_, err := pipe.Exec(ctx)
	if err != nil {
		return connect.NewResponse(&pb.SoftDeleteCommentResponse{Ok: 0, Msg: "Error in hard delete comment in cache:" + err.Error()}), nil
	}

	err = s.db.SoftDeleteCommentById(ctx, req.Msg.CommentId)
	if err != nil {
		return connect.NewResponse(&pb.SoftDeleteCommentResponse{Ok: 0, Msg: "Error in hard delete comment in db:" + err.Error()}), nil
	}

	return connect.NewResponse(&pb.SoftDeleteCommentResponse{Ok: 1}), nil
}

// TODO : not completed
func (s *server) UpdateComment(ctx context.Context, req *connect.Request[pb.UpdateCommentRequest]) (*connect.Response[pb.UpdateCommentResponse], error) {
	// if cached, just modify in cache
	if exist, err := s.postCache.IsKeyExist(ctx, req.Msg.CommentId); err == nil && exist {
		comment, err := s.postCache.GetComment(ctx, req.Msg.CommentId)
		if err != nil {
			return connect.NewResponse(&pb.UpdateCommentResponse{}), nil
		}

		comment.Content = req.Msg.Content
		comment.UpdatedAt = time.Now().Unix()

		err = s.postCache.SetComment(ctx, req.Msg.CommentId, comment)
		if err != nil {
			return connect.NewResponse(&pb.UpdateCommentResponse{}), nil
		}
	}

	// TODO: if not cached, do we need cache? hard to say but... okay just modify it in db. or if the post is cached, we will cache it maybe?
	return connect.NewResponse(&pb.UpdateCommentResponse{}), nil
}

func (s *server) GetCommentsByCommentIds(ctx context.Context, req *connect.Request[pb.GetCommentsByCommentIdsRequest]) (*connect.Response[pb.GetCommentsByCommentIdsResponse], error) {
	return connect.NewResponse(&pb.GetCommentsByCommentIdsResponse{}), nil
}

func (s *server) GetCommentsByPostId(ctx context.Context, req *connect.Request[pb.GetCommentsByPostIdRequest]) (*connect.Response[pb.GetCommentsByPostIdResponse], error) {
	return connect.NewResponse(&pb.GetCommentsByPostIdResponse{}), nil
}

func (s *server) AddBookmark(ctx context.Context, req *connect.Request[pb.AddBookmarkRequest]) (*connect.Response[pb.AddBookmarkResponse], error) {
	return connect.NewResponse(&pb.AddBookmarkResponse{}), nil
}

func (s *server) RemoveBookmark(ctx context.Context, req *connect.Request[pb.RemoveBookmarkRequest]) (*connect.Response[pb.RemoveBookmarkResponse], error) {
	return connect.NewResponse(&pb.RemoveBookmarkResponse{}), nil
}

func (s *server) WriteBackCache(ctx context.Context, req *connect.Request[pb.WriteBackCacheRequest]) (*connect.Response[pb.WriteBackCacheResponse], error) {
	return connect.NewResponse(&pb.WriteBackCacheResponse{}), nil
}

func (s *server) GetUserInfoByOAuth(ctx context.Context, req *connect.Request[pb.GetUserInfoByOAuthRequest]) (*connect.Response[pb.GetUserInfoByOAuthResponse], error) {
	user, err := s.db.GetUserByOauthInfo(ctx, req.Msg.Email, req.Msg.Provider)
	if err != nil {
		return connect.NewResponse(&pb.GetUserInfoByOAuthResponse{Ok: 0, Msg: err.Error()}), nil
	}

	pbUser := &pb.User{
		UserId:             user.UserId,
		UserName:           user.Username,
		Bio:                user.Bio,
		ProfilePictureLink: user.ProfilePictureLink,
		Subscribed:         user.Subscribed,
	}

	return connect.NewResponse(&pb.GetUserInfoByOAuthResponse{Ok: 1, User: pbUser}), nil
}

func (s *server) HardDeleteUser(ctx context.Context, req *connect.Request[pb.HardDeleteUserRequest]) (*connect.Response[pb.HardDeleteUserResponse], error) {
	err := s.userCache.RemoveUser(ctx, req.Msg.UserId)
	if err != nil {
		return connect.NewResponse(&pb.HardDeleteUserResponse{Ok: 0, Msg: "Error in hard delete user in cache:" + err.Error()}), nil
	}

	err = s.db.HardDeleteUserById(ctx, req.Msg.UserId)
	if err != nil {
		return connect.NewResponse(&pb.HardDeleteUserResponse{Ok: 0, Msg: "Error in hard delete user in db:" + err.Error()}), nil
	}

	return connect.NewResponse(&pb.HardDeleteUserResponse{}), nil
}

func (s *server) SoftDeleteUser(ctx context.Context, req *connect.Request[pb.SoftDeleteUserRequest]) (*connect.Response[pb.SoftDeleteUserResponse], error) {
	err := s.userCache.RemoveUser(ctx, req.Msg.UserId)
	if err != nil {
		return connect.NewResponse(&pb.SoftDeleteUserResponse{Ok: 0, Msg: "Error in soft delete user in cache:" + err.Error()}), nil
	}

	err = s.db.SoftDeleteUserById(ctx, req.Msg.UserId)
	if err != nil {
		return connect.NewResponse(&pb.SoftDeleteUserResponse{Ok: 0, Msg: "Error in soft delete user in db:" + err.Error()}), nil
	}
	return connect.NewResponse(&pb.SoftDeleteUserResponse{}), nil
}

// Integrated tested
func (s *server) GetUserInfo(ctx context.Context, req *connect.Request[pb.GetUserInfoRequest]) (*connect.Response[pb.GetUserInfoResponse], error) {
	user, err := getUser(ctx, req.Msg.UserId, s.userCache, s.db)

	if err != nil {
		return connect.NewResponse(&pb.GetUserInfoResponse{Ok: 0, Msg: err.Error()}), nil
	}

	pbUser := &pb.User{
		UserId:             user.UserId,
		UserName:           user.Username,
		Bio:                user.Bio,
		ProfilePictureLink: user.ProfilePictureLink,
		Subscribed:         user.Subscribed,
	}

	return connect.NewResponse(&pb.GetUserInfoResponse{Ok: 1, User: pbUser}), nil
}

// Integrated tested
func (s *server) SignUpUser(ctx context.Context, req *connect.Request[pb.SignUpUserRequest]) (*connect.Response[pb.SignUpUserResponse], error) {
	userId, err := generateUniqueUserID(ctx, s.db)

	if err != nil {
		return connect.NewResponse(&pb.SignUpUserResponse{Ok: 0, Msg: err.Error()}), nil
	}

	user := types.UserPtr{
		Username:           &req.Msg.Name,
		Email:              &req.Msg.Email,
		Bio:                &req.Msg.Bio,
		Subscribed:         commonUtils.Int64Ptr(0),
		ProfilePictureLink: &req.Msg.ProfilePictureLink,
		OauthProvider:      &req.Msg.Provider,
	}

	err = s.db.SetUserById(ctx, userId, &user)
	if err != nil {
		return connect.NewResponse(&pb.SignUpUserResponse{Ok: 0, Msg: err.Error()}), nil
	}

	err = s.userCache.SetUserOptional(ctx, userId, &user)
	if err != nil {
		return connect.NewResponse(&pb.SignUpUserResponse{Ok: 0, Msg: err.Error()}), nil
	}

	return connect.NewResponse(&pb.SignUpUserResponse{Ok: 1, UserId: userId}), nil
}

// Integrated tested
func (s *server) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error) {
	if err := setUser(ctx, req.Msg.UserId, &types.UserPtr{
		Username:           req.Msg.Name,
		Bio:                req.Msg.Bio,
		Subscribed:         req.Msg.Subscribed,
		ProfilePictureLink: req.Msg.ProfilePictureLink,
	}, false, true, s.userCache, s.db); err != nil {
		return connect.NewResponse(&pb.UpdateUserResponse{Ok: 0, Msg: err.Error()}), nil
	}

	return connect.NewResponse(&pb.UpdateUserResponse{Ok: 1}), nil
}

// Integrated tested
func (s *server) TestConnection(ctx context.Context, req *connect.Request[pb.TestConnectionRequest]) (*connect.Response[pb.TestConnectionResponse], error) {
	if req.Msg.Foo == 69.69 {
		return connect.NewResponse(&pb.TestConnectionResponse{Ok: true}), nil
	}
	return connect.NewResponse(&pb.TestConnectionResponse{Ok: false}), nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format(time.RFC1123)
		fmt.Printf("Request received at: %s, Path: %s\n", currentTime, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func StartServer(geoCache *redis.Client, postCache *redis.Client, userCache *redis.Client, db *sql.DB, uploaderS3 *s3manager.Uploader) {
	srv := &server{
		geoCache:   cache.NewGeoRedisCacheService(geoCache),
		postCache:  cache.NewPostRedisCacheService(postCache),
		userCache:  cache.NewUserRedisCacheService(userCache),
		db:         postgres.NewPostgresDatabaseService(db),
		uploaderS3: uploaderS3,
	}

	mux := http.NewServeMux()
	path, handler := pbconnect.NewDataFetcherServiceHandler(srv)

	mux.Handle(path, handler)

	loggedMux := loggingMiddleware(mux)

	http.ListenAndServe(
		":50051",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(loggedMux, &http2.Server{}),
	)
}
