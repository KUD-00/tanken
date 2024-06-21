package cache

import (
	"context"
	"fmt"
	"strconv"

	dbtype "tanken/backend/common/db"
	"tanken/backend/common/types"
	"tanken/backend/common/utils"

	"github.com/redis/go-redis/v9"
)

type PostRedisCacheService struct {
	*RedisBase
}

var _ PostCacheService = (*PostRedisCacheService)(nil)

func NewPostRedisCacheService(client *redis.Client) *PostRedisCacheService {
	return &PostRedisCacheService{
		RedisBase: NewRedisBase(client),
	}
}

func (r *PostRedisCacheService) GetPostDetailsCmd(ctx context.Context, postID string) (*redis.MapStringStringCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post details: pipeliner not found")
	}

	cmd := pipe.HGetAll(ctx, utils.PostPrefix+postID)
	return cmd, nil
}

func (r *PostRedisCacheService) GetPostDetails(ctx context.Context, postID string) (*types.PostDetailsPtr, error) {
	result, err := r.client.HGetAll(ctx, postID).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting post details: %v", err)
	}

	return PostDetailsMapToPostDetailsPtr(result), nil
}

func (r *PostRedisCacheService) SetPostDetails(ctx context.Context, postID string, postDetails *types.PostDetailsPtr) error {
	postDetailsMap := generatePostDetailsMap(postDetails)

	pipe := r.GetPipe(ctx)
	if pipe == nil {
		err := r.client.HSet(ctx, utils.PostPrefix+postID, postDetailsMap).Err()

		if err != nil {
			return fmt.Errorf("error setting post details: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.PostPrefix+postID, postDetailsMap)
	}

	return nil
}

func (r *PostRedisCacheService) GetPostLikedByCmd(ctx context.Context, postID string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post likedBy: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postID+utils.LikedBySuffix), nil
}

func (r *PostRedisCacheService) GetPostLikedBy(ctx context.Context, postID string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postID+utils.LikedBySuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post likedBy: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostLikedBy(ctx context.Context, postID string, userIDs []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(userIDs))
	for i, v := range userIDs {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postID+utils.LikedBySuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post likedBy: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postID+utils.LikedBySuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostLikedBy(ctx context.Context, postID string, userIDs []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(userIDs))
	for i, v := range userIDs {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postID+utils.LikedBySuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post likedBy: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postID+utils.LikedBySuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) GetPostTagsCmd(ctx context.Context, postID string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post tags: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postID+utils.TagsSuffix), nil
}

func (r *PostRedisCacheService) GetPostTags(ctx context.Context, postID string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postID+utils.TagsSuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post tags: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostTags(ctx context.Context, postID string, tags []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(tags))
	for i, v := range tags {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postID+utils.TagsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post tags: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postID+utils.TagsSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostTags(ctx context.Context, postID string, tags []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(tags))
	for i, v := range tags {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postID+utils.TagsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post tags: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postID+utils.TagsSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) GetPostPictureLinksCmd(ctx context.Context, postID string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post pictureLinks: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postID+utils.PictureLinksSuffix), nil
}

func (r *PostRedisCacheService) GetPostPictureLinks(ctx context.Context, postID string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postID+utils.PictureLinksSuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post pictureLinks: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostPictureLinks(ctx context.Context, postID string, pictureLinks []string) error {
	if len(pictureLinks) == 0 {
		return nil
	}

	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(pictureLinks))
	for i, v := range pictureLinks {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postID+utils.PictureLinksSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post pictureLinks: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postID+utils.PictureLinksSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostPictureLinks(ctx context.Context, postID string, pictureLinks []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(pictureLinks))
	for i, v := range pictureLinks {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postID+utils.PictureLinksSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post pictureLinks: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postID+utils.PictureLinksSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) GetPostCommentIdsCmd(ctx context.Context, postID string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post commentIds: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postID+utils.CommentIdsSuffix), nil
}

func (r *PostRedisCacheService) GetPostCommentIds(ctx context.Context, postID string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postID+utils.CommentIdsSuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post commentIds: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostCommentIds(ctx context.Context, postID string, commentIds []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(commentIds))
	for i, v := range commentIds {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postID+utils.CommentIdsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post commentIds: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postID+utils.CommentIdsSuffix, commentIds)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostCommentIds(ctx context.Context, postID string, commentIds []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(commentIds))
	for i, v := range commentIds {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postID+utils.CommentIdsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post commentIds: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postID+utils.CommentIdsSuffix, commentIds)
	}
	return nil
}

func (r *PostRedisCacheService) GetCommentCmd(ctx context.Context, commentID string) (*redis.MapStringStringCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting comment: pipeliner not found")
	}

	return pipe.HGetAll(ctx, utils.CommentPrefix+commentID), nil
}

func (r *PostRedisCacheService) GetComment(ctx context.Context, commentID string) (*types.Comment, error) {
	result, err := r.client.HGetAll(ctx, utils.CommentPrefix+commentID).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting comment: %v", err)
	}

	return CommentMapToComment(result), nil
}

func (r *PostRedisCacheService) SetComment(ctx context.Context, commentID string, comment *types.Comment) error {
	commentMap := map[string]interface{}{
		"createdAt": comment.CreatedAt,
		"updatedAt": comment.UpdatedAt,
		"postId":    comment.PostId,
		"userId":    comment.UserId,
		"content":   comment.Content,
		"likes":     comment.Likes,
	}

	pipe := r.GetPipe(ctx)
	if pipe == nil {
		err := r.client.HSet(ctx, utils.CommentPrefix+commentID, commentMap).Err()
		if err != nil {
			return fmt.Errorf("error setting comment: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.CommentPrefix+commentID, commentMap)
	}

	return nil
}

func (r *PostRedisCacheService) RemoveComment(ctx context.Context, commentID string) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		err := r.client.Del(ctx, utils.CommentPrefix+commentID).Err()
		if err != nil {
			return fmt.Errorf("error removing comment: %v", err)
		}
	} else {
		pipe.Del(ctx, utils.CommentPrefix+commentID)
	}
	return nil
}

func (r *PostRedisCacheService) SetUser(ctx context.Context, userID string, user *types.User) error {
	userMap := map[string]interface{}{
		"userId":             user.UserId,
		"username":           user.Username,
		"profilePictureLink": user.ProfilePictureLink,
		"bio":                user.Bio,
		"subscribed":         user.Subscribed,
	}

	pipe := r.GetPipe(ctx)
	if pipe == nil {
		err := r.client.HSet(ctx, utils.UserPrefix+userID, userMap).Err()
		if err != nil {
			return fmt.Errorf("error setting user: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.UserPrefix+userID, userMap)
	}

	return nil
}

func (r *PostRedisCacheService) GetUserCmd(ctx context.Context, userID string) (*redis.MapStringStringCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting user: pipeliner not found")
	}

	return pipe.HGetAll(ctx, utils.UserPrefix+userID), nil
}

func (r *PostRedisCacheService) GetUser(ctx context.Context, userID string) (*types.User, error) {
	result, err := r.client.HGetAll(ctx, utils.UserPrefix+userID).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	return UserMapToUser(result), nil
}

func (r *PostRedisCacheService) RemoveUser(ctx context.Context, userID string) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		err := r.client.Del(ctx, utils.UserPrefix+userID).Err()
		if err != nil {
			return fmt.Errorf("error removing user: %v", err)
		}
	} else {
		pipe.Del(ctx, utils.UserPrefix+userID)
	}
	return nil
}

func (r *PostRedisCacheService) AddPostCacheScore(ctx context.Context, postID string, score int64) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		err := r.client.HSet(ctx, utils.PostCacheKeys.CacheScore, postID, score).Err()
		if err != nil {
			return fmt.Errorf("error adding post cache score: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.PostCacheKeys.CacheScore, postID, score)
	}
	return nil
}

func (r *PostRedisCacheService) GetNonPopularPosts(ctx context.Context, limit int64) ([]types.Post, error) {
	return nil, nil
}

func (r *PostRedisCacheService) WriteBackToDB(ctx context.Context, db dbtype.DatabaseService, postIds []string) error {
	return nil
}

func PostDetailsMapToPostDetailsPtr(postDetailsMap map[string]string) *types.PostDetailsPtr {
	var createdAt *int64
	if val, ok := postDetailsMap["CreatedAt"]; ok {
		if intValue, err := strconv.ParseInt(val, 10, 64); err == nil {
			createdAt = &intValue
		}
	}

	var updateAt *int64
	if val, ok := postDetailsMap["UpdatedAt"]; ok {
		if intValue, err := strconv.ParseInt(val, 10, 64); err == nil {
			updateAt = &intValue
		}
	}

	var userId *string
	if val, ok := postDetailsMap["UserId"]; ok {
		userId = &val
	}

	var content *string
	if val, ok := postDetailsMap["Content"]; ok {
		content = &val
	}

	var likes *int64
	if val, ok := postDetailsMap["Likes"]; ok {
		if intValue, err := strconv.ParseInt(val, 10, 64); err == nil {
			likes = &intValue
		}
	}

	var bookmarks *int64
	if val, ok := postDetailsMap["Bookmarks"]; ok {
		if intValue, err := strconv.ParseInt(val, 10, 64); err == nil {
			bookmarks = &intValue
		}
	}

	var status *int64
	if val, ok := postDetailsMap["Status"]; ok {
		if intValue, err := strconv.ParseInt(val, 10, 64); err == nil {
			status = &intValue
		}
	}

	return &types.PostDetailsPtr{
		CreatedAt: createdAt,
		UpdatedAt: updateAt,
		UserId:    userId,
		Content:   content,
		Likes:     likes,
		Bookmarks: bookmarks,
		Status:    status,
	}
}

func CommentMapToComment(commentMap map[string]string) *types.Comment {
	return &types.Comment{
		UpdatedAt: utils.StringToInt64(commentMap["updatedAt"], 0),
		CreatedAt: utils.StringToInt64(commentMap["createdAt"], 0),
		PostId:    commentMap["postId"],
		UserId:    commentMap["userId"],
		Content:   commentMap["content"],
		Likes:     utils.StringToInt64(commentMap["likes"], 0),
	}
}

func UserMapToUser(userMap map[string]string) *types.User {
	return &types.User{
		UserId:             userMap["userId"],
		Username:           userMap["username"],
		ProfilePictureLink: userMap["profilePictureLink"],
		Bio:                userMap["bio"],
		Subscribed:         utils.StringToInt64(userMap["subscribed"], 0),
	}
}

func generatePostDetailsMap(postDetails *types.PostDetailsPtr) map[string]interface{} {
	postDetailsMap := make(map[string]interface{})
	if postDetails.UpdatedAt != nil {
		postDetailsMap[utils.PostCacheKeys.UpdatedAt] = strconv.FormatInt(*postDetails.UpdatedAt, 10)
	}
	if postDetails.CreatedAt != nil {
		postDetailsMap[utils.PostCacheKeys.CreatedAt] = strconv.FormatInt(*postDetails.CreatedAt, 10)
	}
	if postDetails.UserId != nil {
		postDetailsMap[utils.PostCacheKeys.UserId] = *postDetails.UserId
	}
	if postDetails.Content != nil {
		postDetailsMap[utils.PostCacheKeys.Content] = *postDetails.Content
	}
	if postDetails.Likes != nil {
		postDetailsMap[utils.PostCacheKeys.Likes] = strconv.FormatInt(*postDetails.Likes, 10)
	}
	if postDetails.Bookmarks != nil {
		postDetailsMap[utils.PostCacheKeys.Bookmarks] = strconv.FormatInt(*postDetails.Bookmarks, 10)
	}
	if postDetails.Status != nil {
		postDetailsMap[utils.PostCacheKeys.Status] = strconv.FormatInt(*postDetails.Status, 10)
	}
	return postDetailsMap
}
