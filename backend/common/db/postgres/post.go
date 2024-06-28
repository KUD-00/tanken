package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"tanken/backend/common/types"

	"github.com/lib/pq"
)

type PostgresDatabaseService struct {
	db *sql.DB
}

func NewPostgresDatabaseService(db *sql.DB) *PostgresDatabaseService {
	return &PostgresDatabaseService{db: db}
}

func (p *PostgresDatabaseService) GetPost(ctx context.Context, postId string) (*types.Post, error) {
	query := `
        SELECT 
            p.post_id, p.user_id, p.content, p.created_at, p.updated_at, p.likes, p.latitude, p.longitude, p.status,
            pl.user_id AS liked_by_user_id,
            pt.tag AS tag,
            ppl.link AS picture_link
        FROM posts p
        LEFT JOIN user_liked_posts pl ON p.post_id = pl.post_id
        LEFT JOIN post_tags pt ON p.post_id = pt.post_id
        LEFT JOIN post_picture_links ppl ON p.post_id = ppl.post_id
        WHERE p.post_id = $1
    `

	rows, err := p.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postDetails := types.PostDetails{}
	var likedBys []string
	var tags []string
	var pictureLinks []string

	for rows.Next() {
		var likedBy sql.NullString
		var tag sql.NullString
		var pictureLink sql.NullString

		err := rows.Scan(
			&postDetails.PostId,
			&postDetails.UserId,
			&postDetails.Content,
			&postDetails.CreatedAt,
			&postDetails.UpdatedAt,
			&postDetails.Likes,
			&postDetails.Location.Latitude,
			&postDetails.Location.Longitude,
			&postDetails.Status,
			&likedBy,
			&tag,
			&pictureLink,
		)

		if err != nil {
			return nil, err
		}

		if likedBy.Valid {
			likedBys = append(likedBys, likedBy.String)
		}

		if tag.Valid {
			tags = append(tags, tag.String)
		}

		if pictureLink.Valid {
			pictureLinks = append(pictureLinks, pictureLink.String)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	postSets := types.PostSets{
		LikedBy:      likedBys,
		Tags:         tags,
		PictureLinks: pictureLinks,
	}

	post := types.Post{
		PostDetails: postDetails,
		PostSets:    postSets,
	}

	return &post, nil
}

func (p *PostgresDatabaseService) SetPost(ctx context.Context, postId string, post *types.Post) error {
	return nil
}

func (p *PostgresDatabaseService) SoftDeletePost(ctx context.Context, postId string) error {
	_, err := p.db.ExecContext(ctx, "UPDATE posts SET status = 0 WHERE user_id = $1", postId)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) HardDeletePost(ctx context.Context, postId string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM posts WHERE user_id = $1", postId)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetPostDetails(ctx context.Context, postId string) (*types.PostDetails, error) {
	var post types.PostDetails

	err := p.db.QueryRowContext(ctx, "SELECT post_id, user_id, content, created_at, updated_at, likes, latitude, longitude, status, FROM posts WHERE post_id = $1", postId).Scan(
		&post.PostId,
		&post.UserId,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Likes,
		&post.Location.Latitude,
		&post.Location.Longitude,
		&post.Status,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostgresDatabaseService) GetPostsDetails(ctx context.Context, postIds []string) (*[]types.PostDetails, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT post_id, user_id, content, created_at, updated_at, likes, latitude, longitude, status FROM posts WHERE post_id = ANY($1)", pq.Array(postIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []types.PostDetails
	for rows.Next() {
		var post types.PostDetails
		if err := rows.Scan(
			&post.PostId,
			&post.UserId,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Likes,
			&post.Location.Latitude,
			&post.Location.Longitude,
			&post.Status,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil
}

func (p *PostgresDatabaseService) SetPostDetails(ctx context.Context, postId string, post *types.PostDetailsPtr) error {
	insertFields := []string{"post_id"}
	insertValues := []string{"$1"}
	updateFields := []string{}
	values := []interface{}{postId}
	placeholderCount := 2

	if post.UserId != nil {
		insertFields = append(insertFields, "user_id")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "user_id = EXCLUDED.user_id")
		values = append(values, *post.UserId)
		placeholderCount++
	}

	if post.Content != nil {
		insertFields = append(insertFields, "content")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "content = EXCLUDED.content")
		values = append(values, *post.Content)
		placeholderCount++
	}

	if post.CreatedAt != nil {
		insertFields = append(insertFields, "created_at")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "created_at = EXCLUDED.created_at")
		values = append(values, *post.CreatedAt)
		placeholderCount++
	}

	if post.UpdatedAt != nil {
		insertFields = append(insertFields, "updated_at")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "updated_at = EXCLUDED.updated_at")
		values = append(values, *post.UpdatedAt)
		placeholderCount++
	}

	if post.Likes != nil {
		insertFields = append(insertFields, "likes")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "likes = EXCLUDED.likes")
		values = append(values, *post.Likes)
		placeholderCount++
	}

	if post.Location != nil {
		insertFields = append(insertFields, "latitude", "longitude")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount), fmt.Sprintf("$%d", placeholderCount+1))
		updateFields = append(updateFields, "latitude = EXCLUDED.latitude", "longitude = EXCLUDED.longitude")
		values = append(values, post.Location.Latitude, post.Location.Longitude)
		placeholderCount += 2
	}

	if post.Status != nil {
		insertFields = append(insertFields, "status")
		insertValues = append(insertValues, fmt.Sprintf("$%d", placeholderCount))
		updateFields = append(updateFields, "status = EXCLUDED.status")
		values = append(values, *post.Status)
		placeholderCount++
	}

	if len(updateFields) == 0 {
		// If no fields to update, just return without doing anything
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO posts (%s)
		VALUES (%s)
		ON CONFLICT (post_id) DO UPDATE SET %s`,
		strings.Join(insertFields, ", "),
		strings.Join(insertValues, ", "),
		strings.Join(updateFields, ", "))

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) SetPostsDetails(ctx context.Context, posts []types.PostDetails) error {
	valueStrings := make([]string, 0, len(posts))
	valueArgs := make([]interface{}, 0, len(posts)*9)

	for i, post := range posts {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9))
		valueArgs = append(valueArgs, post.PostId, post.UserId, post.Content, post.CreatedAt, post.UpdatedAt, post.Likes, post.Location.Latitude, post.Location.Longitude, post.Status)
	}

	stmt := fmt.Sprintf("INSERT INTO posts (post_id, user_id, content, created_at, updated_at, likes, latitude, longitude, status) VALUES %s ON CONFLICT (post_id) DO UPDATE SET user_id = EXCLUDED.user_id, content = EXCLUDED.content, created_at = EXCLUDED.created_at, updated_at = EXCLUDED.updated_at, likes = EXCLUDED.likes, latitude = EXCLUDED.latitude, longitude = EXCLUDED.longitude, status = EXCLUDED.status", strings.Join(valueStrings, ","))

	_, err := p.db.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetPostLikedBy(ctx context.Context, postId string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT user_id FROM user_liked_posts WHERE post_id = $1", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (p *PostgresDatabaseService) AddPostLikedBy(ctx context.Context, postId string, userIds []string) error {
	if len(userIds) == 0 {
		return nil
	}

	query := "INSERT INTO user_liked_posts (post_id, user_id) VALUES "
	values := []interface{}{postId}
	valueStrings := []string{}

	for i, userId := range userIds {
		valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, userId)
	}

	query += strings.Join(valueStrings, ", ") + " ON CONFLICT DO NOTHING"

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) DeletePostLikedBy(ctx context.Context, postId string, userIds []string) error {
	if len(userIds) == 0 {
		return nil
	}

	query := "DELETE FROM user_liked_posts WHERE post_id = $1 AND user_id = ANY($2)"
	_, err := p.db.ExecContext(ctx, query, postId, pq.Array(userIds))
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetPostTags(ctx context.Context, postId string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT tag FROM post_tags WHERE post_id = $1", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (p *PostgresDatabaseService) AddPostTags(ctx context.Context, postId string, tags []string) error {
	if len(tags) == 0 {
		return nil
	}
	query := "INSERT INTO post_tags (post_id, tag) VALUES "
	values := []interface{}{postId}
	valueStrings := []string{}

	for i, tag := range tags {
		valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, tag)
	}

	query += strings.Join(valueStrings, ", ") + " ON CONFLICT DO NOTHING"

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

// 如果数据库中存在大量的标签，且你担心一次性操作可能会引起事务锁或者其他性能问题，可以将操作批量化，即每次删除固定数量的标签
func (p *PostgresDatabaseService) DeletePostTags(ctx context.Context, postId string, tags []string) error {
	query := "DELETE FROM post_tags WHERE post_id = $1 AND tag = ANY($2)"
	_, err := p.db.ExecContext(ctx, query, postId, pq.Array(tags))
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetPostPictureLinks(ctx context.Context, postId string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT link FROM post_picture_links WHERE post_id = $1", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []string
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

func (p *PostgresDatabaseService) AddPostPictureLinks(ctx context.Context, postId string, pictureLinks []string) error {
	if len(pictureLinks) == 0 {
		return nil
	}

	query := "INSERT INTO post_picture_links (post_id, link) VALUES "
	values := []interface{}{postId}
	valueStrings := []string{}

	for i, link := range pictureLinks {
		valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, link)
	}

	query += strings.Join(valueStrings, ", ") + " ON CONFLICT DO NOTHING"

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) DeletePostPictureLinks(ctx context.Context, postId string, pictureLinks []string) error {
	if len(pictureLinks) == 0 {
		return nil
	}

	query := "DELETE FROM post_picture_links WHERE post_id = $1 AND link = ANY($2)"
	_, err := p.db.ExecContext(ctx, query, postId, pq.Array(pictureLinks))
	if err != nil {
		return err
	}

	return nil
}

// check codes below
func (p *PostgresDatabaseService) GetPostBookmarkedBy(ctx context.Context, postId string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT user_id FROM user_bookmarked_posts WHERE post_id = $1", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var user string
		if err := rows.Scan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (p *PostgresDatabaseService) AddPostBookmarkedBy(ctx context.Context, postId string, userIds []string) error {
	if len(userIds) == 0 {
		return nil
	}

	query := "INSERT INTO user_bookmarked_posts (post_id, user_id) VALUES "
	values := []interface{}{postId}
	valueStrings := []string{}

	for i, userId := range userIds {
		valueStrings = append(valueStrings, fmt.Sprintf("($1, $%d)", i+2))
		values = append(values, userId)
	}

	query += strings.Join(valueStrings, ", ") + " ON CONFLICT DO NOTHING"

	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) DeletePostBookmarkedBy(ctx context.Context, postId string, userIds []string) error {
	if len(userIds) == 0 {
		return nil
	}

	query := "DELETE FROM user_bookmarked_posts WHERE post_id = $1 AND user_id = ANY($2)"
	_, err := p.db.ExecContext(ctx, query, postId, pq.Array(userIds))
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabaseService) GetUserBookmarkedPosts(ctx context.Context, userId string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT post_id FROM user_bookmarked_posts WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []string
	for rows.Next() {
		var post string
		if err := rows.Scan(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (p *PostgresDatabaseService) GetUserLikedPosts(ctx context.Context, userId string) ([]string, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT post_id FROM user_liked_posts WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []string
	for rows.Next() {
		var post string
		if err := rows.Scan(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
