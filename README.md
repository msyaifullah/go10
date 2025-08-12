# Go10 Golang Template service

A Go-based REST API service built with Gin framework, following clean architecture principles.

## Features

- 🚀 **Live Reload Development** - Using Air for automatic rebuilds during development
- 🏗️ **Clean Architecture** - Separation of concerns with layers (handlers, services, repositories)
- 🗄️ **Database Integration** - PostgreSQL with migrations support
- 🔧 **Configuration Management** - Environment-based configuration
- 📝 **Structured Logging** - Comprehensive logging with masking
- 🐳 **Docker Support** - Containerized deployment
- 🧪 **Testing** - Comprehensive test coverage

## Design System

[REQUIREMENT](/docs/Requirement.md)

## Prerequisites

- **Golang (go1.23.6):** Ensure you have the correct versions installed go env. You can setup with brew.
   - brew install go
   - go version
   - brew update && brew upgrade go -> to update your golang version
- **A code editor:** Choose a code editor like Visual Studio Code

## Quick Start

1. **Clone the Repository (private):**

    ```bash
    git clone https://github.com/msyaifullah/go10.git
    ```
2. **Navigate to the Project Directory:**

    ```bash
    cd go10
    ```

3. **Install Dependencies:**

    ```bash
    make clean
    make all
    ```

4. **Install Dependencies:**

    ```bash
    # Setup complete development environment (Air + dependencies)
    make setup-dev
    ```

6. **Start the app:**
    
    ```bash
    # Start with live reload (recommended for development)
    make dev

    # Or start without live reload
    make run
    ```

    The server will be available at `http://localhost:8080`



**Additional Notes:**

- **Testing:**

  - To run tests, use the following command:

    ```bash
    make test
    ```
  - To run tests coverage, use the following command:

    ```bash
    make coverage
    ```
    ![coverage](/docs/coverage-1.png 'coverage 1')
    ![coverage](/docs/coverage-2.png 'coverage 2')


## Project Structure

```
loan-svc/
├── cmd/                    # Application entry points
│   ├── cli/                # CLI application
│   └── server/             # HTTP server
├── configs/                # Configuration files
├── internal/               # Private application code
│   ├── application/        # Application layer
│   ├── constant/        # Constants
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Data models
│   ├── repositories/       # Data access layer
│   ├── routes/             # Route definitions
│   └── services/           # Business logic
├── pkg/                    # Public application code
│   ├── adapters/           # External service adapters (use this all for mock 3rd party)
│   ├── config/             # Configuration management
│   ├── database/           # Database connection
│   ├── logger/             # Logging utilities
│   └── response/           # Unify response 
├── migrations/             # Database migrations
├── .air.toml               # Air configuration
├── Dockerfile              # 
├── docker-compose.yml      # 
├── go.mod                  # 
├── go.sum                  # 
├── Makefile                # Build automation
└── docker-compose.yml      # Docker services
```

## Others tools 
  - To run with script, use the following command in project root:
    ```bash
        Available targets:

        General:
            make all                    - Clean, download dependencies, run tests, and build
            make clean                  - Clean build directory
            make deps                   - Download dependencies
            make setup-dev              - Setup development environment (Air + deps)
            make test                   - Run tests
            make bench                  - Run benchmarks
            make coverage               - Run tests coverage
            make coverage-check         - Run coverage check threshold

        Server:
            make build-server           - Build server application
            make run                    - Run server application
            make dev                    - Run server with Air (live reload)
            make install-air            - Install Air for development

        CLI:
            make build-cli              - Build CLI application
            make migrate                - Run database migrations
            make migrate-debug          - Run database migrations with debug
            make migrate-down           - Rollback database migrations
            make migrate-down-debug     - Rollback database migrations with debug

        Build:
            make build                  - Build both server and CLI
    ```  

## Troubleshooting

If you encounter any issues during the setup process, consider the following:

- **Golang version and Postgres version:** Ensure you're using compatible versions.
- **Dependency:** Try resolving conflicts by updating or downgrading dependencies.
- **Check Project-Specific Documentation:** Refer to the project's README or other documentation for any specific setup instructions or requirements.

By following these steps, you should be able to set up the development environment and start working on the project.