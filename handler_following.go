package main

import (
	"context"
	"fmt"

	"github.com/ranckor96/gator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get user feed follows: %w", err)
	}

	for i, feedFollow := range feedFollows {
		fmt.Printf("Feed follow №%d: %s\n", i, feedFollow.FeedName)
	}

	return nil
}
