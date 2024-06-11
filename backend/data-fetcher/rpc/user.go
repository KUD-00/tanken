package rpc

import (
	"context"
	"database/sql"
	"fmt"

	"tanken/backend/common/cache"
	database "tanken/backend/common/db"
	types "tanken/backend/common/types"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// TODO: Do I need `defer pipe.Close()`?
func generateUniqueUserID(ctx context.Context, db database.DatabaseService) (string, error) {
	for {
		id := uuid.NewString()[:8]
		user, err := db.GetUserById(ctx, id)
		if err != nil {
			return "Error generating user id", err
		}
		if user.UserId == "" {
			return id, nil
		}
	}
}

func getUser(ctx context.Context, userId string, pc cache.PostCacheService, db *sql.DB) (*types.User, error) {
	exists, err := pc.IsKeyExist(ctx, "user:"+userId)

	if err != nil {
		return nil, fmt.Errorf("error checking Redis: %v", err)
	}

	if exists {
		user, err := pc.GetUser(ctx, userId)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	user, err := getUserFromDB(ctx, userId, pc, db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getUserFromDB(ctx context.Context, userId string, pc cache.PostCacheService, db *sql.DB) (*types.User, error) {
	var user types.User

	err := db.QueryRowContext(ctx, "SELECT user_id, username, bio, avatar, subscribed FROM users WHERE user_id = $1", userId).Scan(&user.UserId, &user.Username, &user.Bio, &user.Avatar, &user.Subscribed)
	if err != nil {
		return nil, fmt.Errorf("error fetching user from PostgreSQL: %v", err)
	}

	if err := pc.SetUser(ctx, userId, &user); err != nil {
		return &user, err
	}

	return &user, nil
}
