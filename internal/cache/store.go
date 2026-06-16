package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"traceshare/internal/domain"
)
type cachedArtifact struct {
	ID          string `json:"id"`
	ShortCode   string `json:"short_code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        domain.ArtifactType `json:"artifact_type"`
	ServiceName string `json:"service_name"`
	Environment string `json:"environment"`
	Tags        []string `json:"tags"`
	Creator     string  `json:"creator"`
	ObjectKey   string  `json:"object_key"`
	FileName    string  `json:"file_name"`
	ContentType string  `json:"content_type"`
	SizeBytes   int64   `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Preview     string `json:"preview,omitempty"`
}

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client:client}
}

func (s *Store) Get(ctx context.Context, shortCode string) (domain.Artifact, bool) {
	data,err:=s.client.Get(ctx,key(shortCode)).Bytes()
	if err!=nil{
		return domain.Artifact{}, false
	}
	var cached cachedArtifact
	if err:=json.Unmarshal(data,&cached);err!=nil{
		return domain.Artifact{},false
	}
	return fromCached(cached),true
}

func (s *Store) Set(ctx context.Context, artifact domain.Artifact) error {
	ttl:=15*time.Minute
	if artifact.ExpiresAt!=nil {
		untilExpiry:=time.Until(*artifact.ExpiresAt)
		if untilExpiry<=0{
			return nil
		}
		if untilExpiry<ttl {
			ttl=untilExpiry
		}
	}
	data,err:=json.Marshal(toCached(artifact))
	if err!=nil {
		return err
	}
	return s.client.Set(ctx, key(artifact.ShortCode), data, ttl).Err()
}

func key(shortCode string) string {
	return "artifact:"+shortCode
}

func toCached(artifact domain.Artifact) cachedArtifact {
	return cachedArtifact{
		ID:artifact.ID,
		ShortCode:artifact.ShortCode,
		Title:artifact.Title,
		Description:artifact.Description,
		Type:artifact.Type,
		ServiceName: artifact.ServiceName,
		Environment: artifact.Environment,
		Tags:artifact.Tags,
		Creator:artifact.Creator,
		ObjectKey:artifact.ObjectKey,
		FileName:artifact.FileName,
		ContentType:artifact.ContentType,
		SizeBytes:artifact.SizeBytes,
		CreatedAt:artifact.CreatedAt,
		ExpiresAt:artifact.ExpiresAt,
		Preview:artifact.Preview,
	}
}

func fromCached(cached cachedArtifact) domain.Artifact {
	return domain.Artifact{
		ID:cached.ID,
		ShortCode:cached.ShortCode,
		Title:cached.Title,
		Description:cached.Description,
		Type:cached.Type,
		ServiceName:cached.ServiceName,
		Environment:cached.Environment,
		Tags:cached.Tags,
		Creator:cached.Creator,
		ObjectKey:cached.ObjectKey,
		FileName:cached.FileName,
		ContentType:cached.ContentType,
		SizeBytes:cached.SizeBytes,
		CreatedAt:cached.CreatedAt,
		ExpiresAt:cached.ExpiresAt,
		Preview:cached.Preview,
	}
}
