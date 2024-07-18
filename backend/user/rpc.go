package main

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
	"tanken/backend/user/rpc/connectrpc/pbconnect"
	pb "tanken/backend/user/rpc/pb"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type server struct {
	db         *postgres.PostgresDatabaseService
	userCache  *cache.UserRedisCacheService
	uploaderS3 *s3manager.Uploader
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

func StartServer(userCache *redis.Client, db *sql.DB, uploaderS3 *s3manager.Uploader) {
	srv := &server{
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
