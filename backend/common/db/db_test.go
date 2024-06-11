package db

import (
	"context"
	"testing"
	"time"

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "content", "created_at", "updated_at", "likes", "latitude", "longitude", "status"}).
		AddRow("1", "user1", "content1", createdAt, updatedAt, int64(10), 1.23, 4.56, int64(1)).
		AddRow("2", "user2", "content2", createdAt, updatedAt, int64(20), 2.34, 5.67, int64(2))
	mock.ExpectQuery("SELECT post_id, user_id, content, created_at, updated_at, likes, latitude, longitude, status FROM posts WHERE post_id = ANY\\(\\$1\\)").
		WithArgs(pq.Array([]string{"1", "2"})).
		WillReturnRows(rows)

	ctx := context.Background()
	postIDs := []string{"1", "2"}
	posts, err := service.GetPostsDetails(ctx, postIDs)

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

	mockRows := sqlmock.NewRows([]string{"user_id"}).AddRow("user1").AddRow("user2")
	mock.ExpectQuery("SELECT user_id FROM post_likes WHERE post_id = \\$1").
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

	service := NewPostgresDatabaseService(db)

	mock.ExpectExec("INSERT INTO post_likes \\(post_id, user_id\\) VALUES \\(\\$1, \\$2\\), \\(\\$1, \\$3\\) ON CONFLICT DO NOTHING").
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

	service := NewPostgresDatabaseService(db)

	mock.ExpectExec("DELETE FROM post_likes WHERE post_id = \\$1 AND user_id = ANY\\(\\$2\\)").
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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

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

	service := NewPostgresDatabaseService(db)

	rows := sqlmock.NewRows([]string{"comment_id"}).
		AddRow("comment1").
		AddRow("comment2")

	mock.ExpectQuery("SELECT comment_id FROM post_comments WHERE post_id = \\$1").
		WithArgs("post1").
		WillReturnRows(rows)

	ctx := context.Background()
	commentIDs, err := service.GetPostCommentIds(ctx, "post1")

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"comment1", "comment2"}, commentIDs)
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

	service := NewPostgresDatabaseService(db)

	commentIDs := []string{"comment1", "comment2"}

	query := "INSERT INTO post_comments \\(post_id, comment_id\\) VALUES \\(\\$1, \\$2\\), \\(\\$1, \\$3\\) ON CONFLICT DO NOTHING"
	mock.ExpectExec(query).
		WithArgs("post1", "comment1", "comment2").
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.AddPostCommentIds(ctx, "post1", commentIDs)

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

	service := NewPostgresDatabaseService(db)

	commentIDs := []string{"comment1", "comment2"}

	query := "DELETE FROM post_comments WHERE post_id = \\$1 AND comment_id = ANY\\(\\$2\\)"
	mock.ExpectExec(query).
		WithArgs("post1", pq.Array(commentIDs)).
		WillReturnResult(sqlmock.NewResult(1, 2))

	ctx := context.Background()
	err = service.DeletePostCommentIds(ctx, "post1", commentIDs)

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

	service := NewPostgresDatabaseService(db)

	// Mock data
	commentID := "comment1"
	expectedComment := &types.Comment{
		CommentId: commentID,
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
		WithArgs(commentID).
		WillReturnRows(rows)

	ctx := context.Background()
	comment, err := service.GetCommentById(ctx, commentID)

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

	service := NewPostgresDatabaseService(db)

	// Mock data
	commentID := "comment1"
	comment := &types.Comment{
		CommentId: commentID,
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
	err = service.SetCommentById(ctx, commentID, comment)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostgresDatabaseService_DeleteCommentById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := NewPostgresDatabaseService(db)

	commentID := "comment1"

	mock.ExpectExec("DELETE FROM comments WHERE comment_id = \\$1").
		WithArgs(commentID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = service.DeleteCommentById(ctx, commentID)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
