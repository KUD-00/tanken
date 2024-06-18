package main

import (
	"net/http"

	datafetcher "tanken/integreted-local-tests/data-fetcher"
	"tanken/integreted-local-tests/rpc/connectrpc/pbconnect"
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

	userId := datafetcher.TestSignUpUser(client)

	datafetcher.TestGetUserInfo(client, userId)

	datafetcher.TestUpdateUser(client, userId)
}
