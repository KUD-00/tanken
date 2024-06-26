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

func (r *PostRedisCacheService) GetPost(ctx context.Context, postId string) (*types.Post, error) {
	exists, err := r.IsKeyExist(ctx, "post:"+postId)
	if err != nil {
		return nil, fmt.Errorf("error checking key existence: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("post not found in cache")
	}

	ctx, pipe := r.NewPipe(ctx)

	postDetailsCmd, _ := r.GetPostDetailsCmd(ctx, postId)
	tagsCmd, _ := r.GetPostTagsCmd(ctx, postId)
	pictureLinksCmd, _ := r.GetPostPictureLinksCmd(ctx, postId)
	commentsCmd, _ := r.GetPostCommentIdsCmd(ctx, postId)
	likedByCmd, _ := r.GetPostLikedByCmd(ctx, postId)

	_, err = pipe.Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("error getting post from Redis: %v", err)
	}

	post := postDetailsCmd.Val()
	tags := tagsCmd.Val()
	pictureLinks := pictureLinksCmd.Val()
	comments := commentsCmd.Val()
	likedBy := likedByCmd.Val()

	return &types.Post{
		//TODO: make this as function
		PostDetails: types.PostDetails{
			PostId:    postId,
			CreatedAt: utils.StringToInt64(post["CreatedAt"], 0),
			UpdatedAt: utils.StringToInt64(post["UpdatedAt"], 0),
			UserId:    post["UserId"],
			Content:   post["Content"],
			Likes:     utils.StringToInt64(post["Likes"], 0),
			Bookmarks: utils.StringToInt64(post["Bookmarks"], 0),
		},
		PostSets: types.PostSets{
			Tags:         tags,
			PictureLinks: pictureLinks,
			CommentIds:   comments,
			LikedBy:      likedBy,
		},
	}, nil
}

func (r *PostRedisCacheService) RemovePost(ctx context.Context, postId string) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		ctx, pipe = r.NewPipe(ctx)
		r.RemovePostDetails(ctx, postId)
		r.RemovePostTags(ctx, postId, nil)
		r.RemovePostPictureLinks(ctx, postId, nil)
		r.RemoveComments(ctx, nil)
		r.RemovePostCommentIds(ctx, postId, nil)
		r.RemovePostLikedBy(ctx, postId, nil)

		_, err := pipe.Exec(ctx)
		if err != nil {
			return fmt.Errorf("error removing post from cache: %v", err)
		}
	} else {
		r.RemovePostDetails(ctx, postId)
		r.RemovePostTags(ctx, postId, nil)
		r.RemovePostPictureLinks(ctx, postId, nil)
		r.RemoveComments(ctx, nil)
		r.RemovePostCommentIds(ctx, postId, nil)
		r.RemovePostLikedBy(ctx, postId, nil)
	}

	return nil
}

func (r *PostRedisCacheService) GetPostDetailsCmd(ctx context.Context, postId string) (*redis.MapStringStringCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post details: pipeliner not found")
	}

	cmd := pipe.HGetAll(ctx, utils.PostPrefix+postId)
	return cmd, nil
}

func (r *PostRedisCacheService) GetPostDetails(ctx context.Context, postId string) (*types.PostDetailsPtr, error) {
	result, err := r.client.HGetAll(ctx, postId).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting post details: %v", err)
	}

	return PostDetailsMapToPostDetailsPtr(result), nil
}

func (r *PostRedisCacheService) SetPostDetails(ctx context.Context, postId string, postDetails *types.PostDetailsPtr) error {
	postDetailsMap := generatePostDetailsMap(postDetails)

	pipe := r.GetPipe(ctx)
	if pipe == nil {
		err := r.client.HSet(ctx, utils.PostPrefix+postId, postDetailsMap).Err()

		if err != nil {
			return fmt.Errorf("error setting post details: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.PostPrefix+postId, postDetailsMap)
	}

	return nil
}

func (r *PostRedisCacheService) RemovePostDetails(ctx context.Context, postId string) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		err := r.client.Del(ctx, utils.PostPrefix+postId).Err()
		if err != nil {
			return fmt.Errorf("error removing post details: %v", err)
		}
	} else {
		pipe.Del(ctx, utils.PostPrefix+postId)
	}

	return nil
}

func (r *PostRedisCacheService) GetPostLikedByCmd(ctx context.Context, postId string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post likedBy: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postId+utils.LikedBySuffix), nil
}

func (r *PostRedisCacheService) GetPostLikedBy(ctx context.Context, postId string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postId+utils.LikedBySuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post likedBy: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostLikedBy(ctx context.Context, postId string, userIds []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(userIds))
	for i, v := range userIds {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postId+utils.LikedBySuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post likedBy: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postId+utils.LikedBySuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostLikedBy(ctx context.Context, postId string, userIds []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(userIds))
	for i, v := range userIds {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postId+utils.LikedBySuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post likedBy: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postId+utils.LikedBySuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) GetPostTagsCmd(ctx context.Context, postId string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post tags: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postId+utils.TagsSuffix), nil
}

func (r *PostRedisCacheService) GetPostTags(ctx context.Context, postId string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postId+utils.TagsSuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post tags: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostTags(ctx context.Context, postId string, tags []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(tags))
	for i, v := range tags {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postId+utils.TagsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post tags: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postId+utils.TagsSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostTags(ctx context.Context, postId string, tags []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(tags))
	for i, v := range tags {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postId+utils.TagsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post tags: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postId+utils.TagsSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) GetPostPictureLinksCmd(ctx context.Context, postId string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post pictureLinks: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postId+utils.PictureLinksSuffix), nil
}

func (r *PostRedisCacheService) GetPostPictureLinks(ctx context.Context, postId string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postId+utils.PictureLinksSuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post pictureLinks: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostPictureLinks(ctx context.Context, postId string, pictureLinks []string) error {
	if len(pictureLinks) == 0 {
		return nil
	}

	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(pictureLinks))
	for i, v := range pictureLinks {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postId+utils.PictureLinksSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post pictureLinks: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postId+utils.PictureLinksSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostPictureLinks(ctx context.Context, postId string, pictureLinks []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(pictureLinks))
	for i, v := range pictureLinks {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postId+utils.PictureLinksSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post pictureLinks: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postId+utils.PictureLinksSuffix, args...)
	}
	return nil
}

func (r *PostRedisCacheService) GetPostCommentIdsCmd(ctx context.Context, postId string) (*redis.StringSliceCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting post commentIds: pipeliner not found")
	}

	return pipe.SMembers(ctx, utils.PostPrefix+postId+utils.CommentIdsSuffix), nil
}

func (r *PostRedisCacheService) GetPostCommentIds(ctx context.Context, postId string) ([]string, error) {
	result, err := r.client.SMembers(ctx, utils.PostPrefix+postId+utils.CommentIdsSuffix).Result()

	if err != nil {
		return nil, fmt.Errorf("error getting post commentIds: %v", err)
	}

	return result, nil
}

func (r *PostRedisCacheService) AddPostCommentIds(ctx context.Context, postId string, commentIds []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(commentIds))
	for i, v := range commentIds {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SAdd(ctx, utils.PostPrefix+postId+utils.CommentIdsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error adding post commentIds: %v", err)
		}
	} else {
		pipe.SAdd(ctx, utils.PostPrefix+postId+utils.CommentIdsSuffix, commentIds)
	}
	return nil
}

func (r *PostRedisCacheService) RemovePostCommentIds(ctx context.Context, postId string, commentIds []string) error {
	pipe := r.GetPipe(ctx)

	args := make([]interface{}, len(commentIds))
	for i, v := range commentIds {
		args[i] = v
	}

	if pipe == nil {
		err := r.client.SRem(ctx, utils.PostPrefix+postId+utils.CommentIdsSuffix, args...).Err()
		if err != nil {
			return fmt.Errorf("error removing post commentIds: %v", err)
		}
	} else {
		pipe.SRem(ctx, utils.PostPrefix+postId+utils.CommentIdsSuffix, commentIds)
	}
	return nil
}

// TODO: this is not effective, maybe restructure the comment key, let it contains postId maybe
func (r *PostRedisCacheService) RemovePostComments(ctx context.Context, postId string) error {
	commentIds, err := r.GetPostCommentIds(ctx, postId)
	if err != nil {
		return fmt.Errorf("error getting commentIds: %v", err)
	}

	pipe := r.GetPipe(ctx)

	if pipe == nil {
		ctx, pipe = r.NewPipe(ctx)
		err = r.RemoveComments(ctx, commentIds)
		if err != nil {
			return fmt.Errorf("error removing comments: %v", err)
		}

		_, err = pipe.Exec(ctx)
	} else {
		err = r.RemoveComments(ctx, commentIds)
		if err != nil {
			return fmt.Errorf("error removing comments: %v", err)
		}
	}

	return nil
}

func (r *PostRedisCacheService) GetCommentCmd(ctx context.Context, commentId string) (*redis.MapStringStringCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting comment: pipeliner not found")
	}

	return pipe.HGetAll(ctx, utils.CommentPrefix+commentId), nil
}

func (r *PostRedisCacheService) GetComment(ctx context.Context, commentId string) (*types.Comment, error) {
	result, err := r.client.HGetAll(ctx, utils.CommentPrefix+commentId).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting comment: %v", err)
	}

	return CommentMapToComment(result), nil
}

func (r *PostRedisCacheService) SetComment(ctx context.Context, commentId string, comment *types.Comment) error {
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
		err := r.client.HSet(ctx, utils.CommentPrefix+commentId, commentMap).Err()
		if err != nil {
			return fmt.Errorf("error setting comment: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.CommentPrefix+commentId, commentMap)
	}

	return nil
}

func (r *PostRedisCacheService) RemoveComments(ctx context.Context, commentIds []string) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		for _, commentId := range commentIds {
			err := r.client.Del(ctx, utils.CommentPrefix+commentId).Err()
			if err != nil {
				return fmt.Errorf("error removing comment: %v", err)
			}
		}
	} else {
		for _, commentId := range commentIds {
			pipe.Del(ctx, utils.CommentPrefix+commentId)
		}
	}
	return nil
}

func (r *PostRedisCacheService) SetUser(ctx context.Context, userId string, user *types.User) error {
	userMap := map[string]interface{}{
		"userId":             user.UserId,
		"username":           user.Username,
		"profilePictureLink": user.ProfilePictureLink,
		"bio":                user.Bio,
		"subscribed":         user.Subscribed,
	}

	pipe := r.GetPipe(ctx)
	if pipe == nil {
		err := r.client.HSet(ctx, utils.UserPrefix+userId, userMap).Err()
		if err != nil {
			return fmt.Errorf("error setting user: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.UserPrefix+userId, userMap)
	}

	return nil
}

func (r *PostRedisCacheService) GetUserCmd(ctx context.Context, userId string) (*redis.MapStringStringCmd, error) {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		return nil, fmt.Errorf("error getting user: pipeliner not found")
	}

	return pipe.HGetAll(ctx, utils.UserPrefix+userId), nil
}

func (r *PostRedisCacheService) GetUser(ctx context.Context, userId string) (*types.User, error) {
	result, err := r.client.HGetAll(ctx, utils.UserPrefix+userId).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	return UserMapToUser(result), nil
}

func (r *PostRedisCacheService) RemoveUser(ctx context.Context, userId string) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		err := r.client.Del(ctx, utils.UserPrefix+userId).Err()
		if err != nil {
			return fmt.Errorf("error removing user: %v", err)
		}
	} else {
		pipe.Del(ctx, utils.UserPrefix+userId)
	}
	return nil
}

func (r *PostRedisCacheService) AddPostCacheScore(ctx context.Context, postId string, score int64) error {
	pipe := r.GetPipe(ctx)

	if pipe == nil {
		err := r.client.HSet(ctx, utils.PostCacheKeys.CacheScore, postId, score).Err()
		if err != nil {
			return fmt.Errorf("error adding post cache score: %v", err)
		}
	} else {
		pipe.HSet(ctx, utils.PostCacheKeys.CacheScore, postId, score)
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
