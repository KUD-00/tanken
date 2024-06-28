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

func (p *PostgresDatabaseService) SetCommentById(ctx context.Context, commentId string, comment *types.CommentPtr) error {
	insertFields := []string{"comment_id"}
	insertValues := []string{"$1"}
	updateFields := []string{}
	values := []interface{}{commentId}
	placeholderCount := 2

	if comment.PostId != nil {
		insertFields = append(insertFields, "post_id")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "post_id = EXCLUDED.post_id")
		values = append(values, *comment.PostId)
		placeholderCount++
	}

	if comment.UserId != nil {
		insertFields = append(insertFields, "user_id")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "user_id = EXCLUDED.user_id")
		values = append(values, *comment.UserId)
		placeholderCount++
	}

	if comment.Content != nil {
		insertFields = append(insertFields, "content")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "content = EXCLUDED.content")
		values = append(values, *comment.Content)
		placeholderCount++
	}

	if comment.CreatedAt != nil {
		insertFields = append(insertFields, "created_at")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "created_at = EXCLUDED.created_at")
		values = append(values, *comment.CreatedAt)
		placeholderCount++
	}

	if comment.UpdatedAt != nil {
		insertFields = append(insertFields, "updated_at")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "updated_at = EXCLUDED.updated_at")
		values = append(values, *comment.UpdatedAt)
		placeholderCount++
	}

	if comment.Likes != nil {
		insertFields = append(insertFields, "likes")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "likes = EXCLUDED.likes")
		values = append(values, *comment.Likes)
		placeholderCount++
	}

	if len(updateFields) == 0 {
		// If no fields to update, just return without doing anything
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO comments (%s)
		VALUES (%s)
		ON CONFLICT (comment_id) DO UPDATE SET %s`,
		strings.Join(insertFields, ", "),
		strings.Join(insertValues, ", "),
		strings.Join(updateFields, ", "))

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) SoftDeleteCommentById(ctx context.Context, commentId string) error {
	_, err := p.db.ExecContext(ctx, "UPDATE comments SET status = 0 WHERE comment_id = $1", commentId)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) HardDeleteCommentById(ctx context.Context, commentId string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM comments WHERE comment_id = $1", commentId)
	if err != nil {
		return err
	}

	return nil
}
