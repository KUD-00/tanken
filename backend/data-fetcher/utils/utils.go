package utils

import (
	"tanken/backend/common/types"
	pb "tanken/backend/data-fetcher/rpc/pb"
)

func CommonPostToPbPost(commonPost *types.Post) *pb.Post {
	return &pb.Post{
		PostId: commonPost.PostId,
		Location: &pb.Location{
			Latitude:  commonPost.Location.Latitude,
			Longitude: commonPost.Location.Longitude,
		},
		CreatedAt: commonPost.CreatedAt,
		UpdatedAt: commonPost.UpdatedAt,
		Content:   commonPost.Content,
		Likes:     commonPost.Likes,
		Bookmarks: commonPost.Bookmarks,
		Tags:      commonPost.Tags,
	}
}

func CommonPostsToPbPosts(commonPosts []*types.Post) []*pb.Post {
	var pbPosts []*pb.Post
	for _, post := range commonPosts {
		pbPosts = append(pbPosts, CommonPostToPbPost(post))
	}
	return pbPosts
}

func CommonCommentToPbComment(commonComment *types.Comment, commonUser *types.User) *pb.Comment {
	user := CommonUserToPbUser(commonUser)
	return &pb.Comment{
		CommentId: commonComment.CommentId,
		PostId:    commonComment.PostId,
		User:      user,
		Content:   commonComment.Content,
		CreatedAt: commonComment.CreatedAt,
		UpdatedAt: commonComment.UpdatedAt,
		Likes:     commonComment.Likes,
	}
}

func CommonUserToPbUser(commonUser *types.User) *pb.User {
	return &pb.User{
		UserId:             commonUser.UserId,
		UserName:           commonUser.Username,
		Bio:                commonUser.Bio,
		Subscribed:         commonUser.Subscribed,
		ProfilePictureLink: commonUser.ProfilePictureLink,
	}
}
