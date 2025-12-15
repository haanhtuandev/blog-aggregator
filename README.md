# Blog Aggregator

A command-line RSS feed aggregator written in Go that allows users to subscribe to blogs and news sources, aggregate their content, and browse posts from multiple sources in one place.

## Features

- **User management**: Create accounts and login/logout functionality
- **RSS feed subscription**: Add and manage RSS feeds from your favorite blogs
- **Feed following**: Follow/unfollow feeds to customize your content stream
- **Automated scraping**: Regularly fetch new posts from subscribed feeds
- **Content browsing**: View latest posts from all followed feeds
- **Post storage**: Persistently store retrieved posts in a PostgreSQL database

## Prerequisites

- Go 1.24.1 or higher
- PostgreSQL database
- `pq` PostgreSQL driver (installed automatically via go modules)
- `goose` for database migrations (optional, for schema management)

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd blog-aggregator
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the PostgreSQL database and update the connection string in `.gatorconfig.json`:
   ```json
   {
       "db_url": "postgres://username:password@localhost:5432/blog_aggregator?sslmode=disable"
   }
   ```

4. Run database migrations:
   ```bash
   # Using the provided migration script
   ./migration.sh
   ```
   Or manually create the tables:
   - users
   - feeds
   - feed_follows
   - posts

## Configuration

The application looks for a configuration file named `.gatorconfig.json` in your home directory with the following structure:

```json
{
    "db_url": "postgres://username:password@localhost:5432/blog_aggregator?sslmode=disable",
    "current_user_name": ""
}
```

## Available Commands

The application provides several commands accessible through the CLI:

### User Management
- `register <username>` - Create a new user account
- `login <username>` - Login as an existing user
- `users` - List all registered users

### Feed Management
- `addfeed <name> <url>` - Add a new feed to your subscriptions (requires authentication)
- `feeds` - List all available feeds from all users
- `follow <feed_url>` - Subscribe to a feed (requires authentication)
- `unfollow <feed_url>` - Unsubscribe from a feed (requires authentication)
- `following` - List all feeds you're currently following (requires authentication)

### Content Aggregation
- `agg <interval>` - Start fetching posts from followed feeds at specified intervals (e.g., "10m", "1h")
- `browse [limit]` - View latest posts from all followed feeds (defaults to 2 if no limit provided)

### System
- `reset` - Reset user data in the database

## Usage Examples

1. **Register a new user:**
   ```bash
   go run main.go register alice
   ```

2. **Login as a user:**
   ```bash
   go run main.go login alice
   ```

3. **Add a new RSS feed:**
   ```bash
   go run main.go addfeed "Tech Blog" "https://example.com/feed.xml"
   ```

4. **Follow an existing feed:**
   ```bash
   go run main.go follow "https://example.com/feed.xml"
   ```

5. **Start aggregating feeds (every 5 minutes):**
   ```bash
   go run main.go agg 5m
   ```

6. **Browse your followed posts (latest 5):**
   ```bash
   go run main.go browse 5
   ```

## Architecture

The application uses a layered architecture:

- **Main Layer**: Entry point and command routing
- **Command Layer**: Handles individual CLI commands
- **Database Layer**: Interacts with PostgreSQL using generated SQL queries
- **Config Layer**: Manages persistent configuration
- **Parsing Layer**: Handles RSS feed parsing and HTTP requests

### Database Schema

The application uses the following tables:

- `users`: Stores user information
- `feeds`: Contains information about RSS feeds
- `feed_follows`: Tracks which users follow which feeds
- `posts`: Stores individual blog posts retrieved from feeds

## Development

If you want to extend the application, note that database queries are generated using `sqlc`. To regenerate the SQL interface files after making schema changes:

1. Install `sqlc`:
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

2. Generate code:
   ```bash
   sqlc generate
   ```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

If you have questions or suggestions, feel free to open an issue in the repository.

---

Built with Go!