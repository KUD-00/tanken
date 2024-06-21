package datafetcher

import (
	"context"
	"fmt"
	"reflect"
	"tanken/integreted-local-tests/rpc/connectrpc/pbconnect"
	"tanken/integreted-local-tests/rpc/pb"

	"connectrpc.com/connect"
)

func TestConnection(client pbconnect.DataFetcherServiceClient) error {
	req := &pb.TestConnectionRequest{Foo: 69.69}

	res, err := client.TestConnection(context.Background(), connect.NewRequest(req))

	if err != nil || !res.Msg.Ok {
		return fmt.Errorf("error testing connection: %v", err)
	}

	return nil
}

func TestSignUpUser(client pbconnect.DataFetcherServiceClient, users *[]pb.User) error {
	for i := range *users {
		user := &(*users)[i]
		req := &pb.SignUpUserRequest{
			Email:              user.Email,
			Name:               user.UserName,
			Provider:           user.Provider,
			Bio:                user.Bio,
			ProfilePictureLink: user.ProfilePictureLink,
		}

		res, err := client.SignUpUser(context.Background(), connect.NewRequest(req))

		if err != nil || res.Msg.Ok != 1 {
			return fmt.Errorf("error signing up user: %v, response: %v", err, res.Msg.Msg)
		}

		user.UserId = res.Msg.UserId
	}

	return nil
}

func TestGetUserInfo(client pbconnect.DataFetcherServiceClient, user *pb.User) error {
	req := &pb.GetUserInfoRequest{
		UserId: user.UserId,
	}

	res, err := client.GetUserInfo(context.Background(), connect.NewRequest(req))

	if err != nil || res.Msg.Ok != 1 {
		return fmt.Errorf("error getting user info: %v, response: %v", err, res.Msg.Msg)
	}

	expectedResponseUser := &pb.User{
		UserId:             user.UserId,
		UserName:           user.UserName,
		Bio:                user.Bio,
		ProfilePictureLink: user.ProfilePictureLink,
		Subscribed:         user.Subscribed,
	}

	if !reflect.DeepEqual(res.Msg.User, expectedResponseUser) {
		return fmt.Errorf("expected response user in TestGetUserInfo: %v, got: %v", expectedResponseUser, res.Msg.User)
	}

	return nil
}

func TestUpdateUser(client pbconnect.DataFetcherServiceClient, user *pb.User) error {
	bio := "this bio is changed"
	updateUserReq := &pb.UpdateUserRequest{
		UserId: user.UserId,
		Bio:    &bio,
	}

	updateUserRes, err := client.UpdateUser(context.Background(), connect.NewRequest(updateUserReq))

	if err != nil || updateUserRes.Msg.Ok != 1 {
		return fmt.Errorf("error updating user: %v, response: %v", err, updateUserRes.Msg.Msg)
	}

	user.Bio = "this bio is changed"

	getUserReq := &pb.GetUserInfoRequest{
		UserId: user.UserId,
	}

	getUserRes, err := client.GetUserInfo(context.Background(), connect.NewRequest(getUserReq))

	if err != nil || getUserRes.Msg.Ok != 1 {
		return fmt.Errorf("error getting user info: %v, response: %v", err, getUserRes.Msg.Msg)
	}

	expectedResponseUser := &pb.User{
		UserId:             user.UserId,
		UserName:           user.UserName,
		Bio:                user.Bio,
		ProfilePictureLink: user.ProfilePictureLink,
		Subscribed:         user.Subscribed,
	}

	if !reflect.DeepEqual(getUserRes.Msg.User, expectedResponseUser) {
		return fmt.Errorf("expected response user in TestUpdateUser: %v, got: %v", expectedResponseUser, getUserRes.Msg.User)
	}

	return nil
}

func TestAddPost(client pbconnect.DataFetcherServiceClient, user *pb.User, post *pb.Post) error {
	addPostReq := &pb.AddPostRequest{
		UserId:       user.UserId,
		Content:      post.Content,
		Location:     post.Location,
		Tags:         post.Tags,
		PictureChunk: post.PictureChunk,
	}

	addPostRes, err := client.AddPost(context.Background(), connect.NewRequest(addPostReq))

	if err != nil || addPostRes.Msg.Ok != 1 {
		return fmt.Errorf("error adding post: %v, response: %v", err, addPostRes.Msg.Msg)
	}

	post.PostId = addPostRes.Msg.PostId

	return nil
}

func TestGetPostsByPostIds(client pbconnect.DataFetcherServiceClient, post *pb.Post, user *pb.User) error {
	getPostsByPostIdsReq := &pb.GetPostsByPostIdsRequest{
		PostIds: []string{post.PostId},
	}

	getPostsByPostIdsRes, err := client.GetPostsByPostIds(context.Background(), connect.NewRequest(getPostsByPostIdsReq))

	if err != nil || getPostsByPostIdsRes.Msg.Ok != 1 {
		return fmt.Errorf("error getting posts by post ids: %v, response: %v", err, getPostsByPostIdsRes.Msg.Msg)
	}

	expectedResponsePosts := []*pb.Post{
		{
			PostId:       post.PostId,
			Author:       user,
			Content:      post.Content,
			Location:     post.Location,
			Tags:         post.Tags,
			PictureChunk: post.PictureChunk,
		},
	}

	if reflect.DeepEqual(getPostsByPostIdsRes.Msg.Posts[0], expectedResponsePosts) {
		return fmt.Errorf("expected response post: %v, got: %v", expectedResponsePosts, getPostsByPostIdsRes.Msg.Posts[0])
	}

	return nil
}

func TestAddLike(client pbconnect.DataFetcherServiceClient, user *pb.User, post *pb.Post) error {
	addLikeReq := &pb.AddLikeRequest{
		UserId: user.UserId,
		PostId: post.PostId,
	}

	addLikeRes, err := client.AddLike(context.Background(), connect.NewRequest(addLikeReq))

	if err != nil || addLikeRes.Msg.Ok != 1 {
		return fmt.Errorf("error adding like: %v, response: %v", err, addLikeRes.Msg.Msg)
	}

	postRes, err := client.GetPostsByPostIds(context.Background(), connect.NewRequest(&pb.GetPostsByPostIdsRequest{PostIds: []string{post.PostId}}))

	if err != nil || postRes.Msg.Ok != 1 {
		return fmt.Errorf("error getting posts by post ids: %v, response: %v", err, postRes.Msg.Msg)
	}

	if postRes.Msg.Posts[0].Likes != 1 {
		return fmt.Errorf("expected post likes to be 1, got: %v", postRes.Msg.Posts[0].Likes)
	}
	// TODO: won't check likedby cause rpc won't return it

	return nil
}

func TestRemoveLike(client pbconnect.DataFetcherServiceClient, user *pb.User, post *pb.Post) error {
	removeLikeReq := &pb.RemoveLikeRequest{
		UserId: user.UserId,
		PostId: post.PostId,
	}

	removeLikeRes, err := client.RemoveLike(context.Background(), connect.NewRequest(removeLikeReq))

	if err != nil || removeLikeRes.Msg.Ok != 1 {
		return fmt.Errorf("error removing like: %v, response: %v", err, removeLikeRes.Msg.Msg)
	}

	postRes, err := client.GetPostsByPostIds(context.Background(), connect.NewRequest(&pb.GetPostsByPostIdsRequest{PostIds: []string{post.PostId}}))

	if err != nil || postRes.Msg.Ok != 1 {
		return fmt.Errorf("error getting posts by post ids: %v, response: %v", err, postRes.Msg.Msg)
	}

	if postRes.Msg.Posts[0].Likes != 0 {
		return fmt.Errorf("expected post likes to be 0, got: %v", postRes.Msg.Posts[0].Likes)
	}
	// TODO: won't check likedby cause rpc won't return it

	return nil
}

func TestAddComment(client pbconnect.DataFetcherServiceClient, user *pb.User, post *pb.Post, comment *pb.Comment) error {
	addCommentReq := &pb.AddCommentRequest{
		UserId:  user.UserId,
		PostId:  post.PostId,
		Content: comment.Content,
	}

	addCommentRes, err := client.AddComment(context.Background(), connect.NewRequest(addCommentReq))

	if err != nil || addCommentRes.Msg.Ok != 1 {
		return fmt.Errorf("error adding comment: %v, response: %v", err, addCommentRes.Msg.Msg)
	}

	return nil
}
