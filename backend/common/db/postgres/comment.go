package postgres

import (
	"context"
	"fmt"
	"strings"
	"tanken/backend/common/types"

	"github.com/lib/pq"
)

func (p *PostgresDatabaseService) GetPostCommentIds(ctx context.Context, postId string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT comment_id FROM post_comments WHERE post_id = $1", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentIds []string
	for rows.Next() {
		var commentId string
		if err := rows.Scan(&commentId); err != nil {
			return nil, err
		}
		commentIds = append(commentIds, commentId)
	}

	return commentIds, nil
}

func (p *PostgresDatabaseService) AddPostCommentIds(ctx context.Context, postId string, commentIds []string) error {
	if len(commentIds) == 0 {
		return nil
	}

	query := "INSERT INTO post_comments (post_id, comment_id) VALUES "
	values := []interface{}{postId}
	valueStrings := []string{}

	for i, commentId := range commentIds {
		valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, commentId)
	}

	query += strings.Join(valueStrings, ", ") + " ON CONFLICT DO NOTHING"

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) DeletePostCommentIds(ctx context.Context, postId string, commentIds []string) error {
	if len(commentIds) == 0 {
		return nil
	}

	query := "DELETE FROM post_comments WHERE post_id = $1 AND comment_id = ANY($2)"
	_, err := p.db.ExecContext(ctx, query, postId, pq.Array(commentIds))
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetCommentById(ctx context.Context, commentId string) (*types.Comment, error) {
	var comment types.Comment

	err := p.db.QueryRowContext(ctx, "SELECT comment_id, post_id, user_id, content, created_at, updated_at, likes FROM comments WHERE comment_id = $1", commentId).Scan(
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

func (p *PostgresDatabaseService) SetCommentById(ctx context.Context, commentId string, comment *types.Comment) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO comments (comment_id, post_id, user_id, content, created_at, updated_at, likes) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (comment_id) DO UPDATE SET post_id = EXCLUDED.post_id, user_id = EXCLUDED.user_id, content = EXCLUDED.content, created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at, likes = EXCLUDED.likes",
		commentId,
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

func (p *PostgresDatabaseService) DeleteCommentById(ctx context.Context, commentId string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM comments WHERE comment_id = $1", commentId)
	if err != nil {
		return err
	}

	return nil
}
