package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	// db_state := state{db: dbQueries, cfg: &config.Config{Db_url: db_url, Current_user_name: "tuanha"}}
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		log.Fatalf("error opening db with url: %v", err)
	}
	dbQueries := database.New(db)
	programState := &state{db: dbQueries, cfg: &cfg}

	// initialize a struct of commands
	cmds := Commands{Map: make(map[string]func(*state, Command) error)}

	// register commands and check length
	cmds.Register("login", handlerLogin)
	cmds.Register("register", handlerRegister)
	cmds.Register("reset", handlerReset)
	cmds.Register("users", handlerUsers)
	cmds.Register("agg", handlerAgg)
	cmds.Register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.Register("feeds", handlerFeeds)
	cmds.Register("follow", middlewareLoggedIn(handlerFollow))
	cmds.Register("following", middlewareLoggedIn(handlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.Register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	err = cmds.Run(programState, Command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}

}
