package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	config "github.com/ranckor96/gator/internal/config"
	database "github.com/ranckor96/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	programState := state{dbQueries, &cfg}

	cmds := commands{make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	input := os.Args
	if len(input) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{input[1], input[2:]}

	err = cmds.run(&programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
