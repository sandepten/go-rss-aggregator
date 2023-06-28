package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sandepten/go-rss-aggregator/internal/database"
)

func startScraping(
	db *database.Queries, // database connection
	concurrency int, // number of goroutines to use
	timeBetweenScrapes time.Duration, // time in seconds between scrapes
) {
	log.Printf("Starting scraping with %d goroutines and %v seconds between scrapes", concurrency, timeBetweenScrapes)

	ticker := time.NewTicker(timeBetweenScrapes)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error getting feeds to fetch: %v", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go func(feed database.Feed) {
				defer wg.Done()

				_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
				if err != nil {
					log.Printf("Error marking feed as fetched: %v", err)
					return
				}

				rssFeed, err := urlToFeed(feed.Url)
				if err != nil {
					log.Printf("Error fetching feed: %v", err)
					return
				}

				for _, item := range rssFeed.Channel.Item {
					description := sql.NullString{}
					if item.Description != "" {
						description.String = item.Description
						description.Valid = true
					}

					pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
					if err != nil {
						log.Printf("Error parsing date: %v", err)
						continue
					}
					_, err = db.CreatePost(context.Background(), database.CreatePostParams{
						ID:          uuid.New(),
						FeedID:      feed.ID,
						Title:       item.Title,
						Url:         item.Link,
						Description: description,
						PublishedAt: pubAt,
						CreatedAt:   time.Now().UTC(),
						UpdatedAt:   time.Now().UTC(),
					})
					if err != nil {
						if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
							continue
						}
						log.Printf("Error creating item: %v", err)
						continue
					}
				}
				log.Printf("Feed %s collected, %d items", feed.Url, len(rssFeed.Channel.Item))
			}(feed)
		}
		wg.Wait()
	}
}
