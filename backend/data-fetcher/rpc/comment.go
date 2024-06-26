package rpc

import (
	"context"
	"fmt"

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

func generateUniqueCommentID(ctx context.Context, postId string, pc cache.PostCacheService, db database.DatabaseService) (string, error) {
	exist, err := pc.IsKeyExist(ctx, "post:"+postId)

	if exist && err == nil {
		for {
			id := postId + uuid.NewString()
			existIds, err := db.GetPostCommentIds(ctx, postId)

			if err != nil {
				return "", fmt.Errorf("error getting db post comment IDs: %v", err)
			}

			existIdMap := createIDMap(existIds)

			if _, exists := existIdMap[id]; !exists {
				return id, nil
			}
		}
	}

	existIds, err := db.GetPostCommentIds(ctx, postId)
	existIdMap := createIDMap(existIds)

	if err != nil {
		return "", fmt.Errorf("error getting db post comment IDs: %v", err)
	}

	for {
		id := postId + uuid.NewString()[:8]

		if _, exists := existIdMap[id]; !exists {
			return id, nil
		}
	}
}
