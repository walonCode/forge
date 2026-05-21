# Security Policy

## Important

Forge is a personal passion project released under the MIT License. It is provided **as is, with no guarantees**. If you use this code in your own project, you are responsible for your own security — audit it, harden it, and make it fit for your use case.

No formal security support is offered. Vulnerability reports may be read but there is no commitment to respond or patch.

---

## Scope

The following are in scope:

- Authentication and JWT handling (`apps/api/internals/modules/auth/`)
- Authorization — access control between users
- SQL injection or data exposure via API endpoints
- Secrets leaking via logs, headers, or error responses

The following are out of scope:

- Vulnerabilities in third-party dependencies that have no known exploit in this project's context
- Issues that require physical access to a machine
- Denial-of-service via resource exhaustion (rate limiting is in scope, general DoS is not)

---

## Security practices in this codebase

- Passwords hashed with `bcrypt` before storage — plaintext passwords are never stored or logged
- JWTs signed with a secret loaded from environment — never hardcoded
- Database connection string loaded from environment — never committed
- Correlation IDs in logs — no user PII in structured log output
- Return `404` instead of `403` when a user tries to access another user's resource — avoids confirming resource existence
- `.env` is in `.gitignore` — never committed
