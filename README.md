# SaaS Starter

A multi-tenant SaaS backend foundation built with **Go, PostgreSQL and Redis**. The project focuses on backend fundamentals that matter in real products: tenant isolation, authorization, quotas, auditability and operational readiness.

## Features

- **Authentication:** registration and login with bcrypt password hashing and signed JWT access tokens.
- **Tenant provisioning:** registration creates a user, a workspace and an owner membership in one PostgreSQL transaction.
- **Multi-tenancy:** clients select an organization with the X-Organization-ID header; membership is verified against PostgreSQL on every tenant request.
- **RBAC:** owner/admin can write; members have read access to tenant projects.
- **Projects:** tenant-scoped list, create and delete operations.
- **Plan entitlements:** Free, Pro and Business project/member limits; project quotas are enforced by the API.
- **Audit trail:** project creation/deletion events with tenant, actor, resource and metadata; tenant-scoped audit history endpoint.
- **Billing foundation:** subscription status and plan entitlement overview. No payment provider is connected.
- **Rate limiting:** Redis fixed-window IP limiter with a process-local fallback for development.
- **Operations:** PostgreSQL migrations, Docker Compose, liveness/readiness probes, request IDs, structured request logging and graceful shutdown.

## Architecture

- cmd/api: dependency wiring, HTTP routes and graceful shutdown.
- internal/auth: bcrypt and JWT helpers.
- internal/audit: audit event persistence and tenant-scoped reads.
- internal/billing: plan catalog and entitlement limits.
- internal/config: validated environment configuration.
- internal/httpx: handlers, middleware and request context.
- internal/platform: PostgreSQL connection pool and Redis rate limiter.
- internal/projects: tenant-scoped project repository.
- internal/tenancy: organization membership authorization.
- migrations: PostgreSQL schema and indexes.
- api/openapi.yaml: API contract.

## Run locally

Requirements: Go 1.23+, Docker and Docker Compose.

    cp .env.example .env
    docker compose up --build

The Compose setup starts the API, PostgreSQL and Redis. PostgreSQL applies the SQL files mounted into its initialization directory on first startup. For an existing database, apply migrations manually:

    psql "$DATABASE_URL" -f migrations/001_init.sql

API: http://localhost:8080  
Liveness: GET /healthz  
Readiness: GET /readyz

## Configuration

See .env.example. DATABASE_URL and a JWT_SECRET of at least 32 characters are required. Redis is configured through REDIS_ADDR and REDIS_PASSWORD. Never use the development secret in production.

## API

All endpoints are under /v1 unless noted.

| Method | Endpoint | Access | Description |
|---|---|---|---|
| POST | /auth/register | Public | Create account + owner workspace |
| POST | /auth/login | Public | Authenticate and list memberships |
| GET | /me | Bearer token | Current user profile |
| GET | /projects | Tenant member | List projects in selected tenant |
| POST | /projects | Owner/admin | Create a project |
| DELETE | /projects/{id} | Owner/admin | Delete a project |
| GET | /audit-events | Owner/admin | Read latest tenant audit events |
| GET | /billing | Owner/admin | Read subscription status and entitlements |

Tenant routes require both headers:

    Authorization: Bearer <access-token>
    X-Organization-ID: <organization-uuid>

Registration example:

    curl -X POST http://localhost:8080/v1/auth/register \
      -H 'Content-Type: application/json' \
      -d '{"name":"Alex","email":"alex@example.com","password":"a-strong-password"}'

The response includes an access token and the newly created workspace ID. Use both headers above to call tenant endpoints.

## Security and production considerations

- Passwords are bcrypt-hashed; SQL statements are parameterized.
- Tenant-owned queries include the selected organization ID, and membership is revalidated from the database.
- JWTs use HS256 and an expiration claim. Rotate the secret and use a managed secret store in production.
- The in-memory rate-limit fallback is process-local; use Redis in multi-replica deployments.
- Add email verification, password recovery, refresh-token rotation/revocation, MFA, invitations and stronger session controls before production.
- Billing is an entitlement scaffold only. Paid access must be activated by verified payment-provider webhooks, never by trusting a client request.
- Add transactional audit guarantees, OpenTelemetry, backups, TLS and deployment-time migrations before production.

## Development

    go test ./...
    go vet ./...

GitHub Actions runs tests and vet. Tests are configured in the repository but have not been executed in this authoring session.

## License

MIT
