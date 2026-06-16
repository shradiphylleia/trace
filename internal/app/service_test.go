package app

import (
	"bytes"
	"context"
	"io"
	"slices"
	"testing"
	"time"

	"traceshare/internal/domain"
)

func TestServiceCreateNormalizesAndDecorates(t *testing.T) {
	repo:=&fakeRepo{}
	store:=&fakeStore{objects: map[string][]byte{}}
	cache:=&fakeCache{items: map[string]domain.Artifact{}}
	service:=NewService(repo, store, cache, "http://trace.local")

	artifact, err := service.Create(context.Background(), domain.CreateArtifactInput{
		Title:       " Checkout 500 ",
		Type:        domain.ArtifactLog,
		ServiceName: "payments",
		Environment: "staging",
		Tags:        []string{" Sev2 ", "sev2", "checkout"},
		Creator:     "qa@oracle.com",
		Expiration:  "14d",
		FileName:    "trace.log",
		ContentType: "text/plain",
		SizeBytes:   11,
	}, bytes.NewBufferString("hello error"))
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if artifact.Title != "Checkout 500" {
		t.Fatalf("title was not trimmed: %q", artifact.Title)
	}
	if !slices.Equal(artifact.Tags, []string{"sev2", "checkout"}) {
		t.Fatalf("tags not normalized: %#v", artifact.Tags)
	}
	if artifact.ShareURL == "" || artifact.DownloadURL == "" {
		t.Fatalf("expected decorated URLs: %#v", artifact)
	}
	if repo.created.ShortCode == "" {
		t.Fatal("expected repository create")
	}
}

type fakeRepo struct {
	created domain.Artifact
}

func (f *fakeRepo) Create(_ context.Context, artifact domain.Artifact) error {
	f.created = artifact
	return nil
}

func (f *fakeRepo) FindByShortCode(_ context.Context, _ string) (domain.Artifact, error) {
	return f.created, nil
}

func (f *fakeRepo) Search(_ context.Context, _ SearchFilter) ([]domain.Artifact, error) {
	return []domain.Artifact{f.created}, nil
}

func (f *fakeRepo) ListExpired(_ context.Context, _ time.Time, _ int) ([]domain.Artifact, error) {
	return nil, nil
}

func (f *fakeRepo) DeleteExpired(_ context.Context, _ []string) error {
	return nil
}

type fakeStore struct {
	objects map[string][]byte
}

func (f *fakeStore) Save(_ context.Context, key string, reader io.Reader, _ int64, _ string) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	f.objects[key] = data
	return nil
}

func (f *fakeStore) Get(_ context.Context, key string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(f.objects[key])), nil
}

func (f *fakeStore) Delete(_ context.Context, key string) error {
	delete(f.objects, key)
	return nil
}

type fakeCache struct {
	items map[string]domain.Artifact
}

func (f *fakeCache) Get(_ context.Context, shortCode string) (domain.Artifact, bool) {
	artifact, ok := f.items[shortCode]
	return artifact, ok
}

func (f *fakeCache) Set(_ context.Context, artifact domain.Artifact) error {
	f.items[artifact.ShortCode] = artifact
	return nil
}
