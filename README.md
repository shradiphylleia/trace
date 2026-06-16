# Trace Share
sharing debugging artifacts through one short URL.
teams across can upload stack traces, logs, API payloads, validation reports and screenshots, then share links like:

```text
http://localhost:8080/t/abc123
```

## Features
- Upload text or binary artifacts with metadata.
- Store metadata in PostgreSQL.
- Store artifact bodies and screenshots in MinIO.
- Cache short link lookups in Redis.
- Search by title, service, tags and artifact text.
- Expire artifacts after 7 days, 14 days or never expire.
- Background cleanup worker removes expired rows and objects.

## Architecture

```text
cmd/server              HTTP entry point
internal/domain         Core entities and validation
internal/app            Application workflows
internal/httpapi        REST and simple HTML handlers
internal/db             PostgreSQL persistence
internal/storage        MinIO object storage adapter
internal/cache          Redis cache adapter
internal/worker         Expiration cleanup worker
migrations              PostgreSQL schema
```

## Requirements
- Go 1.22+
- Docker and Docker Compose
## Run Locally
```powershell
docker compose up --build
```
The API and web UI will be available at:

```text
http://localhost:8080
```

## Run The React Frontend
The Vite frontend lives in `frontend/` and proxies API requests to the Go service on port `8080`.

```powershell
cd frontend
npm install
npm run dev
```

Then open:
```text
http://localhost:5173
```

MinIO console is available at:
```text
http://localhost:9001
```

Default MinIO credentials:
```text
minioadmin / minioadmin
```

## API

### Upload Artifact

```bash
curl -X POST http://localhost:8080/api/artifacts \
  -F "title=Checkout 500" \
  -F "description=Reservation service stack trace from QA" \
  -F "artifact_type=stack_trace" \
  -F "service_name=reservation" \
  -F "environment=staging" \
  -F "tags=checkout,sev2,company-internal" \
  -F "creator=someone@shraddha.com" \
  -F "expiration=7d" \
  -F "file=@./trace.txt"
```

Response:
```json
{
  "short_url": "http://localhost:8080/t/abc123",
  "short_code": "abc123"
}
```

### Get Artifact Metadata
```bash
curl http://localhost:8080/api/artifacts/abc123
```

### Download Artifact File
```bash
curl -OJ http://localhost:8080/api/artifacts/abc123/download
```

### Search
```bash
curl "http://localhost:8080/api/search?q=checkout&service=payments&tag=sev2"
```

## Verification
Backend:

```powershell
go test ./cmd/server ./internal/...
```

Frontend:

```powershell
cd frontend
npm run lint
npm run build
```

## Expiration Options
Use one of:
- `7d`
- `14d`
- `never`

Expired artifacts are hidden from lookup and removed by the cleanup worker.
