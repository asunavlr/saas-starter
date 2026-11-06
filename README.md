# SaaS Starter

A production-minded multi-tenant SaaS foundation built with Go, PostgreSQL and Redis.

> **Status:** MVP foundation. Authentication, tenant isolation, RBAC, projects, audit events, subscription entitlements and rate limiting are implemented as a modular starting point. Payment-provider integration, email delivery and refresh-token rotation are intentionally left as extension points.

## Features

- Multi-tenant organizations with membership-based access control
- JWT access tokens and bcrypt password hashing
- Role-based authorization: `owner`, `admin`, `member`
- Tenant-scoped project CRUD
- PostgreSQL migrations, foreign keys and indexes
- Audit trail for security-sensitive and business actions
- Plan entitlements and per-tenant project limits
- Redis-backed IP rate limiting (in-memory fallback for local development)
- Structured HTTP middleware, request IDs, graceful shutdown
- Docker Compose, health checks, OpenAPI contract and CI

## Architecture

```text
cmd/api            HTTP server and dependency wiring
internal/
  auth             Password hashing and JWT issuance
  config           Environment configuration
  http             Routes, handlers and middleware
  platform         PostgreSQL, Redis and rate limiter
  tenancy          Organization membership and authorization
  projects         Tenant-scoped project use cases
  audit            Audit event persistence
  billing          Plans and entitlements
migrations/        PostgreSQL schema
api/openapi.yaml   API contract
```

## Quick start

Requirements: Go 1.23+, Docker and Docker Compose.

```bash
cp .env.example .env
docker compose up --build
```

API: `http://localhost:8080`  
Health: `GET /healthz`  
Readiness: `GET /readyz`

## Configuration

| Variable | Description | Default |
|---|---|---|
| `HTTP_ADDR` | HTTP listen address | `:8080` |
| `DATABASE_URL` | PostgreSQL connection string | required |
| `REDIS_ADDR` | Redis host:port | `localhost:6379` |
| `JWT_SECRET` | Signing secret (32+ chars recommended) | required |
| `ACCESS_TOKEN_TTL` | Access token lifetime | `15m` |
| `RATE_LIMIT_RPM` | Requests per IP per minute | `120` |
| `APP_ENV` | Runtime environment | `development` |

## API flow

1. `POST /v1/auth/register` creates a user, an organization and its owner membership.
2. `POST /v1/auth/login` returns a short-lived JWT.
3. Send `Authorization: Bearer <token>` and `X-Organization-ID: <uuid>` to tenant routes.
4. `GET/POST /v1/projects` operates only inside the selected organization.
5. Organization owners/admins can manage projects; members can read them.

Example registration:

```bash
curl -X POST localhost:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Kevin","email":"kevin@example.com","password":"correct-horse-battery"}'
```

## Security notes

- Never use the development JWT secret in production.
- Terminate TLS at a trusted ingress and set secure deployment/network policies.
- Add email verification, password reset, refresh-token rotation and MFA before production launch.
- Payment state must be updated from verified provider webhooks; do not trust client-supplied plan IDs.
- Apply outbound network restrictions if user-controlled URLs are introduced.
- Use a secrets manager and rotate credentials.
- Run migrations as a deployment step, not from every API replica.

## Development

```bash
go test ./...
go vet ./...
```

The integration tests expect PostgreSQL and Redis when those services are enabled. Unit tests should remain isolated from external services.

## Roadmap

- [ ] Stripe/Paddle billing adapter and signed webhook processing
- [ ] Refresh-token rotation and session revocation
- [ ] Email verification and password recovery
- [ ] API keys with scoped permissions
- [ ] Background jobs and transactional outbox
- [ ] OpenTelemetry traces and metrics
- [ ] Admin dashboard / Next.js starter

## License

MIT
