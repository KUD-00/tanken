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

	if err != nil || res.Msg.Ok != true {
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

/*
func TestAddPost(client pbconnect.DataFetcherServiceClient, userId string) (postId string, err error) {
	addPostReq := &pb.AddPostRequest{
		UserId:  userId,
		Content: "body",
		Location: &pb.Location{
			Latitude:  69.69,
			Longitude: 69.69,
		},
		Tags:         []string{"tag1", "tag2"},
		PictureChunk: nil,
	}

	addPostRes, err := client.AddPost(context.Background(), connect.NewRequest(addPostReq))

	if err != nil || addPostRes.Msg.Ok != 1 {
		return "", fmt.Errorf("error adding post: %v, response: %v", err, addPostRes.Msg.Msg)
	}

	return addPostRes.Msg.PostId, nil
}

func TestGetPostsByPostIds(client pbconnect.DataFetcherServiceClient, postIds []string, users []pb.User) error {
	getPostsByPostIdsReq := &pb.GetPostsByPostIdsRequest{
		PostIds: postIds,
	}

	getPostsByPostIdsRes, err := client.GetPostsByPostIds(context.Background(), connect.NewRequest(getPostsByPostIdsReq))

	if err != nil || getPostsByPostIdsRes.Msg.Ok != 1 {
		return fmt.Errorf("error getting posts by post ids: %v, response: %v", err, getPostsByPostIdsRes.Msg.Msg)
	}

	expectedResponsePosts := []*pb.Post{
		{
			PostId:  postIds[0],
			Author:  &users[0],
			Content: "body",
			Location: &pb.Location{
				Latitude:  69.69,
				Longitude: 69.69,
			},
			Tags:         []string{"tag1", "tag2"},
			PictureChunk: nil,
		},
	}

	if reflect.DeepEqual(getPostsByPostIdsRes.Msg.Posts[0], expectedResponsePosts) {
		return fmt.Errorf("expected response post: %v, got: %v", expectedResponsePosts, getPostsByPostIdsRes.Msg.Posts[0])
	}

	return nil
}

*/
