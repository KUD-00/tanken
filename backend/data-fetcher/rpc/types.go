package rpc

type Location struct {
	Latitude  float64
	Longitude float64
}

type Post struct {
	PostId   string
	Location Location

	PostDetails
	PostSets
}

type PostDetails struct {
	Timestamp int64
	UserId    string
	Content   string
	Likes     int64
	Bookmarks int64
	Deleted   bool
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
	Timestamp int64
	Likes     int64
	Deleted   bool
}

type User struct {
	UserId     string
	Username   string
	Bio        string
	Subscribed bool
}
