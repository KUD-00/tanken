package rpc

import (
	"context"
	"fmt"
	"time"

	types "tanken/backend/common/types"

	"tanken/backend/common/cache"
	database "tanken/backend/common/db"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// TODO: Do I need `defer pipe.Close()`?

// TODO: THIS IS NOT GOOD MAYBE
func createIDMap(slice []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, v := range slice {
		m[v] = struct{}{}
	}
	return m
}

func generateUniqueCommentID(ctx context.Context, postID string, pc cache.PostCacheService, db database.DatabaseService) (string, error) {
	exist, err := pc.IsKeyExist(ctx, "post:"+postID)

	if exist && err == nil {
		for {
			id := postID + uuid.NewString()[:8]
			existIds, err := db.GetPostCommentIds(ctx, postID)

			if err != nil {
				return "", fmt.Errorf("error getting db post comment IDs: %v", err)
			}

			existIdMap := createIDMap(existIds)

			if _, exists := existIdMap[id]; !exists {
				return id, nil
			}
		}
	}

	existIds, err := db.GetPostCommentIds(ctx, postID)
	existIdMap := createIDMap(existIds)

	if err != nil {
		return "", fmt.Errorf("error getting db post comment IDs: %v", err)
	}

	for {
		id := postID + uuid.NewString()[:8]

		if _, exists := existIdMap[id]; !exists {
			return id, nil
		}
	}
}

func cacheNewComment(ctx context.Context, postID string, content string, userId string, pc cache.PostCacheService, db database.DatabaseService) error {
	commentID, err := generateUniqueCommentID(ctx, postID, pc, db)

	if err != nil {
		return fmt.Errorf("error generating unique comment ID: %v", err)
	}

	comment := types.Comment{
		CommentId: commentID,
		PostId:    postID,
		UserId:    userId,
		Content:   content,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
		Likes:     0,
		Status:    int64(1),
	}

	return pc.SetComment(ctx, commentID, &comment)
}

func getComment(ctx context.Context, commentID string, pc cache.PostCacheService, db database.DatabaseService) (*types.Comment, error) {
	exists, err := pc.IsKeyExist(ctx, "comment:"+commentID)

	if err != nil {
		return nil, fmt.Errorf("error checking Redis: %v", err)
	}

	if exists {
		comment, err := pc.GetComment(ctx, commentID)
		if err != nil {
			return nil, err
		}
		return comment, nil
	}

	comment, err := db.GetCommentById(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
