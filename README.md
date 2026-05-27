# Go Social API

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)
![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat&logo=docker)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-336791?style=flat&logo=postgresql)
![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?style=flat&logo=swagger)

A layered REST API for a social media platform.

## ⚡ Features

- **Auth:** Secure JWT authentication (Signup, Login, Password Reset).
- **Content:** CRUD for Posts, Comments, and Nested Replies.
- **Social:** Like system, Follow/Unfollow users, and Tagging.
- **Feed:** Personalized feeds.
- **Search:** User search and Tag filtering.
- **Architecture:** Layered architecture with RabbitMQ and Websocket for real-time notifications.

## Motivation

    A learning project for backend with go.I didn't find relevant example projects on github so i want to write my own.That way people can inspect my project and can find useful.Im planning using this project for my future project as a reference.

## 🚀 Quick Start

### 1. Environment Setup

Create a `.env` file in the root directory:

```env
POSTGRES_USER=deneme
POSTGRES_PASSWORD=deneme
POSTGRES_DB=go_social
POSTGRES_PORT=5999
POSTGRES_HOST=localhost
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
JWT_SECRET=supersecretkey
TEST_DB_URL=postgres://deneme:deneme@localhost:5432/go_social_test?sslmode=disable
RABBITMQ_URL=amqp://guest:guest@localhost:5672
```

### 2. Run with Docker

The easiest way to start the API and Database:

```bash
make up
```

## 3. To run locally

```bash
make local
```

_API will be available at `http://localhost:3000`_

## Tech Stack

- **Go**
- **Postgresql**
- **RabbitMQ**
- **Docker**
- **golang-migrate**
- **swaggo**
- **zap**

## Documentation

Interactive API documentation is available via Swagger UI:

**[http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)**

## Project Structure

The project follows a **Clean Architecture** pattern to ensure separation of concerns and scalability:

- **`cmd/`**: Application entry points (API server, Seeder).
- **`internal/controller`**: **Gin** handlers responsible for request validation and response formatting.
- **`internal/services`**: Business layer, using by controllers via dependency injection.
- **`internal/database`**: Data access layer using raw SQL for **PostgreSQL** interactions.
- **`internal/mock`**: Mock database structs for testing purpose.
- **`internal/model`**: Core domain entities and database schemas.
- **`internal/dto`**: Data Transfer Objects for strict API contract definition.
- **`internal/middleware`**: **JWT Auth**, Logging, and Rate Limiting.
- **`internal/routes`**: All routes registered here.
- **`internal/util`**: Config struct and helper functions.
- **`internal/appError`**: Define custom error type.
- **`internal/migration`**: Database migrations for golang-migrate.
- **`internal/broker`**: Contains functions for rabbitmq to connect,declare channel etc.
- **`internal/ws`**: Websocket hub and client handling.
- **`docs/`**: Swagger documentation.
- **`internal/worker`**: RabbitMQ workers.Consuming event here.
