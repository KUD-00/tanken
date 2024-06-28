package db

import (
	"context"
	"tanken/backend/common/types"
)

type DatabaseService interface {
	GetPost(ctx context.Context, postId string) (*types.Post, error)
	SetPost(ctx context.Context, postId string, post *types.Post) error
	SoftDeletePost(ctx context.Context, postId string) error
	HardDeletePost(ctx context.Context, postId string) error

	GetPostDetails(ctx context.Context, postId string) (*types.PostDetails, error)
	SetPostDetails(ctx context.Context, postId string, post *types.PostDetailsPtr) error

	GetPostLikedBy(ctx context.Context, postId string) ([]string, error)
	AddPostLikedBy(ctx context.Context, postId string, userId []string) error
	DeletePostLikedBy(ctx context.Context, postId string, userId []string) error

	GetPostBookmarkedBy(ctx context.Context, postId string) ([]string, error)
	AddPostBookmarkedBy(ctx context.Context, postId string, userId []string) error
	DeletePostBookmarkedBy(ctx context.Context, postId string, userId []string) error

	GetPostTags(ctx context.Context, postId string) ([]string, error)
	AddPostTags(ctx context.Context, postId string, tags []string) error
	DeletePostTags(ctx context.Context, postId string, tags []string) error

	GetPostPictureLinks(ctx context.Context, postId string) ([]string, error)
	AddPostPictureLinks(ctx context.Context, postId string, pictureLinks []string) error
	DeletePostPictureLinks(ctx context.Context, postId string, pictureLinks []string) error

	GetPostCommentIds(ctx context.Context, postId string) ([]string, error)
	AddPostCommentIds(ctx context.Context, postId string, commentIds []string) error
	DeletePostCommentIds(ctx context.Context, postId string, commentIds []string) error

	GetCommentById(ctx context.Context, commentId string) (*types.Comment, error)
	SetCommentById(ctx context.Context, commentId string, comment *types.CommentPtr) error
	SoftDeleteCommentById(ctx context.Context, commentId string) error
	HardDeleteCommentById(ctx context.Context, commentId string) error

	// About User
	GetUserById(ctx context.Context, userId string) (*types.User, error)
	GetUserByOauthInfo(ctx context.Context, email string, oauthProvider string) (*types.User, error)
	SetUserById(ctx context.Context, userId string, user *types.UserPtr) error
	SoftDeleteUserById(ctx context.Context, userId string) error
	HardDeleteUserById(ctx context.Context, userId string) error

	GetUserLikedPosts(ctx context.Context, userId string) ([]string, error)
	GetUserBookmarkedPosts(ctx context.Context, userId string) ([]string, error)
}
