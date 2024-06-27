package db

import (
	"context"
	"testing"
	"time"

	postgres "tanken/backend/common/db/postgres"
	"tanken/backend/common/types"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var (
	createdAt = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	updatedAt = time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC).Unix()
)

// These are unit tests using sqlmock

func TestPostgresDatabaseService_GetPostDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "content", "created_at", "updated_at", "likes", "latitude", "longitude",
		"status"}).
		AddRow("1", "user1", "content", createdAt, updatedAt, 10, 1.23, 4.56, 1)
	mock.ExpectQuery("SELECT post_id, user_id, content, created_at, updated_at, likes, latitude, longitude, status, FROM posts WHERE post_id = \\$1").
		WithArgs("1").
		WillReturnRows(rows)

	ctx := context.Background()
	post, err := service.GetPostDetails(ctx, "1")

	assert.NoError(t, err)

	expectedPost := types.PostDetails{
		PostId:    "1",
		UserId:    "user1",
		Content:   "content",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Likes:     10,
		Status:    1,

		Location: types.Location{
			Latitude:  1.23,
			Longitude: 4.56,
		},
	}

	assert.NotNil(t, post)
	assert.Equal(t, expectedPost, *post)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_GetPostsDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "content", "created_at", "updated_at", "likes", "latitude", "longitude", "status"}).
		AddRow("1", "user1", "content1", createdAt, updatedAt, int64(10), 1.23, 4.56, int64(1)).
		AddRow("2", "user2", "content2", createdAt, updatedAt, int64(20), 2.34, 5.67, int64(2))
	mock.ExpectQuery("SELECT post_id, user_id, content, created_at, updated_at, likes, latitude, longitude, status FROM posts WHERE post_id = ANY\\(\\$1\\)").
		WithArgs(pq.Array([]string{"1", "2"})).
		WillReturnRows(rows)

	ctx := context.Background()
	postIds := []string{"1", "2"}
	posts, err := service.GetPostsDetails(ctx, postIds)

	assert.NoError(t, err)
	assert.NotNil(t, posts)

	expectedPosts := []types.PostDetails{
		{
			PostId:    "1",
			UserId:    "user1",
			Content:   "content1",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Likes:     int64(10),
			Status:    int64(1),
			Location: types.Location{
				Latitude:  1.23,
				Longitude: 4.56,
			},
		},
		{
			PostId:    "2",
			UserId:    "user2",
			Content:   "content2",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Likes:     int64(20),
			Status:    int64(2),
			Location: types.Location{
				Latitude:  2.34,
				Longitude: 5.67,
			},
		},
	}

	assert.Equal(t, expectedPosts, *posts)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_SetPostDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	postId := "1"
	post := &types.PostDetails{
		PostId:    postId,
		UserId:    "user1",
		Content:   "content1",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Likes:     int64(10),
		Status:    int64(1),
		Location: types.Location{
			Latitude:  1.23,
			Longitude: 4.56,
		},
	}

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(postId, post.UserId, post.Content, post.CreatedAt, post.UpdatedAt, post.Likes, post.Location.Latitude, post.Location.Longitude, post.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = service.SetPostDetails(ctx, postId, post)

	if err != nil {
		t.Errorf("error was not expected while setting post details: %s", err)
	}

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_SetPostsDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	posts := []types.PostDetails{
		{
			PostId:    "1",
			UserId:    "user1",
			Content:   "content1",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Likes:     int64(10),
			Status:    int64(1),
			Location: types.Location{
				Latitude:  1.23,
				Longitude: 4.56,
			},
		},
		{
			PostId:    "2",
			UserId:    "user2",
			Content:   "content2",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Likes:     int64(20),
			Status:    int64(2),
			Location: types.Location{
				Latitude:  2.34,
				Longitude: 5.67,
			},
		},
	}

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(
			"1", "user1", "content1", createdAt, updatedAt, int64(10), 1.23, 4.56, int64(1),
			"2", "user2", "content2", createdAt, updatedAt, int64(20), 2.34, 5.67, int64(2),
		).
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.SetPostsDetails(ctx, posts)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_GetPostLikedBy(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	mockRows := sqlmock.NewRows([]string{"user_id"}).AddRow("user1").AddRow("user2")
	mock.ExpectQuery("SELECT user_id FROM user_liked_posts WHERE post_id = \\$1").
		WithArgs("post1").
		WillReturnRows(mockRows)

	ctx := context.Background()
	users, err := service.GetPostLikedBy(ctx, "post1")

	assert.NoError(t, err)
	assert.Equal(t, []string{"user1", "user2"}, users)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_AddPostLikedBy(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	mock.ExpectExec("INSERT INTO user_liked_posts \\(post_id, user_id\\) VALUES \\(\\$1, \\$2\\), \\(\\$1, \\$3\\) ON CONFLICT DO NOTHING").
		WithArgs("post1", "user1", "user2").
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.AddPostLikedBy(ctx, "post1", []string{"user1", "user2"})

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_DeletePostLikedBy(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	mock.ExpectExec("DELETE FROM user_liked_posts WHERE post_id = \\$1 AND user_id = ANY\\(\\$2\\)").
		WithArgs("post1", pq.Array([]string{"user1", "user2"})).
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.DeletePostLikedBy(ctx, "post1", []string{"user1", "user2"})

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_GetPostTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	rows := sqlmock.NewRows([]string{"tag"}).
		AddRow("tag1").
		AddRow("tag2")

	mock.ExpectQuery("SELECT tag FROM post_tags WHERE post_id = \\$1").
		WithArgs("post1").
		WillReturnRows(rows)

	ctx := context.Background()
	tags, err := service.GetPostTags(ctx, "post1")

	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2"}, tags)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_AddPostTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	tags := []string{"tag1", "tag2"}

	// Expect a single batch insert
	mock.ExpectExec("INSERT INTO post_tags \\(post_id, tag\\) VALUES \\(\\$1, \\$2\\), \\(\\$1, \\$3\\) ON CONFLICT DO NOTHING").
		WithArgs("post1", "tag1", "tag2").
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.AddPostTags(ctx, "post1", tags)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_DeletePostTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	tags := []string{"tag1", "tag2"}

	// Expect a single batch delete
	mock.ExpectExec("DELETE FROM post_tags WHERE post_id = \\$1 AND tag = ANY\\(\\$2\\)").
		WithArgs("post1", pq.Array(tags)).
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.DeletePostTags(ctx, "post1", tags)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_GetPostPictureLinks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	rows := sqlmock.NewRows([]string{"link"}).
		AddRow("link1").
		AddRow("link2")

	mock.ExpectQuery("SELECT link FROM post_picture_links WHERE post_id = \\$1").
		WithArgs("post1").
		WillReturnRows(rows)

	ctx := context.Background()
	links, err := service.GetPostPictureLinks(ctx, "post1")

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"link1", "link2"}, links)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_AddPostPictureLinks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	pictureLinks := []string{"link1", "link2"}

	query := "INSERT INTO post_picture_links \\(post_id, link\\) VALUES \\(\\$1, \\$2\\), \\(\\$1, \\$3\\) ON CONFLICT DO NOTHING"
	mock.ExpectExec(query).
		WithArgs("post1", "link1", "link2").
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.AddPostPictureLinks(ctx, "post1", pictureLinks)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_DeletePostPictureLinks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	pictureLinks := []string{"link1", "link2"}

	query := "DELETE FROM post_picture_links WHERE post_id = \\$1 AND link = ANY\\(\\$2\\)"
	mock.ExpectExec(query).
		WithArgs("post1", pq.Array(pictureLinks)).
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.DeletePostPictureLinks(ctx, "post1", pictureLinks)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_GetPostCommentIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	rows := sqlmock.NewRows([]string{"comment_id"}).
		AddRow("comment1").
		AddRow("comment2")

	mock.ExpectQuery("SELECT comment_id FROM post_comments WHERE post_id = \\$1").
		WithArgs("post1").
		WillReturnRows(rows)

	ctx := context.Background()
	commentIds, err := service.GetPostCommentIds(ctx, "post1")

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"comment1", "comment2"}, commentIds)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_AddPostCommentIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	commentIds := []string{"comment1", "comment2"}

	query := "INSERT INTO post_comments \\(post_id, comment_id\\) VALUES \\(\\$1, \\$2\\), \\(\\$1, \\$3\\) ON CONFLICT DO NOTHING"
	mock.ExpectExec(query).
		WithArgs("post1", "comment1", "comment2").
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.AddPostCommentIds(ctx, "post1", commentIds)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_DeletePostCommentIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	commentIds := []string{"comment1", "comment2"}

	query := "DELETE FROM post_comments WHERE post_id = \\$1 AND comment_id = ANY\\(\\$2\\)"
	mock.ExpectExec(query).
		WithArgs("post1", pq.Array(commentIds)).
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.DeletePostCommentIds(ctx, "post1", commentIds)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_GetCommentById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	// Mock data
	commentId := "comment1"
	expectedComment := &types.Comment{
		CommentId: commentId,
		PostId:    "post1",
		UserId:    "user1",
		Content:   "This is a comment.",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Likes:     10,
	}

	rows := sqlmock.NewRows([]string{"comment_id", "post_id", "user_id", "content", "created_at", "updated_at", "likes"}).
		AddRow(expectedComment.CommentId, expectedComment.PostId, expectedComment.UserId, expectedComment.Content, expectedComment.CreatedAt, expectedComment.UpdatedAt, expectedComment.Likes)

	mock.ExpectQuery("SELECT comment_id, post_id, user_id, content, created_at, updated_at, likes FROM comments WHERE comment_id = \\$1").
		WithArgs(commentId).
		WillReturnRows(rows)

	ctx := context.Background()
	comment, err := service.GetCommentById(ctx, commentId)

	assert.NoError(t, err)
	assert.Equal(t, expectedComment, comment)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_SetCommentById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	// Mock data
	commentId := "comment1"
	comment := &types.Comment{
		CommentId: commentId,
		PostId:    "post1",
		UserId:    "user1",
		Content:   "This is a comment.",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Likes:     10,
	}

	mock.ExpectExec("INSERT INTO comments \\(comment_id, post_id, user_id, content, created_at, updated_at, likes\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\) ON CONFLICT \\(comment_id\\) DO UPDATE SET post_id = EXCLUDED.post_id, user_id = EXCLUDED.user_id, content = EXCLUDED.content, created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at, likes = EXCLUDED.likes").
		WithArgs(comment.CommentId, comment.PostId, comment.UserId, comment.Content, comment.CreatedAt, comment.UpdatedAt, comment.Likes).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = service.SetCommentById(ctx, commentId, comment)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_HardDeleteCommentById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	commentId := "comment1"

	mock.ExpectExec("DELETE FROM comments WHERE comment_id = \\$1").
		WithArgs(commentId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = service.HardDeleteCommentById(ctx, commentId)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	expectedUser := &types.User{
		UserId:             "user1",
		Username:           "username1",
		Bio:                "This is a bio",
		ProfilePictureLink: "http://example.com/profile.jpg",
		Subscribed:         int64(1),
	}

	mock.ExpectQuery("SELECT user_id, username, bio, profile_picture_link, subscribed FROM users WHERE user_id = \\$1").
		WithArgs("user1").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "bio", "profile_picture_link", "subscribed"}).
			AddRow(expectedUser.UserId, expectedUser.Username, expectedUser.Bio, expectedUser.ProfilePictureLink, expectedUser.Subscribed))

	ctx := context.Background()
	user, err := service.GetUserById(ctx, "user1")

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_SetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := postgres.NewPostgresDatabaseService(db)

	user := &types.UserPtr{
		Username:           strPtr("new_username"),
		Bio:                strPtr("new_bio"),
		ProfilePictureLink: strPtr("http://example.com/new_profile.jpg"),
		Subscribed:         int64Ptr(int64(1)),
		Email:              strPtr("new_email@example.com"),
		OauthProvider:      strPtr("new_oauth_provider"),
	}

	mock.ExpectExec(`INSERT INTO users \(user_id, username, bio, profile_picture_link, subscribed, email, oauth_provider\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\) ON CONFLICT \(user_id\) DO UPDATE SET username = EXCLUDED.username, bio = EXCLUDED.bio, profile_picture_link = EXCLUDED.profile_picture_link, subscribed = EXCLUDED.subscribed, email = EXCLUDED.email, oauth_provider = EXCLUDED.oauth_provider`).
		WithArgs("user1", "new_username", "new_bio", "http://example.com/new_profile.jpg", int64(1), "new_email@example.com", "new_oauth_provider").
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = service.SetUserById(ctx, "user1", user)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// 辅助函数
func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func int64Ptr(i int64) *int64 {
	return &i
}
