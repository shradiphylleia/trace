# TraceShare Learning Map

This file maps the main pieces of the TraceShare MVP and explains what each one does.

## Product Flow

TraceShare lets an engineer upload a debugging artifact, stores the artifact file and metadata, generates a short code, and exposes everything through a short URL.

Example flow:

```text
Upload file + metadata
  -> store metadata in PostgreSQL
  -> store file in MinIO
  -> generate short code
  -> share /t/{shortCode}
  -> search or download later
```

## Root Files

| File | Purpose |
| --- | --- |
| `README.md` | User-facing project overview, run instructions, and API examples. |
| `go.mod` / `go.sum` | Go module definition and dependency lock data. |
| `Dockerfile` | Builds the Go service into a small container image. |
| `docker-compose.yml` | Runs the full local stack: app, PostgreSQL, Redis, and MinIO. |
| `.dockerignore` | Keeps unnecessary files out of the Docker build context. |
| `frontend/` | React, TypeScript, Vite, Tailwind CSS frontend. |
| `learning.md` | This map of what was built. |

## Frontend

| Path | Purpose |
| --- | --- |
| `frontend/src/App.tsx` | Small route controller for dashboard, upload, details, and search views. |
| `frontend/src/api/client.ts` | Fetch-based API integration layer for upload, lookup, and search. |
| `frontend/src/api/types.ts` | TypeScript types matching the Go API responses. |
| `frontend/src/components/` | Shared app shell, tables, search box, state views, and shadcn-style UI primitives. |
| `frontend/src/pages/dashboard.tsx` | Dashboard with project title, upload action, search, and recent artifacts. |
| `frontend/src/pages/upload.tsx` | Upload form for metadata, expiration policy, and file selection. |
| `frontend/src/pages/artifact-details.tsx` | Metadata, share URL, copy link, download, expiration status, and preview. |
| `frontend/src/pages/search-results.tsx` | Search results list with title, type, service, created date, and preview snippet. |

Run the frontend with:

```powershell
cd frontend
npm install
npm run dev
```

The frontend expects the Go API at:

```text
http://localhost:8080
```

## Backend Entry Point

| Path | Purpose |
| --- | --- |
| `cmd/server/main.go` | Starts the HTTP server, connects to PostgreSQL, Redis, and MinIO, registers routes, and starts the expiration cleanup worker. |

## Configuration

| Path | Purpose |
| --- | --- |
| `internal/config/config.go` | Reads runtime configuration from environment variables with local defaults. |

Important environment variables:

| Variable | Purpose |
| --- | --- |
| `PORT` | HTTP server port. |
| `BASE_URL` | Base URL used when generating share and download links. |
| `DATABASE_URL` | PostgreSQL connection string. |
| `REDIS_ADDR` | Redis host and port. |
| `MINIO_ENDPOINT` | MinIO API endpoint. |
| `MINIO_ACCESS_KEY` / `MINIO_SECRET_KEY` | MinIO credentials. |
| `MINIO_BUCKET` | Bucket where artifact files are stored. |
| `CLEANUP_INTERVAL` | How often expired artifacts are cleaned up. |

## Domain Layer

| Path | Purpose |
| --- | --- |
| `internal/domain/artifact.go` | Defines the core `Artifact` model, allowed artifact types, upload input validation, and expiration policy logic. |
| `internal/domain/artifact_test.go` | Tests expiration and artifact validation behavior. |

Supported artifact types:

- `stack_trace`
- `log`
- `api_payload`
- `validation_report`
- `screenshot`

Supported expiration policies:

- `7d`
- `14d`
- `never`

## Use Case Layer

| Path | Purpose |
| --- | --- |
| `internal/app/service.go` | Contains the main application workflow: create artifact, get by short code, download, search, decorate URLs, normalize tags, and generate short codes. |
| `internal/app/service_test.go` | Tests the artifact creation workflow with fake repository, storage, and cache implementations. |

Key responsibilities:

- Validate upload metadata.
- Enforce the 25 MB upload limit.
- Generate a six-character short code.
- Create object storage keys.
- Save the file to object storage.
- Save metadata to PostgreSQL.
- Read from Redis cache for short-code lookups.
- Hide expired artifacts.

## HTTP Handler Layer

| Path | Purpose |
| --- | --- |
| `internal/httpapi/handler.go` | Defines REST endpoints and the simple HTML web UI. |

Routes:

| Route | Method | Purpose |
| --- | --- | --- |
| `/` | `GET` | Upload/search web UI. |
| `/healthz` | `GET` | Health check endpoint. |
| `/api/artifacts` | `POST` | Upload an artifact and create a short link. |
| `/api/artifacts/{shortCode}` | `GET` | Fetch artifact metadata. |
| `/api/artifacts/{shortCode}/download` | `GET` | Download the stored artifact file. |
| `/api/search` | `GET` | Search artifacts by query, service, or tag. |
| `/t/{shortCode}` | `GET` | Human-readable short-link page. |

## Repository Layer

| Path | Purpose |
| --- | --- |
| `internal/db/artifact_store.go` | PostgreSQL implementation for creating, finding, searching, listing expired, and deleting artifact rows. |

Search supports:

- Title text
- Service name
- Tags
- Description text
- Artifact preview text

## Cache Layer

| Path | Purpose |
| --- | --- |
| `internal/cache/store.go` | Redis cache implementation for short-code artifact lookups. |

The cache stores internal artifact data, including the MinIO object key, while public API responses still hide the object key.

## Object Storage Layer

| Path | Purpose |
| --- | --- |
| `internal/storage/store.go` | MinIO implementation for saving, reading, deleting files, and creating the bucket if needed. |

Artifact files are stored in MinIO under keys like:

```text
artifacts/2026/06/16/{artifactID}.log
```

## Background Worker

| Path | Purpose |
| --- | --- |
| `internal/worker/cleaner.go` | Periodically finds expired artifacts, deletes their MinIO objects, and removes their PostgreSQL rows. |

## Database

| Path | Purpose |
| --- | --- |
| `migrations/001_init.sql` | Creates the `artifacts` table and indexes. |

The `artifacts` table stores:

- Short code
- Title
- Description
- Artifact type
- Service name
- Environment
- Tags
- Creator
- MinIO object key
- File metadata
- Created timestamp
- Expiration timestamp
- Text preview
- Generated PostgreSQL search document

## Local Deployment

`docker-compose.yml` starts:

| Service | Purpose |
| --- | --- |
| `app` | TraceShare Go service. |
| `postgres` | Metadata database. |
| `redis` | Short-link lookup cache. |
| `minio` | Object storage for uploaded artifacts. |

Run with:

```powershell
docker compose up --build
```

Then open:

```text
http://localhost:8080
```

## Verification Done

The Go test suite passed with:

```powershell
go test ./cmd/server ./internal/...
```

Docker Compose configuration was validated with:

```powershell
docker compose config
```

The stack was not started because Docker Desktop was not running on this machine.

The frontend was checked with:

```powershell
cd frontend
npm run lint
npm run build
npm audit --omit=dev
```
