package main

import (
	"net/http"

	datafetcher "tanken/test/data-fetcher"
	"tanken/test/rpc/connectrpc/pbconnect"
)

func main() {
	client := pbconnect.NewDataFetcherServiceClient(
		http.DefaultClient,
		"data-fetcher:50051",
	)

	userId := datafetcher.TestSignUpUser(client)

	datafetcher.TestGetUserInfo(client, userId)

	datafetcher.TestUpdateUser(client, userId)
}
