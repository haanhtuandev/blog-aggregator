package main

import (
	"blog-aggregator/internal/database"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}
type Commands struct {
	Map         map[string]func(*state, Command) error
	CommandList []Command
}

func handlerLogin(s *state, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("no argument found at check length")
	}
	username := cmd.Args[0]
	user, _ := s.db.GetUser(context.Background(), username)
	if user == (database.User{}) {
		return errors.New("user not found, try registering first")
	}

	err := s.cfg.SetUser(username)
	if err != nil {
		return errors.New("no argument found")
	}

	fmt.Println("User has been set!")

	return nil
}

func handlerRegister(s *state, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("wrong arguments passing")
	}
	fmt.Println(cmd.Args[0])
	username := cmd.Args[0]
	user, _ := s.db.GetUser(context.Background(), username)
	if user != (database.User{}) {
		return errors.New("error creating new user, user already exists")
	}

	new_user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: username})
	if err != nil {
		return err
	}
	s.cfg.SetUser(username)
	fmt.Printf("User %s created! Your info: %v %v", new_user.Name, new_user.ID, new_user.CreatedAt)
	return nil
}

func handlerUsers(s *state, cmd Command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		if s.cfg.Current_user_name == u.Name {
			fmt.Printf("* %v (current)\n", u.Name)
		} else {
			fmt.Printf("* %v\n", u.Name)
		}

	}

	return nil
}

func handlerReset(s *state, cmd Command) error {
	err := s.db.ResetUserDB(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerAgg(s *state, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("need to specify the req rate")
	}

	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(duration)
	fmt.Printf("Collecting feeds every %v \n", cmd.Args[0])
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed(s *state, cmd Command, user database.User) error {
	if len(cmd.Args) != 2 {
		return errors.New("wrong argument")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	// user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	// if err != nil {
	// 	return err
	// }
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name, Url: url, UserID: user.ID})
	if err != nil {
		return err
	}
	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)

	_, err = s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd Command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("Feed %v | from %v | belongs to %v \n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func handlerFollow(s *state, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("wrong argument passed")
	}
	url := cmd.Args[0]
	// user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
	// if err != nil {
	// 	return err
	// }
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	feed_user_table, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})
	if err != nil {
		return err
	}
	for _, entry := range feed_user_table {
		fmt.Printf("%v follows %v \n", entry.UserName, entry.FeedName)
	}

	return nil
}

func handlerFollowing(s *state, cmd Command, user database.User) error {
	// current_user := s.cfg.Current_user_name
	user_feed, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}
	for _, feed := range user_feed {
		fmt.Printf("%v - %v \n", feed.FeedName, feed.UserName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return errors.New("wrong arguments")
	}
	url := cmd.Args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	s.db.DeleteFollow(context.Background(), database.DeleteFollowParams{feed.ID, user.ID})
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd Command, user database.User) error) func(*state, Command) error {
	return func(s *state, cmd Command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), next_feed.ID)
	if err != nil {
		return err
	}
	rss_feed, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return err
	}
	feed_id := next_feed.ID

	for _, item := range rss_feed.Channel.Item {
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: item.PubDate,
			FeedID:      feed_id,
		})
		if err != nil {
			return err
		}
		fmt.Println(item.Title)
		fmt.Println(item.Description)
		fmt.Println(item.PubDate)
	}
	return nil

}

func handlerBrowse(s *state, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		if len(cmd.Args) == 0 {
			err := helperBrowse(2, user, s)
			if err != nil {
				return err
			}
			return nil
		} else {
			return errors.New("too many arguments")
		}
	}
	limit, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return err
	}
	err = helperBrowse(limit, user, s)
	if err != nil {
		return err
	}
	return nil

}
func helperBrowse(limit int, user database.User, s *state) error {
	posts, err := s.db.GetPostForUser(context.Background(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Println(post.Name)
		fmt.Println(post.Description)
		fmt.Println(post.CreatedAt)
		fmt.Println(post.PublishedAt)
	}
	return nil
}
func (c *Commands) Run(s *state, cmd Command) error {
	cmd_function, ok := c.Map[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return cmd_function(s, cmd)
}

func (c *Commands) Register(name string, f func(*state, Command) error) {
	_, ok := c.Map[name]
	if !ok {
		c.Map[name] = f
	}
}
