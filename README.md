# DBackend - Investment Platform API

## Overview

DBackend is a robust Go-based API service that powers the investment platform, facilitating connections between investors and startups. The system provides comprehensive deal flow management, matchmaking services, and secure authentication.

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Fiber (HTTP server)
- **Database**: MongoDB
- **Authentication**: JWT-based auth
- **ML Services**: Python FastAPI for matchmaking algorithms
- **Containerization**: Docker & Docker Compose
- **CI/CD**: GitHub Actions

## Features

- **User Management**: Complete user lifecycle with role-based access control
- **Authentication**: Secure JWT-based authentication system
- **Deal Flow Management**: Track and manage investment opportunities
- **Matchmaking**: AI-powered matching between investors and startups
- **Meeting Coordination**: Schedule and manage meetings
- **Task Management**: Track and assign tasks
- **Grant Management**: Handle grant applications and approvals

## Getting Started

### Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- MongoDB (or use the provided Docker setup)
- Python 3.8+ (for MatchMaking service)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-organization/dbackend.git
   cd dbackend
   ```

2. Set up environment variables by creating a `.env` file:
   ```
   APP_ENV=development
   PORT=8080
   BLUEPRINT_DB_HOST=localhost
   BLUEPRINT_DB_PORT=27017
   BLUEPRINT_DB_USERNAME=dbadmin
   BLUEPRINT_DB_ROOT_PASSWORD=securepassword
   BLUEPRINT_DB_DATABASE=ddb
   ```

3. Start the MongoDB container:
   ```bash
   make docker-run
   ```

4. Build and run the application:
   ```bash
   make build
   make run
   ```

5. For development with live reload:
   ```bash
   make watch
   ```

## Project Structure

```
DBackend/
├── cmd/                    # Application entry points
│   └── api/                # Main API server
├── internal/               # Private application code
│   ├── database/           # Database interfaces and implementations
│   ├── models/             # Data models
│   ├── server/             # HTTP server setup
│   │   ├── middleware/     # HTTP middleware
│   │   └── routes/         # API route definitions
│   └── utils/              # Utility functions
├── MatchMakingService/     # Python-based ML service for matching
├── scripts/                # Utility scripts
├── Dockerfile              # Docker configuration
├── docker-compose.yml      # Docker Compose configuration
├── go.mod                  # Go module definition
└── Makefile                # Build and development commands
```

## API Endpoints

The API follows RESTful conventions and is versioned. All endpoints are prefixed with `/api/v1/`.

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Authenticate and receive JWT
- `GET /api/v1/auth/refresh` - Refresh JWT token

### Users
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

### Deal Flow
- `POST /api/v1/dealflow` - Add startup to dealflow
- `GET /api/v1/dealflow` - List all dealflow entries
- `GET /api/v1/dealflow/:id` - Get dealflow by ID
- `PUT /api/v1/dealflow/:id` - Update dealflow entry
- `DELETE /api/v1/dealflow/:id` - Remove from dealflow

For a complete list of endpoints with example requests, see `curl_examples.md`.

## Development

### Makefile Commands

```bash
make all           # Build and run tests
make build         # Build the application
make run           # Run the application
make docker-run    # Create and start DB container
make docker-down   # Shutdown DB container
make test          # Run unit tests
make itest         # Run integration tests
make watch         # Live reload during development
make clean         # Clean up binary from last build
```

### Testing

The project includes both unit tests and integration tests:

```bash
# Run all tests
make test

# Run integration tests (requires DB)
make itest
```

### Adding New Features

1. Create appropriate models in `internal/models`
2. Implement database operations in `internal/database`
3. Add routes in `internal/server/routes`
4. Update the route registration in `internal/server/routes.go`
5. Write tests for your implementation

## Deployment

### Production Setup

1. Build the Docker image:
   ```bash
   docker build -t dbackend:latest .
   ```

2. Deploy using Docker Compose:
   ```bash
   docker-compose up -d
   ```

### Server Configuration

The application is designed to run behind Nginx. A sample configuration is provided in `nginx-config.conf`.

## MatchMaking Service

The MatchMaking service is a Python-based ML service that provides intelligent matching between investors and startups.

### Setup

1. Navigate to the MatchMaking directory:
   ```bash
   cd MatchMakingService
   ```

2. Create a virtual environment:
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   ```

3. Install dependencies:
   ```bash
   pip install fastapi uvicorn tensorflow pandas scikit-learn joblib
   ```

4. Start the service:
   ```bash
   uvicorn app:app --host 0.0.0.0 --port 4040
   ```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -am 'Add new feature'`
4. Push to the branch: `git push origin feature/my-feature`
5. Submit a pull request

## Troubleshooting

### Common Issues

- **Database Connection Issues**: Ensure MongoDB is running and credentials are correct
- **Build Errors**: Make sure you're using Go 1.22 or higher
- **API Errors**: Check the logs for detailed error messages

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For questions or support, please contact the development team at briankimathi94@gmail.com.
