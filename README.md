# AD SOYAD HANIM ATAŞ NO 24080410032
# GoLearn — Remote Learning Platform API

A production-ready Go backend for a remote learning platform with courses, lessons, quizzes, progress tracking and real-time WebSocket classrooms.

## Tech Stack

| Technology | Purpose |
|---|---|
| Go 1.22 + Gin | HTTP framework |
| GORM + SQLite | ORM & persistence |
| JWT (HS256) | Authentication |
| bcrypt | Password hashing |
| gorilla/websocket | Real-time classroom |
| swaggo/swag | Swagger / OpenAPI docs |
| Docker + Compose | Containerisation |

## Local Setup

### Prerequisites
- Go 1.22+
- GCC (for CGO / sqlite3)
- `swag` CLI

### Install swag
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Swagger docs & run
```bash
cd golearn
swag init
go run .
```

API will be available at **http://localhost:8090**

## Docker

```bash
docker compose up --build
```

API: http://localhost:8090  
Swagger UI: http://localhost:8090/swagger/index.html

## Swagger

Visit [http://localhost:8090/swagger/index.html](http://localhost:8090/swagger/index.html)

Click **Authorize** and enter `Bearer <your_jwt_token>`.

## Endpoints

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/api/auth/register` | Register (student/teacher) |
| POST | `/api/auth/login` | Login → JWT |

### Courses
| Method | Path | Auth | Role |
|---|---|---|---|
| GET | `/api/courses` | ✅ | any |
| GET | `/api/courses/:id` | ✅ | any |
| POST | `/api/courses` | ✅ | teacher |
| PUT | `/api/courses/:id` | ✅ | owner teacher |
| DELETE | `/api/courses/:id` | ✅ | owner teacher |

Query params for GET /api/courses: `page`, `limit`, `category`, `sort`

### Lessons
| Method | Path | Auth | Role |
|---|---|---|---|
| GET | `/api/courses/:id/lessons` | ✅ | any |
| POST | `/api/courses/:id/lessons` | ✅ | owner teacher |

### Quiz
| Method | Path | Auth | Role |
|---|---|---|---|
| GET | `/api/lessons/:id/quiz` | ✅ | any |
| POST | `/api/lessons/:id/quiz` | ✅ | owner teacher |
| POST | `/api/quiz/:id/submit` | ✅ | student |

### Progress
| Method | Path | Auth |
|---|---|---|
| POST | `/api/lessons/:id/complete` | ✅ |
| GET | `/api/my/progress` | ✅ |

### WebSocket
```
ws://localhost:8090/ws/classroom/:courseId?token=<JWT>
```

## Example Auth Flow

```bash
# Register a teacher
curl -X POST http://localhost:8090/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@test.com","password":"secret123","role":"teacher"}'

# Login
curl -X POST http://localhost:8090/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"secret123"}'

# Use the returned token
TOKEN="eyJhb..."

# Create a course
curl -X POST http://localhost:8090/api/courses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Go Basics","description":"Learn Go","category":"programming"}'
```

## WebSocket Example (wscat)

```bash
wscat -c "ws://localhost:8090/ws/classroom/1?token=$TOKEN"
# Send a message:
{"text":"Hello class!"}
```

## Git Commit Messages (suggested order)

```
feat: initialize Go module and project structure
feat(config): add environment-based configuration loader
feat(database): add GORM SQLite connection with AutoMigrate
feat(models): define User, Course, Lesson, Quiz, Question, Progress, QuizResult
feat(middleware): add JWT auth middleware
feat(middleware): add teacher-only RBAC guard
feat(middleware): add IP-based rate limiting (5 req/s, burst 10)
feat(handlers): implement register and login with bcrypt + JWT
feat(handlers): implement course CRUD with pagination and ownership checks
feat(handlers): implement lesson list and create
feat(handlers): implement quiz create, get and submit with scoring
feat(handlers): implement lesson complete and course progress tracking
feat(handlers): add WebSocket classroom with per-room goroutine management
feat(swagger): integrate swaggo/swag with annotated handlers
feat(docker): add multi-stage Dockerfile and docker-compose.yml
docs: add comprehensive README with curl examples and WebSocket usage
chore: add .gitignore
```
