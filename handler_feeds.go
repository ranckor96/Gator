package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list feeds: %w", err)
	}

	for i, feed := range feeds {
		fmt.Printf("Feed №%d\n", i+1)
		fmt.Printf(" * Name:		%v\n", feed.Name)
		fmt.Printf(" * URL:			%v\n", feed.Url)
		fmt.Printf(" * user name:	%v\n", feed.UserName)
		fmt.Print("\n")
	}

	return nil
}
