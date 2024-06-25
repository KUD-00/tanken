package postgres

import (
	"context"
	"fmt"
	"strings"
	"tanken/backend/common/types"

	"github.com/lib/pq"
)

func (p *PostgresDatabaseService) GetPostCommentIds(ctx context.Context, postID string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT comment_id FROM post_comments WHERE post_id = $1", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentIDs []string
	for rows.Next() {
		var commentID string
		if err := rows.Scan(&commentID); err != nil {
			return nil, err
		}
		commentIDs = append(commentIDs, commentID)
	}

	return commentIDs, nil
}

func (p *PostgresDatabaseService) AddPostCommentIds(ctx context.Context, postID string, commentIDs []string) error {
	if len(commentIDs) == 0 {
		return nil
	}

	query := "INSERT INTO post_comments (post_id, comment_id) VALUES "
	values := []interface{}{postID}
	valueStrings := []string{}

	for i, commentID := range commentIDs {
		valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, commentID)
	}

	query += strings.Join(valueStrings, ", ") + " ON CONFLICT DO NOTHING"

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) DeletePostCommentIds(ctx context.Context, postID string, commentIDs []string) error {
	if len(commentIDs) == 0 {
		return nil
	}

	query := "DELETE FROM post_comments WHERE post_id = $1 AND comment_id = ANY($2)"
	_, err := p.db.ExecContext(ctx, query, postID, pq.Array(commentIDs))
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetCommentById(ctx context.Context, commentID string) (*types.Comment, error) {
	var comment types.Comment

	err := p.db.QueryRowContext(ctx, "SELECT comment_id, post_id, user_id, content, created_at, updated_at, likes FROM comments WHERE comment_id = $1", commentID).Scan(
		&comment.CommentId,
		&comment.PostId,
		&comment.UserId,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Likes,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (p *PostgresDatabaseService) SetCommentById(ctx context.Context, commentID string, comment *types.Comment) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO comments (comment_id, post_id, user_id, content, created_at, updated_at, likes) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (comment_id) DO UPDATE SET post_id = EXCLUDED.post_id, user_id = EXCLUDED.user_id, content = EXCLUDED.content, created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at, likes = EXCLUDED.likes",
		commentID,
		comment.PostId,
		comment.UserId,
		comment.Content,
		comment.CreatedAt,
		comment.UpdatedAt,
		comment.Likes,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) DeleteCommentById(ctx context.Context, commentID string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM comments WHERE comment_id = $1", commentID)
	if err != nil {
		return err
	}

	return nil
}