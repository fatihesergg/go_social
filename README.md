# Go Social API

Go Social is a backend API for a social media application built with Go. It provides a robust foundation for features like user authentication, content creation, and social interactions. The project follows modern Go development practices, emphasizing a clean, layered architecture and security.

## ✨ Features

### Core Functionality

- **User Management**: Secure user signup and login.
- **JWT Authentication**: Endpoints are protected using JSON Web Tokens.
- **User Profiles**: Fetch user data and profile information.
- **Social Graph**: Users can follow and unfollow each other.
- **Follower/Following Lists**: View lists of who a user follows and who follows them.
- **Post Management**: Full CRUD (Create, Read, Update, Delete) operations for posts.
- **Comment System**: Full CRUD operations for comments on posts.
- **Personalized Feed**: A user-specific feed that aggregates posts from the users they follow.

### Architecture & Design

- **Layered Architecture**: The project is structured into distinct layers (Controller, Database/Store, Model) to enforce separation of concerns.
  - `internal/controller`: Handles incoming HTTP requests, validates input, and calls the appropriate services.
  - `internal/database`: Implements the Repository Pattern. The `...Store` structs abstract all database interactions, making the business logic independent of the database implementation.
  - `internal/model`: Defines the core data structures of the application.
- **Dependency Injection**: Dependencies (like database stores) are created in `main.go` and injected into the controllers. This promotes loose coupling and makes components highly testable.
- **Centralized Routing**: All API routes are clearly defined in `cmd/go_social/main.go` using the Gin Gonic framework, providing a single source of truth for the API's structure.
- **Middleware**: Authentication is handled cleanly using Gin middleware (`internal/middleware/auth.go`), which intercepts requests to protected routes and validates the JWT.

### Security

- **JWT-Based Authentication**: Stateless authentication is implemented using JWTs, which are issued upon successful login.
- **Password Hashing**: The database schema is designed to work with PostgreSQL's `pgcrypto` extension, ensuring that user passwords are securely hashed and never stored in plain text.
- **Authorization Logic**: Operations like updating or deleting a post/comment include checks to ensure that the request is made by the authorized owner of the resource.

### Database

- **PostgreSQL**: A powerful, open-source relational database is used for data persistence.
- **Database Migrations**: The `golang-migrate` tool is used to manage database schema changes in a version-controlled, systematic way. This ensures the database schema is always in sync with the application code.

## 🛠️ Technologies Used

- **Backend**: Go
- **Framework**: Gin Gonic
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Database Migrations**: `golang-migrate`
- **Environment Management**: `godotenv`

## 🚀 Getting Started

### Prerequisites

- Go (1.21 or later)
- PostgreSQL
- `golang-migrate` CLI
- A running PostgreSQL instance.

### Installation & Setup

1.  **Clone the repository:**

    ```bash
    git clone github.com/fatihesergg/go_social
    cd go_social
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Configure Environment Variables:**
    Create a `.env` file in the root directory and add the following variables:

    ```env
    DATABASE_URI="postgres://user:password@localhost:5432/go_social?sslmode=disable"
    JWT_SECRET="your-super-secret-key"
    ```

4.  **Run Database Migrations:**
    Apply all pending migrations to set up your database schema.

    ```bash
    migrate -path ./internal/migration -database "$DATABASE_URI" up
    ```

5.  **Run the Application:**
    ```bash
    go run ./cmd/go_social/main.go
    ```
    The server will start, typically on port `3000`.

## API Endpoints

All endpoints are prefixed with `/api/v1`.

| Method   | Endpoint               | Description                    | Authentication |
| :------- | :--------------------- | :----------------------------- | :------------- |
| `POST`   | `/signup`              | Register a new user            | None           |
| `POST`   | `/login`               | Log in a user and get a JWT    | None           |
| `GET`    | `/users/getMe`         | Get the current user's profile | Required       |
| `GET`    | `/users/:id`           | Get a user's profile by ID     | Required       |
| `POST`   | `/users/:id/follow`    | Follow a user                  | Required       |
| `DELETE` | `/users/:id/unfollow`  | Unfollow a user                | Required       |
| `GET`    | `/users/:id/followers` | Get a user's followers         | Required       |
| `GET`    | `/users/:id/following` | Get users a user is following  | Required       |
| `POST`   | `/posts`               | Create a new post              | Required       |
| `GET`    | `/posts`               | Get all posts                  | Required       |
| `GET`    | `/posts/:id`           | Get a single post by ID        | Required       |
| `PUT`    | `/posts/:id`           | Update a post                  | Required       |
| `DELETE` | `/posts/:id`           | Delete a post                  | Required       |
| `POST`   | `/comments`            | Create a new comment           | Required       |
| `GET`    | `/comments/:post_id`   | Get all comments for a post    | Required       |
| `PUT`    | `/comments/:id`        | Update a comment               | Required       |
| `DELETE` | `/comments/:id`        | Delete a comment               | Required       |
| `GET`    | `/feed`                | Get the personalized user feed | Required       |
