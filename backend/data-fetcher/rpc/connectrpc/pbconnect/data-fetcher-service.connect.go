// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: data-fetcher-service.proto

package pbconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	http "net/http"
	strings "strings"
	pb "tanken/backend/data-fetcher/rpc/pb"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// DataFetcherServiceName is the fully-qualified name of the DataFetcherService service.
	DataFetcherServiceName = "rpc.DataFetcherService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// DataFetcherServiceTestConnectionProcedure is the fully-qualified name of the DataFetcherService's
	// TestConnection RPC.
	DataFetcherServiceTestConnectionProcedure = "/rpc.DataFetcherService/TestConnection"
	// DataFetcherServiceGetPostsByLocationProcedure is the fully-qualified name of the
	// DataFetcherService's GetPostsByLocation RPC.
	DataFetcherServiceGetPostsByLocationProcedure = "/rpc.DataFetcherService/GetPostsByLocation"
	// DataFetcherServiceGetPostsByPostIdsProcedure is the fully-qualified name of the
	// DataFetcherService's GetPostsByPostIds RPC.
	DataFetcherServiceGetPostsByPostIdsProcedure = "/rpc.DataFetcherService/GetPostsByPostIds"
	// DataFetcherServiceGetPostsByUserProcedure is the fully-qualified name of the DataFetcherService's
	// GetPostsByUser RPC.
	DataFetcherServiceGetPostsByUserProcedure = "/rpc.DataFetcherService/GetPostsByUser"
	// DataFetcherServiceAddPostProcedure is the fully-qualified name of the DataFetcherService's
	// AddPost RPC.
	DataFetcherServiceAddPostProcedure = "/rpc.DataFetcherService/AddPost"
	// DataFetcherServiceHardDeletePostProcedure is the fully-qualified name of the DataFetcherService's
	// HardDeletePost RPC.
	DataFetcherServiceHardDeletePostProcedure = "/rpc.DataFetcherService/HardDeletePost"
	// DataFetcherServiceSoftDeletePostProcedure is the fully-qualified name of the DataFetcherService's
	// SoftDeletePost RPC.
	DataFetcherServiceSoftDeletePostProcedure = "/rpc.DataFetcherService/SoftDeletePost"
	// DataFetcherServiceAddLikeProcedure is the fully-qualified name of the DataFetcherService's
	// AddLike RPC.
	DataFetcherServiceAddLikeProcedure = "/rpc.DataFetcherService/AddLike"
	// DataFetcherServiceRemoveLikeProcedure is the fully-qualified name of the DataFetcherService's
	// RemoveLike RPC.
	DataFetcherServiceRemoveLikeProcedure = "/rpc.DataFetcherService/RemoveLike"
	// DataFetcherServiceAddBookmarkProcedure is the fully-qualified name of the DataFetcherService's
	// AddBookmark RPC.
	DataFetcherServiceAddBookmarkProcedure = "/rpc.DataFetcherService/AddBookmark"
	// DataFetcherServiceRemoveBookmarkProcedure is the fully-qualified name of the DataFetcherService's
	// RemoveBookmark RPC.
	DataFetcherServiceRemoveBookmarkProcedure = "/rpc.DataFetcherService/RemoveBookmark"
	// DataFetcherServiceAddCommentProcedure is the fully-qualified name of the DataFetcherService's
	// AddComment RPC.
	DataFetcherServiceAddCommentProcedure = "/rpc.DataFetcherService/AddComment"
	// DataFetcherServiceRemoveCommentProcedure is the fully-qualified name of the DataFetcherService's
	// RemoveComment RPC.
	DataFetcherServiceRemoveCommentProcedure = "/rpc.DataFetcherService/RemoveComment"
	// DataFetcherServiceGetUserInfoProcedure is the fully-qualified name of the DataFetcherService's
	// GetUserInfo RPC.
	DataFetcherServiceGetUserInfoProcedure = "/rpc.DataFetcherService/GetUserInfo"
	// DataFetcherServiceSignUpUserProcedure is the fully-qualified name of the DataFetcherService's
	// SignUpUser RPC.
	DataFetcherServiceSignUpUserProcedure = "/rpc.DataFetcherService/SignUpUser"
	// DataFetcherServiceUpdateUserProcedure is the fully-qualified name of the DataFetcherService's
	// UpdateUser RPC.
	DataFetcherServiceUpdateUserProcedure = "/rpc.DataFetcherService/UpdateUser"
	// DataFetcherServiceHardDeleteUserProcedure is the fully-qualified name of the DataFetcherService's
	// HardDeleteUser RPC.
	DataFetcherServiceHardDeleteUserProcedure = "/rpc.DataFetcherService/HardDeleteUser"
	// DataFetcherServiceSoftDeleteUserProcedure is the fully-qualified name of the DataFetcherService's
	// SoftDeleteUser RPC.
	DataFetcherServiceSoftDeleteUserProcedure = "/rpc.DataFetcherService/SoftDeleteUser"
	// DataFetcherServiceGetUserInfoByOAuthProcedure is the fully-qualified name of the
	// DataFetcherService's GetUserInfoByOAuth RPC.
	DataFetcherServiceGetUserInfoByOAuthProcedure = "/rpc.DataFetcherService/GetUserInfoByOAuth"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	dataFetcherServiceServiceDescriptor                  = pb.File_data_fetcher_service_proto.Services().ByName("DataFetcherService")
	dataFetcherServiceTestConnectionMethodDescriptor     = dataFetcherServiceServiceDescriptor.Methods().ByName("TestConnection")
	dataFetcherServiceGetPostsByLocationMethodDescriptor = dataFetcherServiceServiceDescriptor.Methods().ByName("GetPostsByLocation")
	dataFetcherServiceGetPostsByPostIdsMethodDescriptor  = dataFetcherServiceServiceDescriptor.Methods().ByName("GetPostsByPostIds")
	dataFetcherServiceGetPostsByUserMethodDescriptor     = dataFetcherServiceServiceDescriptor.Methods().ByName("GetPostsByUser")
	dataFetcherServiceAddPostMethodDescriptor            = dataFetcherServiceServiceDescriptor.Methods().ByName("AddPost")
	dataFetcherServiceHardDeletePostMethodDescriptor     = dataFetcherServiceServiceDescriptor.Methods().ByName("HardDeletePost")
	dataFetcherServiceSoftDeletePostMethodDescriptor     = dataFetcherServiceServiceDescriptor.Methods().ByName("SoftDeletePost")
	dataFetcherServiceAddLikeMethodDescriptor            = dataFetcherServiceServiceDescriptor.Methods().ByName("AddLike")
	dataFetcherServiceRemoveLikeMethodDescriptor         = dataFetcherServiceServiceDescriptor.Methods().ByName("RemoveLike")
	dataFetcherServiceAddBookmarkMethodDescriptor        = dataFetcherServiceServiceDescriptor.Methods().ByName("AddBookmark")
	dataFetcherServiceRemoveBookmarkMethodDescriptor     = dataFetcherServiceServiceDescriptor.Methods().ByName("RemoveBookmark")
	dataFetcherServiceAddCommentMethodDescriptor         = dataFetcherServiceServiceDescriptor.Methods().ByName("AddComment")
	dataFetcherServiceRemoveCommentMethodDescriptor      = dataFetcherServiceServiceDescriptor.Methods().ByName("RemoveComment")
	dataFetcherServiceGetUserInfoMethodDescriptor        = dataFetcherServiceServiceDescriptor.Methods().ByName("GetUserInfo")
	dataFetcherServiceSignUpUserMethodDescriptor         = dataFetcherServiceServiceDescriptor.Methods().ByName("SignUpUser")
	dataFetcherServiceUpdateUserMethodDescriptor         = dataFetcherServiceServiceDescriptor.Methods().ByName("UpdateUser")
	dataFetcherServiceHardDeleteUserMethodDescriptor     = dataFetcherServiceServiceDescriptor.Methods().ByName("HardDeleteUser")
	dataFetcherServiceSoftDeleteUserMethodDescriptor     = dataFetcherServiceServiceDescriptor.Methods().ByName("SoftDeleteUser")
	dataFetcherServiceGetUserInfoByOAuthMethodDescriptor = dataFetcherServiceServiceDescriptor.Methods().ByName("GetUserInfoByOAuth")
)

// DataFetcherServiceClient is a client for the rpc.DataFetcherService service.
type DataFetcherServiceClient interface {
	TestConnection(context.Context, *connect.Request[pb.TestConnectionRequest]) (*connect.Response[pb.TestConnectionResponse], error)
	GetPostsByLocation(context.Context, *connect.Request[pb.GetPostsByLocationRequest]) (*connect.Response[pb.GetPostsByLocationResponse], error)
	GetPostsByPostIds(context.Context, *connect.Request[pb.GetPostsByPostIdsRequest]) (*connect.Response[pb.GetPostsByPostIdsResponse], error)
	GetPostsByUser(context.Context, *connect.Request[pb.GetPostsByUserIdRequest]) (*connect.Response[pb.GetPostsByUserIdResponse], error)
	AddPost(context.Context, *connect.Request[pb.AddPostRequest]) (*connect.Response[pb.AddPostResponse], error)
	HardDeletePost(context.Context, *connect.Request[pb.HardDeletePostRequest]) (*connect.Response[pb.HardDeletePostResponse], error)
	SoftDeletePost(context.Context, *connect.Request[pb.SoftDeletePostRequest]) (*connect.Response[pb.SoftDeletePostResponse], error)
	AddLike(context.Context, *connect.Request[pb.AddLikeRequest]) (*connect.Response[pb.AddLikeResponse], error)
	RemoveLike(context.Context, *connect.Request[pb.RemoveLikeRequest]) (*connect.Response[pb.RemoveLikeResponse], error)
	AddBookmark(context.Context, *connect.Request[pb.AddBookmarkRequest]) (*connect.Response[pb.AddBookmarkResponse], error)
	RemoveBookmark(context.Context, *connect.Request[pb.RemoveBookmarkRequest]) (*connect.Response[pb.RemoveBookmarkResponse], error)
	AddComment(context.Context, *connect.Request[pb.AddCommentRequest]) (*connect.Response[pb.AddCommentResponse], error)
	RemoveComment(context.Context, *connect.Request[pb.RemoveCommentRequest]) (*connect.Response[pb.RemoveCommentResponse], error)
	GetUserInfo(context.Context, *connect.Request[pb.GetUserInfoRequest]) (*connect.Response[pb.GetUserInfoResponse], error)
	SignUpUser(context.Context, *connect.Request[pb.SignUpUserRequest]) (*connect.Response[pb.SignUpUserResponse], error)
	UpdateUser(context.Context, *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error)
	HardDeleteUser(context.Context, *connect.Request[pb.HardDeleteUserRequest]) (*connect.Response[pb.HardDeleteUserResponse], error)
	SoftDeleteUser(context.Context, *connect.Request[pb.SoftDeleteUserRequest]) (*connect.Response[pb.SoftDeleteUserResponse], error)
	GetUserInfoByOAuth(context.Context, *connect.Request[pb.GetUserInfoByOAuthRequest]) (*connect.Response[pb.GetUserInfoByOAuthResponse], error)
}

// NewDataFetcherServiceClient constructs a client for the rpc.DataFetcherService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewDataFetcherServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) DataFetcherServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &dataFetcherServiceClient{
		testConnection: connect.NewClient[pb.TestConnectionRequest, pb.TestConnectionResponse](
			httpClient,
			baseURL+DataFetcherServiceTestConnectionProcedure,
			connect.WithSchema(dataFetcherServiceTestConnectionMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getPostsByLocation: connect.NewClient[pb.GetPostsByLocationRequest, pb.GetPostsByLocationResponse](
			httpClient,
			baseURL+DataFetcherServiceGetPostsByLocationProcedure,
			connect.WithSchema(dataFetcherServiceGetPostsByLocationMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getPostsByPostIds: connect.NewClient[pb.GetPostsByPostIdsRequest, pb.GetPostsByPostIdsResponse](
			httpClient,
			baseURL+DataFetcherServiceGetPostsByPostIdsProcedure,
			connect.WithSchema(dataFetcherServiceGetPostsByPostIdsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getPostsByUser: connect.NewClient[pb.GetPostsByUserIdRequest, pb.GetPostsByUserIdResponse](
			httpClient,
			baseURL+DataFetcherServiceGetPostsByUserProcedure,
			connect.WithSchema(dataFetcherServiceGetPostsByUserMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		addPost: connect.NewClient[pb.AddPostRequest, pb.AddPostResponse](
			httpClient,
			baseURL+DataFetcherServiceAddPostProcedure,
			connect.WithSchema(dataFetcherServiceAddPostMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		hardDeletePost: connect.NewClient[pb.HardDeletePostRequest, pb.HardDeletePostResponse](
			httpClient,
			baseURL+DataFetcherServiceHardDeletePostProcedure,
			connect.WithSchema(dataFetcherServiceHardDeletePostMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		softDeletePost: connect.NewClient[pb.SoftDeletePostRequest, pb.SoftDeletePostResponse](
			httpClient,
			baseURL+DataFetcherServiceSoftDeletePostProcedure,
			connect.WithSchema(dataFetcherServiceSoftDeletePostMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		addLike: connect.NewClient[pb.AddLikeRequest, pb.AddLikeResponse](
			httpClient,
			baseURL+DataFetcherServiceAddLikeProcedure,
			connect.WithSchema(dataFetcherServiceAddLikeMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		removeLike: connect.NewClient[pb.RemoveLikeRequest, pb.RemoveLikeResponse](
			httpClient,
			baseURL+DataFetcherServiceRemoveLikeProcedure,
			connect.WithSchema(dataFetcherServiceRemoveLikeMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		addBookmark: connect.NewClient[pb.AddBookmarkRequest, pb.AddBookmarkResponse](
			httpClient,
			baseURL+DataFetcherServiceAddBookmarkProcedure,
			connect.WithSchema(dataFetcherServiceAddBookmarkMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		removeBookmark: connect.NewClient[pb.RemoveBookmarkRequest, pb.RemoveBookmarkResponse](
			httpClient,
			baseURL+DataFetcherServiceRemoveBookmarkProcedure,
			connect.WithSchema(dataFetcherServiceRemoveBookmarkMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		addComment: connect.NewClient[pb.AddCommentRequest, pb.AddCommentResponse](
			httpClient,
			baseURL+DataFetcherServiceAddCommentProcedure,
			connect.WithSchema(dataFetcherServiceAddCommentMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		removeComment: connect.NewClient[pb.RemoveCommentRequest, pb.RemoveCommentResponse](
			httpClient,
			baseURL+DataFetcherServiceRemoveCommentProcedure,
			connect.WithSchema(dataFetcherServiceRemoveCommentMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getUserInfo: connect.NewClient[pb.GetUserInfoRequest, pb.GetUserInfoResponse](
			httpClient,
			baseURL+DataFetcherServiceGetUserInfoProcedure,
			connect.WithSchema(dataFetcherServiceGetUserInfoMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		signUpUser: connect.NewClient[pb.SignUpUserRequest, pb.SignUpUserResponse](
			httpClient,
			baseURL+DataFetcherServiceSignUpUserProcedure,
			connect.WithSchema(dataFetcherServiceSignUpUserMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updateUser: connect.NewClient[pb.UpdateUserRequest, pb.UpdateUserResponse](
			httpClient,
			baseURL+DataFetcherServiceUpdateUserProcedure,
			connect.WithSchema(dataFetcherServiceUpdateUserMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		hardDeleteUser: connect.NewClient[pb.HardDeleteUserRequest, pb.HardDeleteUserResponse](
			httpClient,
			baseURL+DataFetcherServiceHardDeleteUserProcedure,
			connect.WithSchema(dataFetcherServiceHardDeleteUserMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		softDeleteUser: connect.NewClient[pb.SoftDeleteUserRequest, pb.SoftDeleteUserResponse](
			httpClient,
			baseURL+DataFetcherServiceSoftDeleteUserProcedure,
			connect.WithSchema(dataFetcherServiceSoftDeleteUserMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getUserInfoByOAuth: connect.NewClient[pb.GetUserInfoByOAuthRequest, pb.GetUserInfoByOAuthResponse](
			httpClient,
			baseURL+DataFetcherServiceGetUserInfoByOAuthProcedure,
			connect.WithSchema(dataFetcherServiceGetUserInfoByOAuthMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// dataFetcherServiceClient implements DataFetcherServiceClient.
type dataFetcherServiceClient struct {
	testConnection     *connect.Client[pb.TestConnectionRequest, pb.TestConnectionResponse]
	getPostsByLocation *connect.Client[pb.GetPostsByLocationRequest, pb.GetPostsByLocationResponse]
	getPostsByPostIds  *connect.Client[pb.GetPostsByPostIdsRequest, pb.GetPostsByPostIdsResponse]
	getPostsByUser     *connect.Client[pb.GetPostsByUserIdRequest, pb.GetPostsByUserIdResponse]
	addPost            *connect.Client[pb.AddPostRequest, pb.AddPostResponse]
	hardDeletePost     *connect.Client[pb.HardDeletePostRequest, pb.HardDeletePostResponse]
	softDeletePost     *connect.Client[pb.SoftDeletePostRequest, pb.SoftDeletePostResponse]
	addLike            *connect.Client[pb.AddLikeRequest, pb.AddLikeResponse]
	removeLike         *connect.Client[pb.RemoveLikeRequest, pb.RemoveLikeResponse]
	addBookmark        *connect.Client[pb.AddBookmarkRequest, pb.AddBookmarkResponse]
	removeBookmark     *connect.Client[pb.RemoveBookmarkRequest, pb.RemoveBookmarkResponse]
	addComment         *connect.Client[pb.AddCommentRequest, pb.AddCommentResponse]
	removeComment      *connect.Client[pb.RemoveCommentRequest, pb.RemoveCommentResponse]
	getUserInfo        *connect.Client[pb.GetUserInfoRequest, pb.GetUserInfoResponse]
	signUpUser         *connect.Client[pb.SignUpUserRequest, pb.SignUpUserResponse]
	updateUser         *connect.Client[pb.UpdateUserRequest, pb.UpdateUserResponse]
	hardDeleteUser     *connect.Client[pb.HardDeleteUserRequest, pb.HardDeleteUserResponse]
	softDeleteUser     *connect.Client[pb.SoftDeleteUserRequest, pb.SoftDeleteUserResponse]
	getUserInfoByOAuth *connect.Client[pb.GetUserInfoByOAuthRequest, pb.GetUserInfoByOAuthResponse]
}

// TestConnection calls rpc.DataFetcherService.TestConnection.
func (c *dataFetcherServiceClient) TestConnection(ctx context.Context, req *connect.Request[pb.TestConnectionRequest]) (*connect.Response[pb.TestConnectionResponse], error) {
	return c.testConnection.CallUnary(ctx, req)
}

// GetPostsByLocation calls rpc.DataFetcherService.GetPostsByLocation.
func (c *dataFetcherServiceClient) GetPostsByLocation(ctx context.Context, req *connect.Request[pb.GetPostsByLocationRequest]) (*connect.Response[pb.GetPostsByLocationResponse], error) {
	return c.getPostsByLocation.CallUnary(ctx, req)
}

// GetPostsByPostIds calls rpc.DataFetcherService.GetPostsByPostIds.
func (c *dataFetcherServiceClient) GetPostsByPostIds(ctx context.Context, req *connect.Request[pb.GetPostsByPostIdsRequest]) (*connect.Response[pb.GetPostsByPostIdsResponse], error) {
	return c.getPostsByPostIds.CallUnary(ctx, req)
}

// GetPostsByUser calls rpc.DataFetcherService.GetPostsByUser.
func (c *dataFetcherServiceClient) GetPostsByUser(ctx context.Context, req *connect.Request[pb.GetPostsByUserIdRequest]) (*connect.Response[pb.GetPostsByUserIdResponse], error) {
	return c.getPostsByUser.CallUnary(ctx, req)
}

// AddPost calls rpc.DataFetcherService.AddPost.
func (c *dataFetcherServiceClient) AddPost(ctx context.Context, req *connect.Request[pb.AddPostRequest]) (*connect.Response[pb.AddPostResponse], error) {
	return c.addPost.CallUnary(ctx, req)
}

// HardDeletePost calls rpc.DataFetcherService.HardDeletePost.
func (c *dataFetcherServiceClient) HardDeletePost(ctx context.Context, req *connect.Request[pb.HardDeletePostRequest]) (*connect.Response[pb.HardDeletePostResponse], error) {
	return c.hardDeletePost.CallUnary(ctx, req)
}

// SoftDeletePost calls rpc.DataFetcherService.SoftDeletePost.
func (c *dataFetcherServiceClient) SoftDeletePost(ctx context.Context, req *connect.Request[pb.SoftDeletePostRequest]) (*connect.Response[pb.SoftDeletePostResponse], error) {
	return c.softDeletePost.CallUnary(ctx, req)
}

// AddLike calls rpc.DataFetcherService.AddLike.
func (c *dataFetcherServiceClient) AddLike(ctx context.Context, req *connect.Request[pb.AddLikeRequest]) (*connect.Response[pb.AddLikeResponse], error) {
	return c.addLike.CallUnary(ctx, req)
}

// RemoveLike calls rpc.DataFetcherService.RemoveLike.
func (c *dataFetcherServiceClient) RemoveLike(ctx context.Context, req *connect.Request[pb.RemoveLikeRequest]) (*connect.Response[pb.RemoveLikeResponse], error) {
	return c.removeLike.CallUnary(ctx, req)
}

// AddBookmark calls rpc.DataFetcherService.AddBookmark.
func (c *dataFetcherServiceClient) AddBookmark(ctx context.Context, req *connect.Request[pb.AddBookmarkRequest]) (*connect.Response[pb.AddBookmarkResponse], error) {
	return c.addBookmark.CallUnary(ctx, req)
}

// RemoveBookmark calls rpc.DataFetcherService.RemoveBookmark.
func (c *dataFetcherServiceClient) RemoveBookmark(ctx context.Context, req *connect.Request[pb.RemoveBookmarkRequest]) (*connect.Response[pb.RemoveBookmarkResponse], error) {
	return c.removeBookmark.CallUnary(ctx, req)
}

// AddComment calls rpc.DataFetcherService.AddComment.
func (c *dataFetcherServiceClient) AddComment(ctx context.Context, req *connect.Request[pb.AddCommentRequest]) (*connect.Response[pb.AddCommentResponse], error) {
	return c.addComment.CallUnary(ctx, req)
}

// RemoveComment calls rpc.DataFetcherService.RemoveComment.
func (c *dataFetcherServiceClient) RemoveComment(ctx context.Context, req *connect.Request[pb.RemoveCommentRequest]) (*connect.Response[pb.RemoveCommentResponse], error) {
	return c.removeComment.CallUnary(ctx, req)
}

// GetUserInfo calls rpc.DataFetcherService.GetUserInfo.
func (c *dataFetcherServiceClient) GetUserInfo(ctx context.Context, req *connect.Request[pb.GetUserInfoRequest]) (*connect.Response[pb.GetUserInfoResponse], error) {
	return c.getUserInfo.CallUnary(ctx, req)
}

// SignUpUser calls rpc.DataFetcherService.SignUpUser.
func (c *dataFetcherServiceClient) SignUpUser(ctx context.Context, req *connect.Request[pb.SignUpUserRequest]) (*connect.Response[pb.SignUpUserResponse], error) {
	return c.signUpUser.CallUnary(ctx, req)
}

// UpdateUser calls rpc.DataFetcherService.UpdateUser.
func (c *dataFetcherServiceClient) UpdateUser(ctx context.Context, req *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error) {
	return c.updateUser.CallUnary(ctx, req)
}

// HardDeleteUser calls rpc.DataFetcherService.HardDeleteUser.
func (c *dataFetcherServiceClient) HardDeleteUser(ctx context.Context, req *connect.Request[pb.HardDeleteUserRequest]) (*connect.Response[pb.HardDeleteUserResponse], error) {
	return c.hardDeleteUser.CallUnary(ctx, req)
}

// SoftDeleteUser calls rpc.DataFetcherService.SoftDeleteUser.
func (c *dataFetcherServiceClient) SoftDeleteUser(ctx context.Context, req *connect.Request[pb.SoftDeleteUserRequest]) (*connect.Response[pb.SoftDeleteUserResponse], error) {
	return c.softDeleteUser.CallUnary(ctx, req)
}

// GetUserInfoByOAuth calls rpc.DataFetcherService.GetUserInfoByOAuth.
func (c *dataFetcherServiceClient) GetUserInfoByOAuth(ctx context.Context, req *connect.Request[pb.GetUserInfoByOAuthRequest]) (*connect.Response[pb.GetUserInfoByOAuthResponse], error) {
	return c.getUserInfoByOAuth.CallUnary(ctx, req)
}

// DataFetcherServiceHandler is an implementation of the rpc.DataFetcherService service.
type DataFetcherServiceHandler interface {
	TestConnection(context.Context, *connect.Request[pb.TestConnectionRequest]) (*connect.Response[pb.TestConnectionResponse], error)
	GetPostsByLocation(context.Context, *connect.Request[pb.GetPostsByLocationRequest]) (*connect.Response[pb.GetPostsByLocationResponse], error)
	GetPostsByPostIds(context.Context, *connect.Request[pb.GetPostsByPostIdsRequest]) (*connect.Response[pb.GetPostsByPostIdsResponse], error)
	GetPostsByUser(context.Context, *connect.Request[pb.GetPostsByUserIdRequest]) (*connect.Response[pb.GetPostsByUserIdResponse], error)
	AddPost(context.Context, *connect.Request[pb.AddPostRequest]) (*connect.Response[pb.AddPostResponse], error)
	HardDeletePost(context.Context, *connect.Request[pb.HardDeletePostRequest]) (*connect.Response[pb.HardDeletePostResponse], error)
	SoftDeletePost(context.Context, *connect.Request[pb.SoftDeletePostRequest]) (*connect.Response[pb.SoftDeletePostResponse], error)
	AddLike(context.Context, *connect.Request[pb.AddLikeRequest]) (*connect.Response[pb.AddLikeResponse], error)
	RemoveLike(context.Context, *connect.Request[pb.RemoveLikeRequest]) (*connect.Response[pb.RemoveLikeResponse], error)
	AddBookmark(context.Context, *connect.Request[pb.AddBookmarkRequest]) (*connect.Response[pb.AddBookmarkResponse], error)
	RemoveBookmark(context.Context, *connect.Request[pb.RemoveBookmarkRequest]) (*connect.Response[pb.RemoveBookmarkResponse], error)
	AddComment(context.Context, *connect.Request[pb.AddCommentRequest]) (*connect.Response[pb.AddCommentResponse], error)
	RemoveComment(context.Context, *connect.Request[pb.RemoveCommentRequest]) (*connect.Response[pb.RemoveCommentResponse], error)
	GetUserInfo(context.Context, *connect.Request[pb.GetUserInfoRequest]) (*connect.Response[pb.GetUserInfoResponse], error)
	SignUpUser(context.Context, *connect.Request[pb.SignUpUserRequest]) (*connect.Response[pb.SignUpUserResponse], error)
	UpdateUser(context.Context, *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error)
	HardDeleteUser(context.Context, *connect.Request[pb.HardDeleteUserRequest]) (*connect.Response[pb.HardDeleteUserResponse], error)
	SoftDeleteUser(context.Context, *connect.Request[pb.SoftDeleteUserRequest]) (*connect.Response[pb.SoftDeleteUserResponse], error)
	GetUserInfoByOAuth(context.Context, *connect.Request[pb.GetUserInfoByOAuthRequest]) (*connect.Response[pb.GetUserInfoByOAuthResponse], error)
}

// NewDataFetcherServiceHandler builds an HTTP handler from the service implementation. It returns
// the path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewDataFetcherServiceHandler(svc DataFetcherServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	dataFetcherServiceTestConnectionHandler := connect.NewUnaryHandler(
		DataFetcherServiceTestConnectionProcedure,
		svc.TestConnection,
		connect.WithSchema(dataFetcherServiceTestConnectionMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceGetPostsByLocationHandler := connect.NewUnaryHandler(
		DataFetcherServiceGetPostsByLocationProcedure,
		svc.GetPostsByLocation,
		connect.WithSchema(dataFetcherServiceGetPostsByLocationMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceGetPostsByPostIdsHandler := connect.NewUnaryHandler(
		DataFetcherServiceGetPostsByPostIdsProcedure,
		svc.GetPostsByPostIds,
		connect.WithSchema(dataFetcherServiceGetPostsByPostIdsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceGetPostsByUserHandler := connect.NewUnaryHandler(
		DataFetcherServiceGetPostsByUserProcedure,
		svc.GetPostsByUser,
		connect.WithSchema(dataFetcherServiceGetPostsByUserMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceAddPostHandler := connect.NewUnaryHandler(
		DataFetcherServiceAddPostProcedure,
		svc.AddPost,
		connect.WithSchema(dataFetcherServiceAddPostMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceHardDeletePostHandler := connect.NewUnaryHandler(
		DataFetcherServiceHardDeletePostProcedure,
		svc.HardDeletePost,
		connect.WithSchema(dataFetcherServiceHardDeletePostMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceSoftDeletePostHandler := connect.NewUnaryHandler(
		DataFetcherServiceSoftDeletePostProcedure,
		svc.SoftDeletePost,
		connect.WithSchema(dataFetcherServiceSoftDeletePostMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceAddLikeHandler := connect.NewUnaryHandler(
		DataFetcherServiceAddLikeProcedure,
		svc.AddLike,
		connect.WithSchema(dataFetcherServiceAddLikeMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceRemoveLikeHandler := connect.NewUnaryHandler(
		DataFetcherServiceRemoveLikeProcedure,
		svc.RemoveLike,
		connect.WithSchema(dataFetcherServiceRemoveLikeMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceAddBookmarkHandler := connect.NewUnaryHandler(
		DataFetcherServiceAddBookmarkProcedure,
		svc.AddBookmark,
		connect.WithSchema(dataFetcherServiceAddBookmarkMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceRemoveBookmarkHandler := connect.NewUnaryHandler(
		DataFetcherServiceRemoveBookmarkProcedure,
		svc.RemoveBookmark,
		connect.WithSchema(dataFetcherServiceRemoveBookmarkMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceAddCommentHandler := connect.NewUnaryHandler(
		DataFetcherServiceAddCommentProcedure,
		svc.AddComment,
		connect.WithSchema(dataFetcherServiceAddCommentMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceRemoveCommentHandler := connect.NewUnaryHandler(
		DataFetcherServiceRemoveCommentProcedure,
		svc.RemoveComment,
		connect.WithSchema(dataFetcherServiceRemoveCommentMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceGetUserInfoHandler := connect.NewUnaryHandler(
		DataFetcherServiceGetUserInfoProcedure,
		svc.GetUserInfo,
		connect.WithSchema(dataFetcherServiceGetUserInfoMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceSignUpUserHandler := connect.NewUnaryHandler(
		DataFetcherServiceSignUpUserProcedure,
		svc.SignUpUser,
		connect.WithSchema(dataFetcherServiceSignUpUserMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceUpdateUserHandler := connect.NewUnaryHandler(
		DataFetcherServiceUpdateUserProcedure,
		svc.UpdateUser,
		connect.WithSchema(dataFetcherServiceUpdateUserMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceHardDeleteUserHandler := connect.NewUnaryHandler(
		DataFetcherServiceHardDeleteUserProcedure,
		svc.HardDeleteUser,
		connect.WithSchema(dataFetcherServiceHardDeleteUserMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceSoftDeleteUserHandler := connect.NewUnaryHandler(
		DataFetcherServiceSoftDeleteUserProcedure,
		svc.SoftDeleteUser,
		connect.WithSchema(dataFetcherServiceSoftDeleteUserMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	dataFetcherServiceGetUserInfoByOAuthHandler := connect.NewUnaryHandler(
		DataFetcherServiceGetUserInfoByOAuthProcedure,
		svc.GetUserInfoByOAuth,
		connect.WithSchema(dataFetcherServiceGetUserInfoByOAuthMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/rpc.DataFetcherService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case DataFetcherServiceTestConnectionProcedure:
			dataFetcherServiceTestConnectionHandler.ServeHTTP(w, r)
		case DataFetcherServiceGetPostsByLocationProcedure:
			dataFetcherServiceGetPostsByLocationHandler.ServeHTTP(w, r)
		case DataFetcherServiceGetPostsByPostIdsProcedure:
			dataFetcherServiceGetPostsByPostIdsHandler.ServeHTTP(w, r)
		case DataFetcherServiceGetPostsByUserProcedure:
			dataFetcherServiceGetPostsByUserHandler.ServeHTTP(w, r)
		case DataFetcherServiceAddPostProcedure:
			dataFetcherServiceAddPostHandler.ServeHTTP(w, r)
		case DataFetcherServiceHardDeletePostProcedure:
			dataFetcherServiceHardDeletePostHandler.ServeHTTP(w, r)
		case DataFetcherServiceSoftDeletePostProcedure:
			dataFetcherServiceSoftDeletePostHandler.ServeHTTP(w, r)
		case DataFetcherServiceAddLikeProcedure:
			dataFetcherServiceAddLikeHandler.ServeHTTP(w, r)
		case DataFetcherServiceRemoveLikeProcedure:
			dataFetcherServiceRemoveLikeHandler.ServeHTTP(w, r)
		case DataFetcherServiceAddBookmarkProcedure:
			dataFetcherServiceAddBookmarkHandler.ServeHTTP(w, r)
		case DataFetcherServiceRemoveBookmarkProcedure:
			dataFetcherServiceRemoveBookmarkHandler.ServeHTTP(w, r)
		case DataFetcherServiceAddCommentProcedure:
			dataFetcherServiceAddCommentHandler.ServeHTTP(w, r)
		case DataFetcherServiceRemoveCommentProcedure:
			dataFetcherServiceRemoveCommentHandler.ServeHTTP(w, r)
		case DataFetcherServiceGetUserInfoProcedure:
			dataFetcherServiceGetUserInfoHandler.ServeHTTP(w, r)
		case DataFetcherServiceSignUpUserProcedure:
			dataFetcherServiceSignUpUserHandler.ServeHTTP(w, r)
		case DataFetcherServiceUpdateUserProcedure:
			dataFetcherServiceUpdateUserHandler.ServeHTTP(w, r)
		case DataFetcherServiceHardDeleteUserProcedure:
			dataFetcherServiceHardDeleteUserHandler.ServeHTTP(w, r)
		case DataFetcherServiceSoftDeleteUserProcedure:
			dataFetcherServiceSoftDeleteUserHandler.ServeHTTP(w, r)
		case DataFetcherServiceGetUserInfoByOAuthProcedure:
			dataFetcherServiceGetUserInfoByOAuthHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedDataFetcherServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedDataFetcherServiceHandler struct{}

func (UnimplementedDataFetcherServiceHandler) TestConnection(context.Context, *connect.Request[pb.TestConnectionRequest]) (*connect.Response[pb.TestConnectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.TestConnection is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) GetPostsByLocation(context.Context, *connect.Request[pb.GetPostsByLocationRequest]) (*connect.Response[pb.GetPostsByLocationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.GetPostsByLocation is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) GetPostsByPostIds(context.Context, *connect.Request[pb.GetPostsByPostIdsRequest]) (*connect.Response[pb.GetPostsByPostIdsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.GetPostsByPostIds is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) GetPostsByUser(context.Context, *connect.Request[pb.GetPostsByUserIdRequest]) (*connect.Response[pb.GetPostsByUserIdResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.GetPostsByUser is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) AddPost(context.Context, *connect.Request[pb.AddPostRequest]) (*connect.Response[pb.AddPostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.AddPost is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) HardDeletePost(context.Context, *connect.Request[pb.HardDeletePostRequest]) (*connect.Response[pb.HardDeletePostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.HardDeletePost is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) SoftDeletePost(context.Context, *connect.Request[pb.SoftDeletePostRequest]) (*connect.Response[pb.SoftDeletePostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.SoftDeletePost is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) AddLike(context.Context, *connect.Request[pb.AddLikeRequest]) (*connect.Response[pb.AddLikeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.AddLike is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) RemoveLike(context.Context, *connect.Request[pb.RemoveLikeRequest]) (*connect.Response[pb.RemoveLikeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.RemoveLike is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) AddBookmark(context.Context, *connect.Request[pb.AddBookmarkRequest]) (*connect.Response[pb.AddBookmarkResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.AddBookmark is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) RemoveBookmark(context.Context, *connect.Request[pb.RemoveBookmarkRequest]) (*connect.Response[pb.RemoveBookmarkResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.RemoveBookmark is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) AddComment(context.Context, *connect.Request[pb.AddCommentRequest]) (*connect.Response[pb.AddCommentResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.AddComment is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) RemoveComment(context.Context, *connect.Request[pb.RemoveCommentRequest]) (*connect.Response[pb.RemoveCommentResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.RemoveComment is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) GetUserInfo(context.Context, *connect.Request[pb.GetUserInfoRequest]) (*connect.Response[pb.GetUserInfoResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.GetUserInfo is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) SignUpUser(context.Context, *connect.Request[pb.SignUpUserRequest]) (*connect.Response[pb.SignUpUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.SignUpUser is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) UpdateUser(context.Context, *connect.Request[pb.UpdateUserRequest]) (*connect.Response[pb.UpdateUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.UpdateUser is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) HardDeleteUser(context.Context, *connect.Request[pb.HardDeleteUserRequest]) (*connect.Response[pb.HardDeleteUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.HardDeleteUser is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) SoftDeleteUser(context.Context, *connect.Request[pb.SoftDeleteUserRequest]) (*connect.Response[pb.SoftDeleteUserResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.SoftDeleteUser is not implemented"))
}

func (UnimplementedDataFetcherServiceHandler) GetUserInfoByOAuth(context.Context, *connect.Request[pb.GetUserInfoByOAuthRequest]) (*connect.Response[pb.GetUserInfoByOAuthResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("rpc.DataFetcherService.GetUserInfoByOAuth is not implemented"))
}
