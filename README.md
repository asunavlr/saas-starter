# SaaS Starter

A Go backend foundation for multi-tenant SaaS products, built with PostgreSQL and Redis.

## Implemented foundation
- Email/password authentication with bcrypt and signed JWT access tokens.
- Organizations as tenants and membership-based authorization.
- Role-based access control: owner, admin and member.
- Tenant-scoped project endpoints.
- PostgreSQL schema for users, organizations, memberships, projects, plans, subscriptions and audit logs.
- Audit events for organization and project creation.
- Subscription and plan primitives (demo only; no payment provider is connected).
- Redis-backed rate limiting for authentication endpoints, with an in-memory fallback.
- Health/readiness endpoints, Docker Compose and CI workflow.

## Structure
- cmd/api: server startup and graceful shutdown.
- internal/auth: password hashing and JWT.
- internal/config: environment validation.
- internal/domain: request-scoped identity.
- internal/httpx: HTTP handlers, middleware and routing.
- internal/store: parameterized PostgreSQL operations.
- migrations: schema and indexes.

## Run
Requirements: Go 1.23+, Docker and Docker Compose.

    cp .env.example .env
    docker compose up --build

Apply the database schema before using the API:

    psql "$DATABASE_URL" -f migrations/001_init.sql

API: http://localhost:8080
Liveness: GET /healthz
Readiness: GET /readyz

## API
| Method | Endpoint | Access |
|---|---|---|
| POST | /v1/auth/register | Public |
| POST | /v1/auth/login | Public |
| POST | /v1/organizations | Authenticated |
| GET | /v1/organizations/current | Tenant member |
| GET | /v1/projects | Tenant member |
| POST | /v1/projects | Owner / admin |
| GET | /v1/audit-logs | Owner / admin |
| GET | /v1/billing/subscription | Owner / admin |
| POST | /v1/billing/subscription | Owner |

Register, then create an organization. The organization response includes a tenant-scoped access token. Use that token for tenant endpoints.

## Security
Tenant identity is read from a verified JWT and membership is rechecked against PostgreSQL on every tenant request. Tenant-owned SQL statements include organization_id. Passwords are bcrypt-hashed and SQL is parameterized.

Before production, add email verification, password recovery, refresh-token rotation/revocation, MFA, invitations, payment-provider webhook verification, TLS, secret rotation, backups and observability. The plan-change endpoint is a demo scaffold and does not prove payment.

## Development
    go test ./...
    go vet ./...

CI checks formatting, vet and tests. Tests have not yet been executed in this repository's runtime.

## License
MIT
