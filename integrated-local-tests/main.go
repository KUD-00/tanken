package main

import (
	"net/http"

	datafetcher "tanken/integreted-local-tests/data-fetcher"
	"tanken/integreted-local-tests/rpc/connectrpc/pbconnect"
	"tanken/integreted-local-tests/rpc/pb"
)

func main() {
	client := pbconnect.NewDataFetcherServiceClient(
		http.DefaultClient,
		"http://data-fetcher:50051",
	)

	err := datafetcher.TestConnection(client)
	if err != nil {
		panic(err)
	}

	users := []pb.User{
		{
			Email:              "alice@gotanken.com",
			UserName:           "alice",
			Provider:           "github",
			Bio:                "this bio will change",
			ProfilePictureLink: "TODO: it need to be a link",
		},
		{
			Email:              "bob@gotanken.com",
			UserName:           "bob",
			Provider:           "github",
			Bio:                "this is bob's bio",
			ProfilePictureLink: "TODO: it need to be a link",
		},
	}

	err = datafetcher.TestSignUpUser(client, &users)
	if err != nil {
		panic(err)
	}

	err = datafetcher.TestGetUserInfo(client, &users[0])
	if err != nil {
		panic(err)
	}

	err = datafetcher.TestUpdateUser(client, &users[0])
	if err != nil {
		panic(err)
	}
}
