package db

import (
	"context"
	"tanken/backend/common/types"
)

type DatabaseService interface {
	GetPost(ctx context.Context, postID string) (*types.Post, error)
	SetPost(ctx context.Context, postID string, post *types.Post) error
	DeletePost(ctx context.Context, postID string) error

	GetPostDetails(ctx context.Context, postID string) (*types.PostDetails, error)
	SetPostDetails(ctx context.Context, postID string, post *types.PostDetails) error

	GetPostLikedBy(ctx context.Context, postID string) ([]string, error)
	AddPostLikedBy(ctx context.Context, postID string, userId []string) error
	DeletePostLikedBy(ctx context.Context, postID string, userId []string) error

	GetPostTags(ctx context.Context, postID string) ([]string, error)
	AddPostTags(ctx context.Context, postID string, tags []string) error
	DeletePostTags(ctx context.Context, postID string, tags []string) error

	GetPostPictureLinks(ctx context.Context, postID string) ([]string, error)
	AddPostPictureLinks(ctx context.Context, postID string, pictureLinks []string) error
	DeletePostPictureLinks(ctx context.Context, postID string, pictureLinks []string) error

	GetPostCommentIds(ctx context.Context, postID string) ([]string, error)
	AddPostCommentIds(ctx context.Context, postID string, commentIds []string) error
	DeletePostCommentIds(ctx context.Context, postID string, commentIds []string) error

	GetCommentById(ctx context.Context, commentID string) (*types.Comment, error)
	SetCommentById(ctx context.Context, commentID string, comment *types.Comment) error
	DeleteCommentById(ctx context.Context, commentID string) error

	GetUserById(ctx context.Context, userID string) (*types.User, error)
	SetUserById(ctx context.Context, userID string, user *types.UserPtr) error
	DeleteUserById(ctx context.Context, userID string) error
}
