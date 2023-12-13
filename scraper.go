package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/danny-ell77/rssagg/internal/database"
	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			
			//  Worker
			go worker(db, wg, feed)
		}
		wg.Wait()
	}
}

func worker(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedsAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error updating feed")
		return
	}

	rssFeed, err := getRemoteFeeds(feed.Url)
	if err != nil {
		log.Println("Error fetching feed")
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		
		// Specific to Wagslane
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("Could parse timestamp for post in feed")
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Description: description,
			Url: item.Link,
			PublishedAt: t,
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key"){
				continue
			}
			log.Println("Failed to create Post:", err)
		}
	}

	log.Printf("Feed %v collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}