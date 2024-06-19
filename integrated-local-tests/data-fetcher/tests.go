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

func TestSignUpUser(client pbconnect.DataFetcherServiceClient) (string, error) {
	req := &pb.SignUpUserRequest{
		Email:              "alice@gotanken.com",
		Name:               "alice",
		Provider:           "github",
		Bio:                "this bio will change",
		ProfilePictureLink: "TODO: it need to be a link",
	}

	res, err := client.SignUpUser(context.Background(), connect.NewRequest(req))

	if err != nil || res.Msg.Ok != 1 {
		return "", fmt.Errorf("error signing up user: %v, response: %v", err, res.Msg.Msg)
	}

	return res.Msg.UserId, nil
}

func TestGetUserInfo(client pbconnect.DataFetcherServiceClient, userId string) error {
	req := &pb.GetUserInfoRequest{
		UserId: userId,
	}

	res, err := client.GetUserInfo(context.Background(), connect.NewRequest(req))

	if err != nil || res.Msg.Ok != 1 {
		return fmt.Errorf("error getting user info: %v, response: %v", err, res.Msg.Msg)
	}

	expectedResponseUser := &pb.User{
		UserId:             userId,
		UserName:           "alice",
		Bio:                "this bio will change",
		ProfilePictureLink: "TODO: it need to be a link",
		Subscribed:         0,
	}

	if reflect.DeepEqual(res.Msg.User, expectedResponseUser) {
		return fmt.Errorf("expected response user: %v, got: %v", expectedResponseUser, res.Msg.User)
	}

	return nil
}

func TestUpdateUser(client pbconnect.DataFetcherServiceClient, userId string) error {
	updateUserReq := &pb.UpdateUserRequest{
		UserId: userId,
		Bio:    "this bio is changed",
	}

	updateUserRes, err := client.UpdateUser(context.Background(), connect.NewRequest(updateUserReq))

	if err != nil || updateUserRes.Msg.Ok != 1 {
		return fmt.Errorf("error updating user: %v, response: %v", err, updateUserRes.Msg.Msg)
	}

	getUserReq := &pb.GetUserInfoRequest{
		UserId: userId,
	}

	getUserRes, err := client.GetUserInfo(context.Background(), connect.NewRequest(getUserReq))

	if err != nil || getUserRes.Msg.Ok != 1 {
		return fmt.Errorf("error getting user info: %v, response: %v", err, getUserRes.Msg.Msg)
	}

	expectedResponseUser := &pb.User{
		UserId:             userId,
		UserName:           "alice",
		Bio:                "this bio is changed",
		ProfilePictureLink: "TODO: it need to be a link",
		Subscribed:         0,
	}

	if reflect.DeepEqual(getUserRes.Msg.User, expectedResponseUser) {
		return fmt.Errorf("expected response user: %v, got: %v", expectedResponseUser, getUserRes.Msg.User)
	}

	return nil
}
