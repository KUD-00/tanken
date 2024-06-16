package main

import (
	"net/http"

	datafetcher "tanken/backend/test/data-fetcher"
	"tanken/backend/test/rpc/connectrpc/pbconnect"
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
