# SpotSync API

Smart Parking & EV Charging Reservation System built with Go, Echo, GORM, and PostgreSQL.

## Features

- User authentication (register, login) with JWT and bcrypt
- Role-based access control (driver, admin)
- Parking zone management (CRUD for admins)
- Public parking zone browsing with real-time availability
- Reservation system with concurrency-safe booking using GORM transactions and row-level locking (`FOR UPDATE`)
- Clean Architecture with layered separation (DTO, Handler, Service, Repository, Models)

## Tech Stack

- Go 1.25+
- Echo v5
- GORM with PostgreSQL
- JWT v5
- bcrypt
- go-playground/validator

## Architecture

```
cmd/main.go                 -> Entry point
internal/
  auth/                     -> JWT service
  config/                   -> Environment & DB config
  httpresponse/             -> Standardized JSON responses
  middlewares/              -> Auth & Role middleware
  server/                   -> HTTP server setup
  domain/
    user/                   -> Auth module
    parkingzone/            -> Parking zones module
    reservation/            -> Reservations module
```

Each domain follows Clean Architecture:
- **DTO** — Request/Response structures
- **Handler** — HTTP layer, binds & validates input, returns JSON
- **Service** — Business logic
- **Repository** — GORM database operations
- **Entity** — GORM models

## API Endpoints

### Authentication
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/api/v1/auth/register` | Public | Register a new user |
| POST | `/api/v1/auth/login` | Public | Login and get JWT token |

### Parking Zones
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| GET | `/api/v1/zones` | Public | List all zones with availability |
| GET | `/api/v1/zones/:id` | Public | Get single zone details |
| POST | `/api/v1/zones` | Admin | Create a new zone |
| PUT | `/api/v1/zones/:id` | Admin | Update a zone |
| DELETE | `/api/v1/zones/:id` | Admin | Delete a zone |

### Reservations
| Method | Endpoint | Access | Description |
|--------|----------|--------|-------------|
| POST | `/api/v1/reservations` | Authenticated | Reserve a parking spot |
| GET | `/api/v1/reservations/my-reservations` | Authenticated | View my reservations |
| DELETE | `/api/v1/reservations/:id` | Authenticated | Cancel a reservation |
| GET | `/api/v1/reservations` | Admin | View all reservations |

## Setup (Local)

1. **Clone the repository**

2. **Set up PostgreSQL** (local or cloud like NeonDB/Supabase)

3. **Create `.env` file**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your database credentials and JWT secret.

4. **Install dependencies**
   ```bash
   go mod tidy
   ```

5. **Run the server**
   ```bash
   go run cmd/main.go
   ```

   Or use Air for hot-reloading:
   ```bash
   air
   ```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | Server port (default: 8080) |
| `DSN` | PostgreSQL connection string |
| `JWT_SECRET` | Secret key for JWT signing |

## Concurrency Safety

The reservation system uses **GORM Transactions** combined with **Row-Level Locking (`FOR UPDATE`)** on the parking zone record to prevent race conditions when multiple users simultaneously try to book the last available spot.

```go
db.Transaction(func(tx *gorm.DB) error {
    // Lock the zone row
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)
    // Count active reservations
    // Check capacity
    // Create reservation atomically
})
```

## License

MIT
