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

	userId, err := datafetcher.TestSignUpUser(client)
	if err != nil {
		panic(err)
	}

	err = datafetcher.TestGetUserInfo(client, userId)
	if err != nil {
		panic(err)
	}

	err = datafetcher.TestUpdateUser(client, userId)
	if err != nil {
		panic(err)
	}
}
