package types

type Location struct {
	Latitude  float64
	Longitude float64
}

type Post struct {
	PostDetails
	PostSets
}

type PostDetails struct {
	PostId     string
	CreatedAt  int64
	UpdatedAt  int64
	UserId     string
	Content    string
	Likes      int64
	Bookmarks  int64
	Status     int64
	Location   Location
	CacheScore int64
}

type PostDetailsPtr struct {
	CreatedAt  *int64
	UpdatedAt  *int64
	UserId     *string
	Content    *string
	Likes      *int64
	Bookmarks  *int64
	Status     *int64
	CacheScore *int64
	Changed    *bool
}

type PostSets struct {
	Tags         []string
	PictureLinks []string
	CommentIds   []string
	LikedBy      []string
}

type Comment struct {
	CommentId string
	PostId    string
	UserId    string
	Content   string
	CreatedAt int64
	UpdatedAt int64
	Likes     int64
	Status    int64
}

type User struct {
	UserId             string
	Username           string
	Email              string
	Bio                string
	CreatedAt          int64
	Subscribed         int64
	ProfilePictureLink string
	OauthProvider      string
}

type UserPtr struct {
	// for edit
	Username           *string
	Email              *string
	Bio                *string
	Subscribed         *int64
	ProfilePictureLink *string
	OauthProvider      *string
	Changed            *bool
}

type UserSets struct {
	CheckedPostIds    []string
	LikedPostIds      []string
	BookmarkedPostIds []string
}
