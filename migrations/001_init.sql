CREATE TABLE IF NOT EXISTS artifacts (
    id text PRIMARY KEY,
    short_code text NOT NULL UNIQUE,
    title text NOT NULL,
    description text NOT NULL DEFAULT '',
    artifact_type text NOT NULL CHECK (artifact_type IN ('stack_trace', 'log', 'api_payload', 'validation_report', 'screenshot')),
    service_name text NOT NULL,
    environment text NOT NULL,
    tags text[] NOT NULL DEFAULT '{}',
    creator text NOT NULL,
    object_key text NOT NULL,
    file_name text NOT NULL,
    content_type text NOT NULL,
    size_bytes bigint NOT NULL CHECK (size_bytes > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz,
    preview text NOT NULL DEFAULT '',
    search_document tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('simple', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(service_name, '')), 'A') ||
        setweight(to_tsvector('simple', coalesce(array_to_string(tags, ' '), '')), 'B') ||
        setweight(to_tsvector('simple', coalesce(description, '')), 'C') ||
        setweight(to_tsvector('simple', coalesce(preview, '')), 'D')
    ) STORED
);

CREATE INDEX IF NOT EXISTS idx_artifacts_service_name ON artifacts (service_name);
CREATE INDEX IF NOT EXISTS idx_artifacts_tags ON artifacts USING gin (tags);
CREATE INDEX IF NOT EXISTS idx_artifacts_search ON artifacts USING gin (search_document);
CREATE INDEX IF NOT EXISTS idx_artifacts_expires_at ON artifacts (expires_at) WHERE expires_at IS NOT NULL;
