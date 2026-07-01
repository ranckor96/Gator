package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ranckor96/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("wrong time format, use 's', 'm', 'h'")
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feedToFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	currTime := time.Now().UTC()
	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		UpdatedAt: currTime,
		ID:        feedToFetch.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't mark feed fetched: %w", err)
	}

	fmt.Printf("Fetching feed from %s:\n", feedToFetch.Url)
	feed, err := fetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	for i, item := range feed.Channel.Item {
		fmt.Printf("Title %d:\n", i+1)
		fmt.Println(item.Title)
	}
	fmt.Println()

	return nil
}
