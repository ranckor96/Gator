package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ranckor96/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feedURL := cmd.Args[0]
	feedFollow, err := feedFollow(s, feedURL, user)
	if err != nil {
		return err
	}

	fmt.Println("Feed followed successfully")
	printFeedFollow(feedFollow)
	return nil
}

func printFeedFollow(feedFollow database.CreateFeedFollowRow) {
	fmt.Printf(" * Feed Name:	%v\n", feedFollow.FeedName)
	fmt.Printf(" * User Name:	%v\n", feedFollow.UserName)
}

func feedFollow(s *state, feedURL string, user database.User) (database.CreateFeedFollowRow, error) {
	feed, err := s.db.GetFeed(context.Background(), feedURL)
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("couldn't find feed: %w", err)
	}

	feedFollowID := uuid.New()
	currTime := time.Now().UTC()

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        feedFollowID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("couldn't follow feed: %w", err)
	}

	return feedFollow, nil
}
