```
# Go - Backend Development Task
## User API with DOB and Calculated Age

A RESTful API built in Go to manage users with their `name` and `dob` (date of birth). The API calculates and returns a user's `age` dynamically when fetching user details.

---

## Tech Stack

| Technology | Version | Purpose |
|---|---|---|
| Go | 1.22 | Programming language |
| GoFiber | v2 | Web framework |
| MySQL | 8.0 | Database |
| SQLC | latest | Type-safe SQL query generation |
| Uber Zap | latest | Structured logging |
| go-playground/validator | v10 | Input validation |
| golang-migrate | v4.17.0 | Database migrations |
| Docker | latest | Containerization |

---

## Prerequisites

Make sure you have the following installed:

- [Go 1.22+](https://go.dev/dl/)
- [MySQL 8.0](https://dev.mysql.com/downloads/)
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) *(optional, for Docker setup)*

---

## Setup Without Docker

### 1. Clone the repository
```bash
git clone https://github.com/aadityahammad-2002/Go---Backend-Development-Task.git
cd Go---Backend-Development-Task
```

### 2. Configure environment
```bash
cp .env.example .env
```
Open `.env` and fill in your values:
```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=usersdb
APP_PORT=3000
```

### 3. Create the database
```bash
mysql -u root -p
```
```sql
CREATE DATABASE usersdb;
EXIT;
```

### 4. Run migrations
```bash
migrate -path db/migrations -database "mysql://root:your_password@tcp(localhost:3306)/usersdb" up
```

### 5. Install dependencies
```bash
go mod tidy
```

### 6. Generate SQLC code
```bash
sqlc generate -f db/sqlc/sqlc.yaml
```

### 7. Start the server
```bash
go run ./cmd/server/main.go
```

Server runs at: `http://localhost:3000`

---

## Setup With Docker

### 1. Clone the repository
```bash
git clone https://github.com/aadityahammad-2002/Go---Backend-Development-Task.git
cd Go---Backend-Development-Task
```

### 2. Configure environment
```bash
cp .env.example .env
```
Set `DB_HOST=db` (not localhost) in `.env`:
```
DB_HOST=db
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=usersdb
APP_PORT=3000
```

### 3. Start everything
```bash
docker-compose up --build
```

### 4. Run migrations (first time only)
Open a new terminal:
```bash
docker-compose exec app migrate -path db/migrations -database "mysql://root:root@tcp(db:3306)/usersdb" up
```

Server runs at: `http://localhost:3000`

---

## API Endpoints

### Create User
**POST** `/users`
```json
Request:
{
  "name": "Alice",
  "dob": "1990-05-10"
}

Response 201:
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10"
}
```

### Get User by ID
**GET** `/users/:id`
```json
Response 200:
{
  "id": 1,
  "name": "Alice",
  "dob": "1990-05-10",
  "age": 35
}
```

### Update User
**PUT** `/users/:id`
```json
Request:
{
  "name": "Alice Updated",
  "dob": "1991-03-15"
}

Response 200:
{
  "id": 1,
  "name": "Alice Updated",
  "dob": "1991-03-15"
}
```

### Delete User
**DELETE** `/users/:id`
```
Response: 204 No Content
```

### List All Users
**GET** `/users`
```
Optional: ?page=1&limit=10
```
```json
Response 200:
[
  {
    "id": 1,
    "name": "Alice",
    "dob": "1990-05-10",
    "age": 35
  }
]
```

---

## Running Tests

```bash
go test ./...
```

For verbose output:
```bash
go test -v ./internal/service/...
```

Tests cover:
- Birthday is today
- Birthday already passed this year
- Birthday not yet this year
- Leap year birthday (Feb 29)

---

## Environment Variables

| Variable | Description | Example |
|---|---|---|
| DB_HOST | MySQL host | localhost |
| DB_PORT | MySQL port | 3306 |
| DB_USER | MySQL username | root |
| DB_PASSWORD | MySQL password | root |
| DB_NAME | Database name | usersdb |
| APP_PORT | App listen port | 3000 |

---

## Project Structure

```
├── cmd/server/main.go              → App entry point, wires dependencies
├── config/config.go                → Loads environment variables
├── db/
│   ├── migrations/                 → SQL migration files
│   └── sqlc/                       → SQLC config and generated DB code
├── internal/
│   ├── handler/user_handler.go     → HTTP request handlers
│   ├── repository/user_repository.go → Database access layer
│   ├── service/user_service.go     → Business logic + age calculation
│   ├── service/user_service_test.go → Unit tests
│   ├── routes/routes.go            → Route registration
│   ├── middleware/request_id.go    → Injects X-Request-ID header
│   ├── middleware/logger.go        → Logs request duration
│   ├── models/user.go              → Request/response structs
│   └── logger/logger.go            → Uber Zap logger setup
├── .env.example                    → Environment variable template
├── Dockerfile                      → Multi-stage Docker build
└── docker-compose.yml              → App + MySQL services
```

---

## Age Calculation

Age is never stored in the database. It is always calculated at runtime using Go's `time` package in `internal/service/user_service.go`. The logic compares today's date to the user's `dob` and subtracts 1 if the birthday hasn't occurred yet this year.
```
