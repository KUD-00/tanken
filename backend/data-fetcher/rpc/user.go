package rpc

import (
	"context"
	"fmt"

	"tanken/backend/common/cache"
	database "tanken/backend/common/db"
	types "tanken/backend/common/types"
	utils "tanken/backend/common/utils"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// TODO: Do I need `defer pipe.Close()`?
func generateUniqueUserID(ctx context.Context, db database.DatabaseService) (string, error) {
	for {
		id := uuid.NewString()
		user, err := db.GetUserById(ctx, id)
		if err != nil {
			return "Error generating user id", err
		}
		if user.UserId == "" {
			return id, nil
		}
	}
}

func getUser(ctx context.Context, userId string, uc cache.UserCacheService, db database.DatabaseService) (*types.User, error) {
	exists, err := uc.IsKeyExist(ctx, "user:"+userId)

	if err != nil {
		return nil, fmt.Errorf("error checking Redis: %v", err)
	}

	if exists {
		user, err := uc.GetUser(ctx, userId)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	user, err := db.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func setUser(ctx context.Context, userId string, user *types.UserPtr, needWriteBackNow bool, needCache bool, uc cache.UserCacheService, db database.DatabaseService) error {
	exists, err := uc.IsKeyExist(ctx, "user:"+userId)
	if err != nil {
		return fmt.Errorf("error checking Redis: %v", err)
	}

	if exists {
		if needWriteBackNow {
			if err = db.SetUserById(ctx, userId, user); err != nil {
				return err
			}

			if err = uc.SetUserOptional(ctx, userId, user); err != nil {
				return err
			}
		} else {
			user.Changed = utils.BoolPtr(true)

			if err = uc.SetUserOptional(ctx, userId, user); err != nil {
				return err
			}
		}
		return nil

	} else {
		if err = db.SetUserById(ctx, userId, user); err != nil {
			return err
		}

		if needCache {
			user, err := db.GetUserById(ctx, userId)
			if err != nil {
				return err
			}

			if err = uc.SetUser(ctx, userId, user); err != nil {
				return err
			}
		}
	}

	return nil
}
