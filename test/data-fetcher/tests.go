package datafetcher

import (
	"context"
	"fmt"
	"tanken/test/rpc/connectrpc/pbconnect"
	"tanken/test/rpc/pb"

	"connectrpc.com/connect"
)

func TestSignUpUser(client pbconnect.DataFetcherServiceClient) string {
	req := &pb.SignUpUserRequest{
		Email:              "alice@gotanken.com",
		Name:               "alice",
		Provider:           "github",
		Bio:                "this bio will change",
		ProfilePictureLink: "TODO: it need to be a link",
	}

	res, err := client.SignUpUser(context.Background(), connect.NewRequest(req))

	if err != nil || res.Msg.Ok != 1 {
		fmt.Errorf("error signing up user: %v", err)
	}

	return res.Msg.UserId
}

func TestGetUserInfo(client pbconnect.DataFetcherServiceClient, userId string) {
	req := &pb.GetUserInfoRequest{
		UserId: userId,
	}

	res, err := client.GetUserInfo(context.Background(), connect.NewRequest(req))

	if err != nil || res.Msg.Ok != 1 {
		fmt.Errorf("error getting user info: %v", err)
	}

	expectedResponseUser := &pb.User{
		UserId:             userId,
		UserName:           "alice",
		Bio:                "this bio will change",
		ProfilePictureLink: "TODO: it need to be a link",
		Subscribed:         0,
	}

	if res.Msg.User != expectedResponseUser {
		fmt.Errorf("expected response user: %v, got: %v", expectedResponseUser, res.Msg.User)
	}
}

func TestUpdateUser(client pbconnect.DataFetcherServiceClient, userId string) {
	updateUserReq := &pb.UpdateUserRequest{
		Bio: "this bio is changed",
	}

	updateUserRes, err := client.UpdateUser(context.Background(), connect.NewRequest(updateUserReq))

	if err != nil || updateUserRes.Msg.Ok != 1 {
		fmt.Errorf("error updating user: %v", err)
	}

	getUserReq := &pb.GetUserInfoRequest{
		UserId: userId,
	}

	getUserRes, err := client.GetUserInfo(context.Background(), connect.NewRequest(getUserReq))

	if err != nil || getUserRes.Msg.Ok != 1 {
		fmt.Errorf("error getting user info: %v", err)
	}

	expectedResponseUser := &pb.User{
		UserId:             userId,
		UserName:           "alice",
		Bio:                "this bio is changed",
		ProfilePictureLink: "TODO: it need to be a link",
		Subscribed:         0,
	}

	if getUserRes.Msg.User != expectedResponseUser {
		fmt.Errorf("expected response user: %v, got: %v", expectedResponseUser, getUserRes.Msg.User)
	}
}
