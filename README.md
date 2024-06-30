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

## Development

use `--debug` flag when running the command to automatically erase the database on startup.

## Testing

`tests/intergration` tests are run using JetBrain's GoLand HTTP request module