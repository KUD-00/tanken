package postgres

import (
	"context"
	"database/sql"
	"strconv"
	"tanken/backend/common/types"
)

func (p *PostgresDatabaseService) GetUserById(ctx context.Context, userID string) (*types.User, error) {
	var user types.User

	err := p.db.QueryRowContext(ctx, "SELECT user_id, username, bio, profile_picture_link, subscribed FROM users WHERE user_id = $1", userID).Scan(
		&user.UserId,
		&user.Username,
		&user.Bio,
		&user.ProfilePictureLink,
		&user.Subscribed,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.User{}, nil
		}
		return nil, err
	}

	return &user, nil
}

func (p *PostgresDatabaseService) SetUserById(ctx context.Context, userID string, user *types.UserPtr) error {
	query := "INSERT INTO users (user_id"
	values := []interface{}{userID}
	updates := " ON CONFLICT (user_id) DO UPDATE SET"

	index := 2 // 第一个占位符是userID

	if user.Username != nil {
		query += ", username"
		values = append(values, *user.Username)
		updates += " username = EXCLUDED.username,"
		index++
	}
	if user.Bio != nil {
		query += ", bio"
		values = append(values, *user.Bio)
		updates += " bio = EXCLUDED.bio,"
		index++
	}
	if user.ProfilePictureLink != nil {
		query += ", profile_picture_link"
		values = append(values, *user.ProfilePictureLink)
		updates += " profile_picture_link = EXCLUDED.profile_picture_link,"
		index++
	}
	if user.Subscribed != nil {
		query += ", subscribed"
		values = append(values, *user.Subscribed)
		updates += " subscribed = EXCLUDED.subscribed,"
		index++
	}
	if user.Email != nil {
		query += ", email"
		values = append(values, *user.Email)
		updates += " email = EXCLUDED.email,"
		index++
	}
	if user.OauthProvider != nil {
		query += ", oauth_provider"
		values = append(values, *user.OauthProvider)
		updates += " oauth_provider = EXCLUDED.oauth_provider,"
		index++
	}

	query += ") VALUES ($1"
	for i := 2; i < index; i++ {
		query += ", $" + strconv.Itoa(i)
	}
	query += ")"

	updates = updates[:len(updates)-1] // 去掉最后的逗号
	query += updates

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) HardDeleteUserById(ctx context.Context, userID string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) SoftDeleteUserById(ctx context.Context, userId string) error {
	_, err := p.db.ExecContext(ctx, "UPDATE users SET status = 0 WHERE user_id = $1", userId)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetUserByOauthInfo(ctx context.Context, email string, oauthProvider string) (*types.User, error) {
	var user types.User

	err := p.db.QueryRowContext(ctx, "SELECT user_id, username, bio, profile_picture_link, subscribed FROM users WHERE email = $1 AND oauth_provider = $2", email, oauthProvider).Scan(
		&user.UserId,
		&user.Username,
		&user.Bio,
		&user.ProfilePictureLink,
		&user.Subscribed,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
