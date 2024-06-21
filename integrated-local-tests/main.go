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

	post := pb.Post{
		Content: "content body",
		Location: &pb.Location{
			Latitude:  34.99233214592428,
			Longitude: 135.8173205715178,
		},
		Tags:         []string{"tag1", "tag2"},
		PictureChunk: nil,
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

	err = datafetcher.TestAddPost(client, &users[0], &post)
	if err != nil {
		panic(err)
	}

	err = datafetcher.TestGetPostsByPostIds(client, &post, &users[0])
	if err != nil {
		panic(err)
	}

	err = datafetcher.TestAddLike(client, &users[1], &post)
	if err != nil {
		panic(err)
	}

	err = datafetcher.TestRemoveLike(client, &users[1], &post)
	if err != nil {
		panic(err)
	}
}
