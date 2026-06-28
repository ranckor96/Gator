package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ranckor96/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	userUUID := user.ID

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	feedUUID := uuid.New()
	currTime := time.Now().UTC()

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        feedUUID,
		Name:      feedName,
		Url:       feedURL,
		UserID:    userUUID,
		CreatedAt: currTime,
		UpdatedAt: currTime,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	_, err = followHelper(s, feedURL)
	if err != nil {
		return err
	}

	fmt.Println("Feed created successfully")
	printFeed(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      	 %v\n", feed.ID)
	fmt.Printf(" * Name:    	 %v\n", feed.Name)
	fmt.Printf(" * URL:      	 %v\n", feed.Url)
	fmt.Printf(" * user ID:      %v\n", feed.UserID)
}
