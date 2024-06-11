-- V1__initial_schema.sql

CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    profile_picture_link TEXT,
    bio TEXT,
    oauth_provider TEXT NOT NULL,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    subscribed INT DEFAULT 0
);

CREATE TABLE posts (
    post_id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(user_id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    likes INT DEFAULT 0,
    latitude FLOAT,
    longitude FLOAT,
    status SMALLINT DEFAULT 1
);

CREATE TABLE comments (
    comment_id UUID PRIMARY KEY,
    post_id UUID REFERENCES posts(post_id),
    user_id UUID REFERENCES users(user_id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    likes INT DEFAULT 0
);

CREATE TABLE post_tags (
    post_id UUID NOT NULL,
    tag TEXT NOT NULL,
    PRIMARY KEY (post_id, tag),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE TABLE post_pictures (
    post_id UUID NOT NULL,
    picture_link TEXT NOT NULL,
    PRIMARY KEY (post_id, picture_link),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE TABLE post_comments (
    post_id UUID NOT NULL,
    comment_id UUID NOT NULL,
    PRIMARY KEY (post_id, comment_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE
);

CREATE TABLE post_likes (
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE user_checked_posts (
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    PRIMARY KEY (user_id, post_id),
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE TABLE user_liked_posts (
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);

CREATE TABLE user_bookmarked_posts (
    user_id UUID NOT NULL,
    post_id UUID NOT NULL,
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE
);
