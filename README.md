# Sorcerer's Labyrinth Backend

REST API для интерактивной игры-книги с системой боёв, бонусов и геймплейных механик.

## Stack

- **Language:** Go 1.23.0
- **Framework:** Gin
- **Database:** PostgreSQL 13+ (Docker)
- **Architecture:** Clean Architecture + DI (samber/do)
- **Authentication:** JWT tokens

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23.0+
- Make

### Installation

```bash
# Clone repository
git clone <repo-url>
cd labirint-kolduna.ru-go

# Start PostgreSQL
make up

# Run migrations
make migrate

# (Optional) Seed database
make migrate-seed

# Start application
make run
```

Application runs on `http://localhost:8888`

### Docker Services

| Service | Port | Description |
|---------|------|-------------|
| app     | 8888 | Go Gin application |
| nginx   | 80   | Reverse proxy |
| postgres| 5432 | PostgreSQL database |

## Development

### Commands

```bash
# Run application
make run              # go run cmd/main.go

# Build
make build            # go build -o main cmd/main.go
make run-build        # build + ./main

# Testing
make test             # go test -v ./tests
make test-all         # go test -v ./modules/.../tests/...
make test-coverage    # coverage report

# Migrations
make migrate              # run migrations
make migrate-rollback     # rollback migrations
make migrate-status       # show migration status
make migrate-seed         # run seeders

# Docker
make up               # docker-compose up -d
make down             # docker-compose down
make logs             # docker-compose logs -f
```

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── modules/
│   ├── auth/                   # Authentication
│   │   ├── controller/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── dto/
│   │   └── routes.go
│   ├── user/                   # User management
│   │   ├── controller/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── dto/
│   │   └── routes.go
│   └── game/                   # Game mechanics
│       ├── controller/
│       ├── service/
│       ├── repository/
│       ├── dto/
│       ├── battle/
│       ├── bonus/
│       ├── bribe/
│       ├── dice/
│       ├── sleep/
│       └── routes.go
├── providers/                  # DI container (samber/do)
├── middlewares/                # CORS, auth, logging
├── database/
│   ├── migrations/             # GORM migrations
│   └── seeders/                # Seed data
├── pkg/
│   ├── constants/
│   └── helpers/
└── .claude/                    # AI/Documentation tools
```

## API Documentation

### Endpoints

#### Auth Module (`/api/auth`)
- `POST /register` - Register new user
- `POST /login` - Authenticate and receive tokens
- `POST /refresh` - Refresh access token
- `POST /logout` - Invalidate tokens

#### User Module (`/api/user`)
- `GET /` - Get all users
- `GET /me` - Get current user profile
- `PUT /:id` - Update user
- `DELETE /:id` - Delete user

#### Game Module (`/api/game`)
- `GET /get-section` - Get current game section
- `GET /role-the-dice` - Roll dice
- `POST /choice` - Make game choice
- `POST /move` - Move to new section
- `POST /battle` - Execute battle mechanics
- `GET /profile` - Get player profile
- `POST /ability/meds` - Use meds ability
- `POST /ability/bonus` - Use bonus ability
- `POST /ability/sleep` - Enter Sleepy Kingdom
- `POST /ability/sleep/choice` - Make Sleepy Kingdom choice
- `POST /ability/bribe` - Use bribe ability
- `GET /map` - Get game map

**Total: 20 endpoints**

### Detailed API Specs

Full API documentation available in:
- **Auth:** [`.claude/output/contracts/auth.md`](.claude/output/contracts/auth.md)
- **User:** [`.claude/output/contracts/user.md`](.claude/output/contracts/user.md)
- **Game:** [`.claude/output/contracts/game.md`](.claude/output/contracts/game.md)

### QA Documentation

- **Test Checklist:** [`.claude/output/qa/checklist.md`](.claude/output/qa/checklist.md) - 100+ test scenarios
- **Postman Collection:** [`.claude/output/qa/postman-collection.json`](.claude/output/qa/postman-collection.json) - Automated testing suite

## Authentication

All game and user endpoints (except auth endpoints) require JWT token:

```bash
# Login to get tokens
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"name":"wizard","password":"secret123"}'

# Use access_token in Authorization header
curl -X GET http://localhost:8888/api/game/profile \
  -H "Authorization: Bearer {access_token}"
```

## Architecture

### Clean Architecture Pattern

```
Controller → Service → Repository → Database
```

- **Controller:** HTTP request handling, validation
- **Service:** Business logic
- **Repository:** Database access
- **Database:** PostgreSQL with GORM ORM

### Dependency Injection

All dependencies managed via `samber/do`:

```go
// providers/
func ProvideService(injector *do.Injector) *SomeService {
    return &SomeService{
        repo: do.MustInvoke[*Repository](injector),
    }
}
```

## Game Mechanics

### Core Features

- **Sections:** Interactive story nodes with choices
- **Battles:** Combat system with weapons and dice
- **Abilities:** Special powers (meds, bonus, sleep, bribe)
- **Sleepy Kingdom:** Special 12-section dream realm
- **Dice Rolling:** Randomized outcomes (d6)
- **Map:** Track visited sections and available paths

### Player Stats

- Health & Max Health
- Gold currency
- Meds (healing items)
- Weapons & inventory
- Buffs & debuffs
- Special bonuses

## Code Style

- Package names: lowercase, single word
- Exported symbols: PascalCase
- Private symbols: camelCase
- No comments/docstrings unless required
- `gofmt` formatting
- Context first parameter for server operations

## Contributing

1. Create feature branch from `develop`
2. Follow Clean Architecture pattern
3. Add tests for new features
4. Update API documentation
5. Submit PR to `develop`

## Troubleshooting

### Database Connection Issues
```bash
# Check PostgreSQL status
docker-compose logs postgres

# Restart services
make down && make up

# Access PostgreSQL directly
docker exec -it go-gin-clean-starter-db /bin/sh
```

### Migration Errors
```bash
# Check migration status
make migrate-status

# Rollback and retry
make migrate-rollback
make migrate
```

## License

[License information]

## Contact

[Contact information]
