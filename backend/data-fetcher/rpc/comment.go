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

func generateUniqueCommentID(ctx context.Context, postID string, pc cache.PostCacheService, db database.DatabaseService) (string, error) {
	exist, err := pc.IsKeyExist(ctx, "post:"+postID)

	if exist && err == nil {
		for {
			id := postID + uuid.NewString()
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
