# Knowledge Capsule API

Knowledge Capsule API is a Go-based backend service designed to manage "knowledge capsules"—bite-sized pieces of information categorized by topics and tags. It provides a RESTful API for creating, retrieving, and searching these capsules, along with user authentication and topic management.

## Features

- **User Authentication**: Secure registration and login using JWT (JSON Web Tokens).
- **Capsule Management**: Create and retrieve knowledge capsules with support for private/public visibility.
- **Topic Organization**: Categorize capsules into topics.
- **Search**: Search functionality to find specific capsules.
- **Tagging**: Add tags to capsules for better organization.
- **File-based Storage**: Simple JSON file-based persistence for users, topics, and capsules (easy to set up, no database required).

## Tech Stack

- **Language**: [Go](https://go.dev/) (1.23+)
- **Containerization**: [Docker](https://www.docker.com/)
- **Build Tool**: [Make](https://www.gnu.org/software/make/)
- **Live Reload**: [Air](https://github.com/air-verse/air)
- **Git Hooks**: [Lefthook](https://github.com/evilmartians/lefthook)

## Prerequisites

Ensure you have the following installed on your system:

- [Go](https://go.dev/dl/) (1.23 or later)
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Make](https://www.gnu.org/software/make/)

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd knowledge-capsule-api
```

### 2. Environment Setup

Create a `.env` file in the root directory. You can copy the example below:

```bash
# .env
PORT=8080
GO_ENV=development
JWT_SECRET=your_super_secret_key_here
```

> **Note**: You can generate a secure JWT secret using `make g-jwt`.

### 3. Run with Docker (Recommended)

To start the application in **development mode** (with live reload):

```bash
make up-dev
```

The API will be available at `http://localhost:8081`.

To start in **production mode**:

```bash
make up
```

The API will be available at `http://localhost:8080`.

To stop the containers:

```bash
make down-dev  # for dev
# or
make down      # for prod
```

### 4. Run Locally

If you prefer to run without Docker:

1. **Install dependencies**:

    ```bash
    make install
    ```

2. **Run the server**:

    ```bash
    make run
    ```

    This uses `air` for live reloading.

    Or build and run the binary directly:

    ```bash
    make build-local
    ./tmp/server
    ```

## API Documentation

Interactive Swagger documentation is available at `/docs/index.html` when the server is running.

### Authentication

- **POST** `/api/auth/register`
  - Register a new user.
  - Body: `{ "name": "John Doe", "email": "john@example.com", "password": "securepassword" }`
- **POST** `/api/auth/login`
  - Login and receive a JWT token.
  - Body: `{ "email": "john@example.com", "password": "securepassword" }`

### Topics

*Requires Authentication Header: `Authorization: Bearer <token>`*

- **GET** `/api/topics`
  - Get all topics.
- **POST** `/api/topics`
  - Create a new topic.
  - Body: `{ "name": "Golang", "description": "All things Go" }`

### Capsules

*Requires Authentication Header: `Authorization: Bearer <token>`*

- **GET** `/api/capsules`
  - Get all capsules for the logged-in user.
- **POST** `/api/capsules`
  - Create a new capsule.
  - Body:

        ```json
        {
          "title": "Interfaces in Go",
          "content": "Interfaces are named collections of method signatures...",
          "topic": "Golang",
          "tags": ["programming", "go"],
          "is_private": false
        }
        ```

### Search

*Requires Authentication Header: `Authorization: Bearer <token>`*

- **GET** `/api/search?q=<query>`
  - Search capsules by title or content.

### Health Check

- **GET** `/health`
  - Check if the service is running.

## Project Structure

```
knowledge-capsule-api/
├── config/         # Configuration loading
├── handlers/       # HTTP request handlers
├── middleware/     # HTTP middleware (Auth, Logger, etc.)
├── models/         # Data structures
├── store/          # Data persistence logic (JSON file store)
├── utils/          # Utility functions
├── data/           # JSON data storage (users.json, etc.)
├── scripts/        # Helper scripts
├── Dockerfile      # Production Dockerfile
├── Dockerfile.dev  # Development Dockerfile
├── compose.yaml    # Docker Compose configuration
├── Makefile        # Build and run commands
└── main.go         # Application entry point
```

## Development Commands

The `Makefile` provides several useful commands:

- `make help`: Show all available commands.
- `make run`: Run the app locally with live reload.
- `make build-local`: Build the binary locally.
- `make fmt`: Format code.
- `make vet`: Run `go vet`.
- `make tidy`: Run `go mod tidy`.
- `make test`: Run tests (if available).
- `make g-jwt`: Generate a random JWT secret.
