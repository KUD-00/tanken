package cache

import (
	"context"
	"strconv"
	"tanken/backend/common/types"
	"tanken/backend/common/utils"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

var (
	createdAt = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	updatedAt = time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC).Unix()
)

func TestPostRedisCacheService_GetPostDetailsCmd(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	ctx, pipe := service.NewPipe(ctx)

	postId := "post123"
	expectedResult := map[string]string{"title": "Test Post", "content": "This is a test post"}

	mock.ExpectHGetAll(postId).SetVal(expectedResult)

	cmd, err := service.GetPostDetailsCmd(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	_, err = pipe.Exec(ctx)
	assert.NoError(t, err)

	result, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostDetails(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := map[string]string{
		"UpdatedAt": strconv.FormatInt(updatedAt, 10),
		"CreatedAt": strconv.FormatInt(createdAt, 10),
		"UserId":    "user123",
		"Content":   "This is a test post",
		"Likes":     "100",
		"Bookmarks": "50",
		"Status":    "1",
	}

	mock.ExpectHGetAll(postId).SetVal(expectedResult)

	details, err := service.GetPostDetails(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, details)

	assert.Equal(t, updatedAt, *details.UpdatedAt)
	assert.Equal(t, createdAt, *details.CreatedAt)
	assert.Equal(t, "user123", *details.UserId)
	assert.Equal(t, "This is a test post", *details.Content)
	assert.Equal(t, int64(100), *details.Likes)
	assert.Equal(t, int64(50), *details.Bookmarks)
	assert.Equal(t, int64(1), *details.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_SetPostDetails(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	postDetails := &types.PostDetailsPtr{
		UpdatedAt: utils.Int64Ptr(1672444800), // 2023-01-01T00:00:00Z
		CreatedAt: utils.Int64Ptr(1640995200), // 2022-01-01T00:00:00Z
		UserId:    utils.StringPtr("user123"),
		Content:   utils.StringPtr("This is a test post"),
		Likes:     utils.Int64Ptr(100),
		Bookmarks: utils.Int64Ptr(50),
		Status:    utils.Int64Ptr(1),
	}

	postDetailsMap := map[string]interface{}{
		utils.PostCacheKeys.UpdatedAt: "1672444800",
		utils.PostCacheKeys.CreatedAt: "1640995200",
		utils.PostCacheKeys.UserId:    "user123",
		utils.PostCacheKeys.Content:   "This is a test post",
		utils.PostCacheKeys.Likes:     "100",
		utils.PostCacheKeys.Bookmarks: "50",
		utils.PostCacheKeys.Status:    "1",
	}

	mock.ExpectHSet(postId, postDetailsMap).SetVal(int64(1))

	err := service.SetPostDetails(ctx, postId, postDetails)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostLikedByCmd(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"user1", "user2", "user3"}

	ctx, pipe := service.NewPipe(ctx)

	mock.ExpectSMembers(postId + utils.LikedBySuffix).SetVal(expectedResult)

	cmd, err := service.GetPostLikedByCmd(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	_, err = pipe.Exec(ctx)
	assert.NoError(t, err)

	result, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostLikedBy(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"user1", "user2", "user3"}

	mock.ExpectSMembers(postId + utils.LikedBySuffix).SetVal(expectedResult)

	result, err := service.GetPostLikedBy(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_AddPostLikedBy(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	userIds := []string{"user1", "user2", "user3"}

	args := make([]interface{}, len(userIds))
	for i, v := range userIds {
		args[i] = v
	}

	mock.ExpectSAdd(postId+utils.LikedBySuffix, args...).SetVal(3)

	err := service.AddPostLikedBy(ctx, postId, userIds)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_RemovePostLikedBy(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	userIds := []string{"user1", "user2", "user3"}

	args := make([]interface{}, len(userIds))
	for i, v := range userIds {
		args[i] = v
	}

	mock.ExpectSRem(postId+utils.LikedBySuffix, args...).SetVal(3)

	err := service.RemovePostLikedBy(ctx, postId, userIds)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostTagsCmd(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"tag1", "tag2", "tag3"}

	ctx, pipe := service.NewPipe(ctx)

	mock.ExpectSMembers(postId + utils.TagsSuffix).SetVal(expectedResult)

	cmd, err := service.GetPostTagsCmd(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	_, err = pipe.Exec(ctx)
	assert.NoError(t, err)

	result, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostTags(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"tag1", "tag2", "tag3"}

	mock.ExpectSMembers(postId + utils.TagsSuffix).SetVal(expectedResult)

	result, err := service.GetPostTags(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_AddPostTags(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	tags := []string{"tag1", "tag2", "tag3"}

	args := make([]interface{}, len(tags))
	for i, v := range tags {
		args[i] = v
	}

	mock.ExpectSAdd(postId+utils.TagsSuffix, args...).SetVal(3)

	err := service.AddPostTags(ctx, postId, tags)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_RemovePostTags(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	tags := []string{"tag1", "tag2", "tag3"}

	args := make([]interface{}, len(tags))
	for i, v := range tags {
		args[i] = v
	}

	mock.ExpectSRem(postId+utils.TagsSuffix, args...).SetVal(3)

	err := service.RemovePostTags(ctx, postId, tags)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostPictureLinksCmd(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"link1", "link2", "link3"}

	ctx, pipe := service.NewPipe(ctx)
	mock.ExpectSMembers(postId + utils.PictureLinksSuffix).SetVal(expectedResult)

	cmd, err := service.GetPostPictureLinksCmd(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	_, err = pipe.Exec(ctx)
	assert.NoError(t, err)

	result, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostPictureLinks(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"link1", "link2", "link3"}

	mock.ExpectSMembers(postId + utils.PictureLinksSuffix).SetVal(expectedResult)

	result, err := service.GetPostPictureLinks(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_AddPostPictureLinks(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	pictureLinks := []string{"link1", "link2", "link3"}

	args := make([]interface{}, len(pictureLinks))
	for i, v := range pictureLinks {
		args[i] = v
	}

	mock.ExpectSAdd(postId+utils.PictureLinksSuffix, args...).SetVal(3)

	err := service.AddPostPictureLinks(ctx, postId, pictureLinks)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_RemovePostPictureLinks(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	pictureLinks := []string{"link1", "link2", "link3"}

	args := make([]interface{}, len(pictureLinks))
	for i, v := range pictureLinks {
		args[i] = v
	}

	mock.ExpectSRem(postId+utils.PictureLinksSuffix, args...).SetVal(3)

	err := service.RemovePostPictureLinks(ctx, postId, pictureLinks)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostCommentIdsCmd(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"comment1", "comment2", "comment3"}

	ctx, pipe := service.NewPipe(ctx)
	mock.ExpectSMembers(postId + utils.CommentIdsSuffix).SetVal(expectedResult)

	cmd, err := service.GetPostCommentIdsCmd(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	_, err = pipe.Exec(ctx)
	assert.NoError(t, err)

	result, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetPostCommentIds(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	expectedResult := []string{"comment1", "comment2", "comment3"}

	mock.ExpectSMembers(postId + utils.CommentIdsSuffix).SetVal(expectedResult)

	result, err := service.GetPostCommentIds(ctx, postId)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_AddPostCommentIds(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	commentIds := []string{"comment1", "comment2", "comment3"}

	args := make([]interface{}, len(commentIds))
	for i, v := range commentIds {
		args[i] = v
	}

	mock.ExpectSAdd(postId+utils.CommentIdsSuffix, args...).SetVal(3)

	err := service.AddPostCommentIds(ctx, postId, commentIds)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_RemovePostCommentIds(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	postId := "post123"
	commentIds := []string{"comment1", "comment2", "comment3"}

	args := make([]interface{}, len(commentIds))
	for i, v := range commentIds {
		args[i] = v
	}

	mock.ExpectSRem(postId+utils.CommentIdsSuffix, args...).SetVal(3)

	err := service.RemovePostCommentIds(ctx, postId, commentIds)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetCommentCmd(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	commentId := "comment123"
	expectedResult := map[string]string{
		"createdAt": "2023-01-01T00:00:00Z",
		"updatedAt": "2023-01-02T00:00:00Z",
		"postId":    "post123",
		"userId":    "user123",
		"content":   "This is a comment",
		"likes":     "10",
	}

	ctx, pipe := service.NewPipe(ctx)
	mock.ExpectHGetAll(utils.CommentPrefix + commentId).SetVal(expectedResult)

	cmd, err := service.GetCommentCmd(ctx, commentId)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	_, err = pipe.Exec(ctx)
	assert.NoError(t, err)

	result, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetComment(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	commentId := "comment123"
	expectedResult := map[string]string{
		"createdAt": strconv.FormatInt(createdAt, 10),
		"updatedAt": strconv.FormatInt(updatedAt, 10),
		"postId":    "post123",
		"userId":    "user123",
		"content":   "This is a comment",
		"likes":     "10",
	}

	mock.ExpectHGetAll(utils.CommentPrefix + commentId).SetVal(expectedResult)

	result, err := service.GetComment(ctx, commentId)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedComment := &types.Comment{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		PostId:    "post123",
		UserId:    "user123",
		Content:   "This is a comment",
		Likes:     10,
	}
	assert.Equal(t, expectedComment, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_SetComment(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	commentId := "comment123"
	comment := &types.Comment{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		PostId:    "post123",
		UserId:    "user123",
		Content:   "This is a comment",
		Likes:     10,
	}

	commentMap := map[string]interface{}{
		"createdAt": comment.CreatedAt,
		"updatedAt": comment.UpdatedAt,
		"postId":    comment.PostId,
		"userId":    comment.UserId,
		"content":   comment.Content,
		"likes":     comment.Likes,
	}

	mock.ExpectHSet(utils.CommentPrefix+commentId, commentMap).SetVal(1)

	err := service.SetComment(ctx, commentId, comment)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_RemoveComment(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	commentIds := []string{"comment123"}

	mock.ExpectDel(utils.CommentPrefix + commentIds[0]).SetVal(1)

	err := service.RemoveComments(ctx, commentIds)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_SetUser(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	userId := "user123"
	user := &types.User{
		UserId:             "user123",
		Username:           "testuser",
		ProfilePictureLink: "avatar.png",
		Bio:                "This is a test user",
		Subscribed:         1,
	}

	userMap := map[string]interface{}{
		"userId":             user.UserId,
		"username":           user.Username,
		"profilePictureLink": user.ProfilePictureLink,
		"bio":                user.Bio,
		"subscribed":         user.Subscribed,
	}

	mock.ExpectHSet(utils.UserPrefix+userId, userMap).SetVal(1)

	err := service.SetUser(ctx, userId, user)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetUserCmd(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	userId := "user123"
	expectedResult := map[string]string{
		"userId":             "user123",
		"username":           "testuser",
		"profilePictureLink": "avatar.png",
		"bio":                "This is a test user",
		"subscribed":         strconv.FormatInt(1, 10),
	}

	ctx, pipe := service.NewPipe(ctx)
	mock.ExpectHGetAll(utils.UserPrefix + userId).SetVal(expectedResult)

	cmd, err := service.GetUserCmd(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	_, err = pipe.Exec(ctx)
	assert.NoError(t, err)

	result, err := cmd.Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_GetUser(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	userId := "user123"
	expectedResult := map[string]string{
		"userId":             "user123",
		"username":           "testuser",
		"profilePictureLink": "avatar.png",
		"bio":                "This is a test user",
		"subscribed":         strconv.FormatInt(1, 10),
	}

	mock.ExpectHGetAll(utils.UserPrefix + userId).SetVal(expectedResult)

	result, err := service.GetUser(ctx, userId)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	expectedUser := &types.User{
		UserId:             "user123",
		Username:           "testuser",
		ProfilePictureLink: "avatar.png",
		Bio:                "This is a test user",
		Subscribed:         1,
	}
	assert.Equal(t, expectedUser, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRedisCacheService_RemoveUser(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	service := &PostRedisCacheService{RedisBase: &RedisBase{client: client}}

	userId := "user123"

	mock.ExpectDel(utils.UserPrefix + userId).SetVal(1)

	err := service.RemoveUser(ctx, userId)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
