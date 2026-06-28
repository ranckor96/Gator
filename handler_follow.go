package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ranckor96/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feedURL := cmd.Args[0]
	feedFollow, err := followHelper(s, feedURL)
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

func followHelper(s *state, feedURL string) (database.CreateFeedFollowRow, error) {
	feed, err := s.db.GetFeed(context.Background(), feedURL)
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("couldn't find feed: %w", err)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("couldn't find user: %w", err)
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
