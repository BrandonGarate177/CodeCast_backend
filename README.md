# CodeCast Backend

CodeCast Backend is a real-time code sharing and collaboration platform designed for developers and educators. It enables anonymous code sharing sessions where participants can join, view live code snippets, and collaborate without requiring user accounts, making it perfect for quick demos, code reviews, and educational purposes.

## Key Features

- **Anonymous Sessions**: Create and join code sharing sessions without user registration
- **Real-time Code Sharing**: Push and pull code snippets in real-time during sessions
- **Session Management**: Create, join, and end sessions with unique session codes
- **Participant Management**: Support for up to 20 participants per session
- **RESTful API**: Clean and well-documented API endpoints for all functionality
- **Database Integration**: PostgreSQL integration with proper data modeling

## Technical Stack

- **Backend Framework**: Go 1.23.0 with Gin web framework
- **Database**: PostgreSQL with lib/pq driver
- **Authentication**: Token-based authentication for session creators
- **Database Migration**: Atlas for database schema management
- **Environment Management**: godotenv for configuration
- **UUID Generation**: Google UUID library for unique identifiers

## API Endpoints

### Anonymous Sessions

| Endpoint | Method | Authentication | Purpose |
|----------|--------|----------------|---------|
| `/api/v1/sessions/anon` | `POST` | None | Create a new anonymous session |
| `/api/v1/sessions/anon/:code/join` | `POST` | None | Join an existing session as participant |
| `/api/v1/sessions/anon/:code/snippets` | `POST` | Creator token | Push a code snippet to session |
| `/api/v1/sessions/anon/:code/snippets` | `GET` | None | Fetch the latest code snippet |
| `/api/v1/sessions/anon/:code/end` | `POST` | Creator token | End the session and clean up data |

## Installation & Usage

### Prerequisites

- Go 1.23.0 or later
- PostgreSQL database
- Environment variables configured

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd CodeCast_backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   Create a `.env` file in the parent directory with:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=codecast
   ```

4. **Set up the database**
   ```bash
   # Using Atlas for migrations
   atlas migrate apply --env local
   ```

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

The server will start on `http://localhost:8080`

## Database Schema

### anon_sessions
- `id` (UUID, Primary Key)
- `session_code` (TEXT, Unique)
- `creator_token` (TEXT, Secret)
- `display_name` (TEXT)
- `is_active` (BOOLEAN)
- `created_at` (TIMESTAMP)

### anon_participants
- `id` (UUID, Primary Key)
- `session_id` (UUID, Foreign Key)
- `display_name` (TEXT)
- `joined_at` (TIMESTAMP)

### anon_snippets
- `id` (UUID, Primary Key)
- `session_id` (UUID, Foreign Key)
- `file_name` (TEXT)
- `content` (TEXT)
- `pushed_at` (TIMESTAMP)


### Key Components

- **Main Application** (`cmd/main.go`): Sets up the Gin router, initializes database, and registers API routes
- **Database Layer** (`db/db.go`): Handles PostgreSQL connection and initialization
- **Anonymous Sessions Module** (`modules/anon_sessions/`): Complete module for session management
- **Utilities** (`utils/`): Common utilities for authentication and helper functions

### Development Workflow

1. Make changes to the codebase
2. Test locally with `go run cmd/main.go`
3. Update database schema using Atlas migrations if needed
4. Test API endpoints using tools like Postman or curl

## License

All packaged builds are © 2025 Brandon Garate.
Distributed binaries are for personal or team use only. Please see the main CodeCast repository for source code licensing terms.

## Support

For issues, feature requests, or contributions, please visit the main CodeCast repository or contact the development team.
