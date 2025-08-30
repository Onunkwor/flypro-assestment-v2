# FlyPro Backend API

## 📌 Overview

This project is a backend API for Travel Expense Management, built for the FlyPro Backend Engineer Assessment. It demonstrates:

- Clean architecture with repository, service, and handler layers
- Request validation using DTOs and custom validators
- Currency conversion integration with USD normalization
- Redis caching for performance
- Goose for database migrations
- testing with mocks

## 🏗 Project Structure

```
cmd/
└── server/main.go          # Entry point
internal/
├── config/                 # DB, Redis, env configs
├── dto/                    # Request/response DTOs
├── handlers/               # HTTP handlers
├── repository/             # Data access layer
├── services/               # Business logic
├── models/                 # GORM models
├── middleware/             # Logging, CORS, rate limiting
└── utils/                  # Helpers (error formatting, etc.)
migrations/                 # Goose migration files
tests/                      # Mock repositories
docker-compose.yml          # Postgres + Redis setup
Makefile                    # Migration & test commands
.env                        # Environment variables
```

## ⚙️ Setup & Run

### 1. Clone Repo & Install Dependencies

```bash
git clone https://github.com/yourusername/flypro-backend.git
cd flypro-backend
go mod tidy
```

### 2. Environment Variables

Create a `.env` file in the project root:

```
PORT=
DB_URL=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
REDIS_ADDR=
CURRENCY_API_KEY=
CURRENCY_API=https://v6.exchangerate-api.com/v6/CURRENCY_API_KEY
```

### 3. Start Dependencies

```bash
docker-compose up -d
```

### 4. Run Migrations

```bash
make migrate-up
make migrate-status
```

### 5. Run Server

```bash
go run ./cmd/server
```

## 🗄 Database Schema

- **User**: Manages user accounts
- **Expense**: Tracks expenses with original Amount + Currency and computed AmountUSD
- **ExpenseReport**: Groups multiple expenses, stores Total in USD

**Design Decision**: I chose to persist both the original amount + currency and a converted AmountUSD. This preserves data integrity while enabling USD-based reporting.

## 📡 API Endpoints

### Users

- `POST /api/users` – Create user
- `GET /api/users/:id` – Get user details

### Expenses

- `POST /api/expenses` – Create expense
- `GET /api/expenses` – List expenses (pagination, filters)
- `GET /api/expenses/:id` – Get expense details
- `PUT /api/expenses/:id` – Update expense
- `DELETE /api/expenses/:id` – Delete expense

### Reports

- `POST /api/reports` – Create report
- `POST /api/reports/:id/expenses` – Add expenses to report
- `GET /api/reports` – List reports (pagination)
- `PUT /api/reports/:id/submit` – Submit report

### Download Postman Collection

[📥 FlyPro Assessment Collection](./postman/flypro-assestment.postman_collection.json)

## 🔄 Currency Conversion & Caching

- Integrated with a third-party currency API
- All expenses normalized to USD in upon expense creation
- Cached exchange rates in Redis (6-hour TTL)

## 🧪 Testing

I implemented service-layer tests using GoMock to validate core business logic with mocked repositories and external services:

- **User Service**: Prevents duplicate email registration, Successfully creates a new user.
- **Expense Service**: Handles currency conversion failures gracefully, Creates an expense with USD conversion & exchange rate applied, Updates an expense while preserving correct conversion values
- **Report Service**: Validates user existence before creating a report, Successfully creates a report with a valid user, Adds expenses to reports and increments totals, Submits reports, ensuring correct status transitions, Retrieves a user’s report expenses with pagination

**Note**: Focused on reporting tests instead of CreateExpense, as the design stores both original values and computed AmountUSD, showcasing USD conversion at the reporting stage.

Run tests with:

```bash
make test
```

## ✅ Submission Deliverables

- Clean Git history & commits
- Dockerized local setup (Postgres + Redis)
- Goose migrations with indexes & constraints
- Postman collection for API testing
