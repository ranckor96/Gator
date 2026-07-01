package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ranckor96/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <titime_between_reqsme>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	feedToFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Println("couldn't get next feed to fetch: %w", err)
		return
	}

	currTime := time.Now().UTC()
	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		UpdatedAt: currTime,
		ID:        feedToFetch.ID,
	})
	if err != nil {
		log.Println("couldn't mark feed fetched: %w", err)
		return
	}

	fmt.Printf("Fetching feed from %s:\n", feedToFetch.Url)
	feed, err := fetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		log.Println("couldn't fetch feed: %w", err)
		return
	}

	for _, item := range feed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   currTime,
			UpdatedAt:   currTime,
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publishedAt,
			FeedID:      feedToFetch.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	fmt.Println("Feed fetched")
	log.Printf("Feed %s collected, %v posts found", feedToFetch.Name, len(feed.Channel.Item))
}
