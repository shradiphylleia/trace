package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"traceshare/internal/app"
	"traceshare/internal/domain"
)

type ArtifactStore struct {
	db *pgxpool.Pool
}

func NewArtifactStore(pool *pgxpool.Pool) *ArtifactStore {
	return &ArtifactStore{db: pool}
}

func (s *ArtifactStore) Create(ctx context.Context, artifact domain.Artifact) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO artifacts (
			id, short_code, title, description, artifact_type, service_name, environment,
			tags, creator, object_key, file_name, content_type, size_bytes, created_at, expires_at, preview
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, artifact.ID, artifact.ShortCode, artifact.Title, artifact.Description, artifact.Type,
		artifact.ServiceName, artifact.Environment, artifact.Tags, artifact.Creator, artifact.ObjectKey,
		artifact.FileName, artifact.ContentType, artifact.SizeBytes, artifact.CreatedAt, artifact.ExpiresAt,
		artifact.Preview)
	return err
}

func (s *ArtifactStore) FindByShortCode(ctx context.Context, shortCode string) (domain.Artifact, error) {
	row := s.db.QueryRow(ctx, baseSelect()+`
		WHERE short_code = $1
		  AND (expires_at IS NULL OR expires_at > now())
	`, shortCode)
	return scanArtifact(row)
}

func (s *ArtifactStore) Search(ctx context.Context, filter app.SearchFilter) ([]domain.Artifact, error) {
	query := baseSelect() + `
		WHERE (expires_at IS NULL OR expires_at > now())
		  AND ($1 = '' OR search_document @@ plainto_tsquery('simple', $1) OR title ILIKE '%' || $1 || '%' OR preview ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR service_name ILIKE $2)
		  AND ($3 = '' OR tags @> ARRAY[$3]::text[])
		ORDER BY created_at DESC
		LIMIT $4
	`
	rows, err := s.db.Query(ctx, query, strings.TrimSpace(filter.Query), servicePattern(filter.Service), strings.ToLower(strings.TrimSpace(filter.Tag)), filter.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artifacts:= make([]domain.Artifact, 0)
	for rows.Next() {
		artifact, err := scanArtifact(rows)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, rows.Err()
}

func (s *ArtifactStore) ListExpired(ctx context.Context, now time.Time, limit int) ([]domain.Artifact, error) {
	rows, err := s.db.Query(ctx, baseSelect()+`
		WHERE expires_at IS NOT NULL
		  AND expires_at <= $1
		ORDER BY expires_at ASC
		LIMIT $2
	`, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artifacts []domain.Artifact
	for rows.Next() {
		artifact, err := scanArtifact(rows)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, rows.Err()
}

func (s *ArtifactStore) DeleteExpired(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := s.db.Exec(ctx, `DELETE FROM artifacts WHERE id = ANY($1)`, ids)
	return err
}

func baseSelect() string {
	return `
		SELECT id, short_code, title, description, artifact_type, service_name, environment,
			tags, creator, object_key, file_name, content_type, size_bytes, created_at, expires_at, preview
		FROM artifacts
	`
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanArtifact(row rowScanner) (domain.Artifact, error) {
	var artifact domain.Artifact
	err := row.Scan(&artifact.ID, &artifact.ShortCode, &artifact.Title, &artifact.Description,
		&artifact.Type, &artifact.ServiceName, &artifact.Environment, &artifact.Tags, &artifact.Creator,
		&artifact.ObjectKey, &artifact.FileName, &artifact.ContentType, &artifact.SizeBytes,
		&artifact.CreatedAt, &artifact.ExpiresAt, &artifact.Preview)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Artifact{}, app.ErrNotFound
	}
	return artifact, err
}

func servicePattern(service string) string {
	service = strings.TrimSpace(service)
	if service == "" {
		return ""
	}
	return "%" + service + "%"
}
