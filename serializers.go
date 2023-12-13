package main

import (
	"time"

	"github.com/danny-ell77/rssagg/internal/database"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time	`json:"createdd_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Name      string	`json:"name"`
	APIKey	  string    `json:"api_key"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string 	`json:"name"`
	Url       string	`json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"crated_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID        uuid.UUID 	`json:"id"`
	CreatedAt time.Time 	`json:"crated_at"`
	UpdatedAt time.Time 	`json:"updated_at"`
	Title       string		`json:"title"`
	Description	*string		`json:"description"`
	PublishedAt time.Time	`json:"published_at"`
	Url         string		`json:"url"`
	FeedID      uuid.UUID	`json:"feed_id"`
}

func userSerializer(dbUser database.User) User {
	return User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name: dbUser.Name,
		APIKey: dbUser.ApiKey,
	}
}

func feedSerializer(dbFeed database.Feed) Feed {
	return Feed{
		ID: dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name: dbFeed.Name,
		Url: dbFeed.Url,
		UserID: dbFeed.UserID,
	}
}

func manyFeedsSerializer(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, feedSerializer(dbFeed))
	}
	return feeds
}

func feedFollowSerializer(dbFeedFollows database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID: dbFeedFollows.ID,
		CreatedAt: dbFeedFollows.CreatedAt,
		UpdatedAt: dbFeedFollows.UpdatedAt,
		UserID: dbFeedFollows.UserID,
		FeedID: dbFeedFollows.FeedID,
	}
}

func manyFeedFollowSerializer(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feed_follows := []FeedFollow{}
	for _, dbFeedFollow := range dbFeedFollows {
		feed_follows = append(feed_follows, feedFollowSerializer(dbFeedFollow))
	}
	return feed_follows
}

func postsSerializer(dbPost database.Post) Post {
	var description *string

	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}

	return Post {
		ID: dbPost.ID,
		CreatedAt: dbPost.CreatedAt,
		UpdatedAt: dbPost.UpdatedAt,
		PublishedAt: dbPost.PublishedAt,
		Title: dbPost.Title,
		Url: dbPost.Url,
		Description: description,
		FeedID: dbPost.FeedID,
	}
}

func manyPostsSerializer(dbPosts []database.Post) []Post {
	posts := []Post{} 
	for _, post := range dbPosts {
		posts = append(posts, postsSerializer(post))
	}
	return posts
}