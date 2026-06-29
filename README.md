# SpotSync

A smart parking reservation API built with Go. Drivers can browse parking zones and make reservations, while admins manage zones and view all bookings — all secured with JWT auth and concurrency-safe booking logic.

**Live URL:** _Deploy to Railway / Render and add your URL here_

---

## Features

- User registration and login with bcrypt password hashing
- JWT-based authentication with role claims (`driver`, `admin`)
- Role-based access control via middleware
- Parking zone CRUD (admin only) with zone types: `general`, `ev_charging`, `covered`
- Real-time available spots calculation on every zone response
- Reservation creation with row-level locking to prevent overbooking under concurrent load
- Cancel own reservations (drivers) or view all reservations (admin)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25+ |
| HTTP Framework | Echo v5 |
| ORM | GORM |
| Database | PostgreSQL (NeonDB / Supabase / local) |
| Authentication | JWT v5 + bcrypt |
| Validation | go-playground/validator v10 |
| Config | godotenv |

---

## Architecture

The project follows **Clean Architecture** with a domain-driven folder structure. Each domain (`user`, `parkingzone`, `reservation`) is fully self-contained.

```
cmd/
└── main.go                  Entry point — wires config, DB, and server

internal/
├── auth/
│   └── jwt.go               JWT generation & validation (with role support)
├── config/
│   ├── config.go            Loads environment variables
│   └── db.go                Connects to PostgreSQL via GORM, runs AutoMigrate
├── httpresponse/
│   └── response.go          Standardized JSON response helpers
├── middlewares/
│   └── auth.go              AuthMiddleware (JWT) and RoleMiddleware
├── server/
│   └── http.go              Echo server setup, validator wiring, route registration
└── domain/
    ├── user/                Auth module
    ├── parkingzone/         Parking zones module
    └── reservation/         Reservations module
```

### How the layers interact

```
HTTP Request
     │
     ▼
  Handler          Binds & validates input, returns JSON response
     │
     ▼
  Service          Business logic, error mapping, response building
     │
     ▼
 Repository        GORM database queries
     │
     ▼
  Entity           GORM model (maps to DB table)
```

Dependency injection is done manually in each domain's `register.go`, which wires the repository → service → handler chain and registers routes on the Echo instance.

### Concurrency Safety

Reservations use a **GORM transaction with `FOR UPDATE` row-level locking** on the parking zone row. This ensures that when two users simultaneously try to book the last spot, only one succeeds — the other gets a `409 Conflict`.

```
BeginTransaction
  → Lock parking_zone row (FOR UPDATE)
  → Count active reservations
  → Check capacity
  → Insert reservation (or return ErrZoneFull)
CommitTransaction
```

---

## Setup (Local)

### Prerequisites

- [Go 1.21+](https://go.dev/dl)
- PostgreSQL (local install, [NeonDB](https://neon.tech), or [Supabase](https://supabase.com))

### Steps

**1. Clone the repository**
```bash
git clone <your-repo-url>
cd spotsync
```

**2. Create your `.env` file**
```bash
cp .env.example .env
```

**3. Fill in your environment variables** (see table below)

**4. Download dependencies**
```bash
go mod tidy
```

**5. Run the server**
```bash
go run cmd/main.go
```

The server starts on `http://localhost:8080`. Tables are auto-migrated on startup.

**Optional — hot reload with Air**
```bash
air
```

---

## Environment Variables

| Variable | Required | Description | Example |
|---|---|---|---|
| `PORT` | Yes | Port the server listens on | `8080` |
| `DSN` | Yes | PostgreSQL connection string | `postgresql://user:pass@host/db?sslmode=require` |
| `JWT_SECRET` | Yes | Secret key used to sign JWT tokens | `change-me-in-production` |

---

## API Endpoints

Base URL: `http://localhost:8080/api/v1`

All authenticated and admin endpoints require the `Authorization` header:
```
Authorization: Bearer <token>
```

### Auth

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/auth/register` | Public | Register a new user |
| `POST` | `/auth/login` | Public | Login and receive a JWT token |

**Register request body:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "secret123",
  "role": "driver"
}
```
> `role` is optional — defaults to `driver`. Accepted values: `driver`, `admin`.

**Login request body:**
```json
{
  "email": "jane@example.com",
  "password": "secret123"
}
```

---

### Parking Zones

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `GET` | `/zones` | Public | List all zones with available spots |
| `GET` | `/zones/:id` | Public | Get a single zone by ID |
| `POST` | `/zones` | Admin | Create a new parking zone |
| `PUT` | `/zones/:id` | Admin | Update an existing zone |
| `DELETE` | `/zones/:id` | Admin | Delete a zone |

**Create / Update zone request body:**
```json
{
  "name": "Zone A",
  "type": "general",
  "total_capacity": 50,
  "price_per_hour": 5.00
}
```
> `type` accepted values: `general`, `ev_charging`, `covered`

---

### Reservations

| Method | Endpoint | Access | Description |
|---|---|---|---|
| `POST` | `/reservations` | Authenticated | Reserve a parking spot |
| `GET` | `/reservations/my-reservations` | Authenticated | View your own reservations |
| `DELETE` | `/reservations/:id` | Authenticated | Cancel your reservation |
| `GET` | `/reservations` | Admin | View all reservations |

**Create reservation request body:**
```json
{
  "zone_id": 1,
  "license_plate": "DHK-1234"
}
```

---

## License

MIT
