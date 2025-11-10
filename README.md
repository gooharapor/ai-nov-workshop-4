# User Management API

A clean and scalable REST API built with Go, Fiber framework, and GORM ORM for managing users.

## 🏗️ Project Structure

```
.
├── main.go                 # Application entry point
├── models/                 # Data models
│   └── user.go            # User model with GORM tags
├── database/              # Database configuration
│   └── database.go        # GORM connection and migration
├── handlers/              # HTTP request handlers
│   └── user_handler.go    # User CRUD handlers
├── routes/                # Route definitions
│   └── routes.go          # API routes setup
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependency checksums
└── users.db              # SQLite database file
```

## 🚀 Tech Stack

- **Go** 1.25.4
- **Fiber** v2.52.9 - Express-inspired web framework
- **GORM** v1.31.1 - ORM library
- **SQLite** - Database

## 📋 Features

- ✅ Clean Architecture with separated layers
- ✅ RESTful API design
- ✅ GORM ORM for database operations
- ✅ Auto-migration of database schema
- ✅ CORS enabled
- ✅ Request logging middleware
- ✅ Soft delete support (via GORM)
- ✅ JSON responses
- ✅ Input validation

## 🔧 Installation

1. Install dependencies:

```bash
go mod download
```

2. Run the application:

```bash
go run main.go
```

Server will start on `http://localhost:3000`

## 📡 API Endpoints

### Root

- `GET /` - Hello world endpoint

### Users

- `GET /users` - Get all users
- `GET /users/:id` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

## 📝 API Examples

### Get all users

```bash
curl http://localhost:3000/users
```

### Get user by ID

```bash
curl http://localhost:3000/users/1
```

### Create user

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "081-234-5678",
    "address": "123 Main St, Bangkok",
    "avatar": "https://i.pravatar.cc/150?img=1"
  }'
```

### Update user

```bash
curl -X PUT http://localhost:3000/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "phone": "081-999-9999",
    "address": "456 New Address, Bangkok",
    "avatar": "https://i.pravatar.cc/150?img=2"
  }'
```

### Delete user

```bash
curl -X DELETE http://localhost:3000/users/1
```

## 🗃️ Database Schema

### Users Table

| Field      | Type     | Description                     |
| ---------- | -------- | ------------------------------- |
| id         | INTEGER  | Primary key (auto-increment)    |
| name       | TEXT     | User's full name (required)     |
| email      | TEXT     | User's email (unique, required) |
| phone      | TEXT     | Phone number                    |
| address    | TEXT     | Address                         |
| avatar     | TEXT     | Avatar/profile image URL        |
| created_at | DATETIME | Timestamp of creation           |
| updated_at | DATETIME | Timestamp of last update        |
| deleted_at | DATETIME | Soft delete timestamp           |

## 🎯 Architecture Benefits

1. **Separation of Concerns**: Each layer has a specific responsibility
2. **Maintainability**: Easy to locate and modify code
3. **Testability**: Handlers and models can be tested independently
4. **Scalability**: Easy to add new features and endpoints
5. **GORM Benefits**:
   - Automatic migrations
   - Soft deletes
   - Query builder
   - Associations support

## 🔄 Migration from Raw SQL

This project was refactored from raw SQL to GORM, providing:

- Type-safe database operations
- Automatic schema migrations
- Cleaner and more maintainable code
- Better error handling
- Soft delete capability

## 📦 Dependencies

```go
require (
    github.com/gofiber/fiber/v2 v2.52.9
    gorm.io/gorm v1.31.1
    gorm.io/driver/sqlite v1.6.0
)
```

## 🤝 Contributing

Feel free to submit issues and enhancement requests!

## 📄 License

MIT License
