# Movie Application

A modern web application for managing and discovering movies, built with Go and Preact.

## 🚀 Features

- RESTful API with Swagger documentation
- User authentication with sign-up/sign-in functionality
- Movie browsing and favorites system
- Permission management
- Docker containerization
- Development environment with hot reload via Vite
- Production-ready static file serving with embed

## 🛠 Tech Stack

### Backend

- Go 1.23.5
- Custom micro-framework for routing and middleware
- GORM for database operations
- Swagger for API documentation
- Fiber for HTTP server

### Frontend

- Preact for UI components
- Vite for build tooling and development server
- TypeScript for type safety
- Tailwind CSS for styling

### Infrastructure

- Docker & Docker Compose
- Make for build automation
- PostgreSQL for data storage

## 📦 Prerequisites

- Go 1.23.5
- Node.js 18 or higher
- Docker & Docker Compose
- Make

## 🏃‍♂️ Quick Start

1. Clone the repository:

```bash
git clone https://github.com/guemidiborhane/factorydigitale.tech.git movie-app
cd movie-app
```

2. Configure your environment:

```bash
make setup
```

3. Start the development environment:

```bash
make dev
```

This will:

- Start PostgreSQL in Docker
- Run database migrations
- Start the Go API server with hot reload
- Start the Vite development server
- Set up reverse proxy from Go to Vite

Visit http://localhost:3000 to access the application (all requests will be properly routed through the Go server).

## 🔄 Development Mode

In development mode, the application:

- Serves the frontend through Vite's development server
- Proxies all non-API requests to Vite for hot module replacement
- Routes API requests (`/api/*`) to the Go backend
- Enables hot reload for both frontend and backend changes

```go
// Development proxy configuration
if app.AppConfig.IsDev() {
    url := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
    app.Use(
        proxy.Balancer(proxy.Config{
            Servers: []string{url},
            Next: func(c *fiber.Ctx) bool {
                return len(c.Path()) >= 4 && c.Path()[:4] == "/api"
            },
        }),
    )
}
```

## 🏗 Production Mode

In production, the application:

- Embeds the built frontend assets using Go's `embed` package
- Serves static files directly from the Go binary
- Handles all routing through the Go server
- Enables compression and recovery middleware

```go
//go:embed all:build
var static embed.FS

var FS = filesystem.New(filesystem.Config{
    Root:         http.FS(static),
    Browse:       false,
    PathPrefix:   "/build",
    Index:        "index.html",
    NotFoundFile: "/build/index.html",
})
```

## 📁 Project Structure

```
├── internal/           # Internal packages
│   ├── config/        # Application configuration
│   └── router/        # Router setup and middleware
├── pkg/               # Frontend packages
│   ├── movies/        # Movie-related features
│   │   └── pages/
│   │       ├── index.tsx        # Movie listing
│   │       └── favourites/      # Favorites management
│   ├── permissions/   # Permission management
│   │   └── pages/
│   │       └── index.tsx
│   └── users/         # User authentication
│       └── pages/
│           ├── sign_in/
│           ├── sign_out/
│           └── sign_up/
├── docs/              # Swaggo generated documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── docker/         # Docker related files (Dockerfiles, shell history, ...)
├── docker-compose.yml
└── Makefile
```

## 🗺 Application Routes

Frontend routes are organized by feature:

### Movies

- `/movies` - Main movie listing and management
- `/movies/favourites` - User's favorite movies

### User Authentication

- `/users/sign_in` - User login
- `/users/sign_up` - New user registration
- `/users/sign_out` - User logout

### Administration

- `/permissions` - Permission management interface

## 🔧 Available Make Commands

```bash
make        # List possible commands
make dev    # Start development environment
make build  # Build production Docker image
```

## 🚢 Deployment

1. Build the production Docker image:

```bash
make build
```

## 📚 API Documentation

API documentation is available through Swagger UI at `/swagger` when the server is running. You can also find the OpenAPI specification at `/swagger/doc.json`.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
