package utils

const (
	LikedBySuffix      = ":likedBy"
	TagsSuffix         = ":tags"
	PictureLinksSuffix = ":pictureLinks"
	CommentIdsSuffix   = ":commentIds"

	CommentPrefix = "comment:"
	UserPrefix    = "user:"
)

type PostCacheKeysType struct {
	UpdatedAt  string
	CreatedAt  string
	UserId     string
	Content    string
	Likes      string
	Bookmarks  string
	Status     string
	CacheScore string
}

var PostCacheKeys = PostCacheKeysType{
	UpdatedAt:  "UpdatedAt",
	CreatedAt:  "CreatedAt",
	UserId:     "UserId",
	Content:    "Content",
	Likes:      "Likes",
	Bookmarks:  "Bookmarks",
	Status:     "Status",
	CacheScore: "CacheScore",
}
