# Chirpy
boot.dev course project

Base URL: `http://localhost:8080`

## Endpoint namespaces

- `/app/` is for static assets
- `/admin/` mostly for metrics. Format: `html`
- `/api/` REST API. Format: `json`

## API endpoints

- `/api/chirps` methods: GET, POST
- `/api/chirps/{id}` methods: GET
- `/api/users` methods: POST
- `/api/login` methods: POST

## Development

use `--debug` flag when running the command to automatically erase the database on startup.

## Environment

`JWT_SECRET` can be generated using `openssl rand -base64 64`

## Testing

`tests/intergration` tests are run using JetBrain's GoLand HTTP request module