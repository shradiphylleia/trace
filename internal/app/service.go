package app

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"traceshare/internal/domain"
)

const maxUploadBytes=25<<20
const defaultSearchLimit=25
const maxSearchLimit=100
const shortCodeAlphabet="abcdefghijklmnopqrstuvwxyz0123456789"

var ErrNotFound=errors.New("artifact not found")

type ArtifactStore interface {
	Create(ctx context.Context,artifact domain.Artifact) error
	FindByShortCode(ctx context.Context, shortCode string) (domain.Artifact, error)
	Search(ctx context.Context, filter SearchFilter) ([]domain.Artifact, error)
	ListExpired(ctx context.Context, now time.Time, limit int) ([]domain.Artifact, error)
	DeleteExpired(ctx context.Context, ids []string) error
}

type FileStore interface {
	Save(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

type Cache interface {
	Get(ctx context.Context, shortCode string) (domain.Artifact, bool)
	Set(ctx context.Context, artifact domain.Artifact) error
}

type SearchFilter struct {
	Query   string
	Service string
	Tag     string
	Limit   int
}

type Service struct {
	artifacts ArtifactStore
	files     FileStore
	cache     Cache
	baseURL   string
}

func NewService(artifacts ArtifactStore, files FileStore, cache Cache, baseURL string) *Service {
	return &Service{artifacts: artifacts, files: files, cache: cache, baseURL: strings.TrimRight(baseURL, "/")}
}

func (s *Service) Create(ctx context.Context, input domain.CreateArtifactInput, reader io.Reader) (domain.Artifact, error) {
	now:=time.Now().UTC()
	if err:=input.Validate(); err != nil {
		return domain.Artifact{}, err
	}

	expiresAt,err:=domain.ExpirationTime(input.Expiration, now)
	if err!=nil{
		return domain.Artifact{}, err
	}

	shortCode,err:=generateShortCode()
	if err!=nil {
		return domain.Artifact{},err
	}
	if input.SizeBytes>maxUploadBytes {
		return domain.Artifact{},errors.New("file exceeds 25MB limit")
	}
	id:=newID()
	ext:=filepath.Ext(input.FileName)
	objectKey:=fmt.Sprintf("artifacts/%s/%s%s",now.Format("2006/01/02"),id,ext)

	var previewBuffer bytes.Buffer
	tee:=io.TeeReader(io.LimitReader(reader,maxUploadBytes+1),&previewBuffer)
	if err:=s.files.Save(ctx,objectKey,tee,input.SizeBytes,input.ContentType);err!=nil{
		return domain.Artifact{}, err
	}
	if previewBuffer.Len()>maxUploadBytes {
		_ =s.files.Delete(ctx,objectKey)
		return domain.Artifact{},errors.New("file exceeds 25MB limit")
	}

	artifact := domain.Artifact{
		ID:          id,
		ShortCode:   shortCode,
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Type:        input.Type,
		ServiceName: strings.TrimSpace(input.ServiceName),
		Environment: strings.TrimSpace(input.Environment),
		Tags:        normalizeTags(input.Tags),
		Creator:     strings.TrimSpace(input.Creator),
		ObjectKey:   objectKey,
		FileName:    cleanFileName(input.FileName),
		ContentType: input.ContentType,
		SizeBytes:   input.SizeBytes,
		CreatedAt:   now,
		ExpiresAt:   expiresAt,
		Preview:     previewText(input.ContentType, previewBuffer.Bytes()),
	}

	if err:=s.artifacts.Create(ctx, artifact); err != nil {
		_ =s.files.Delete(ctx, objectKey)
		return domain.Artifact{}, err
	}
	return s.decorate(artifact), nil
}

func (s *Service) Get(ctx context.Context, shortCode string) (domain.Artifact, error) {
	shortCode=strings.TrimSpace(shortCode)
	if artifact,ok:=s.cache.Get(ctx, shortCode); ok {
		return s.decorate(artifact), nil
	}

	artifact,err:=s.artifacts.FindByShortCode(ctx, shortCode)
	if err != nil {
		return domain.Artifact{}, err
	}
	if artifact.ExpiresAt!=nil&&artifact.ExpiresAt.Before(time.Now().UTC()) {
		return domain.Artifact{}, ErrNotFound
	}
	_=s.cache.Set(ctx, artifact)
	return s.decorate(artifact), nil
}

func (s *Service) Download(ctx context.Context, shortCode string) (domain.Artifact, io.ReadCloser, error) {
	artifact,err:=s.Get(ctx, shortCode)
	if err!=nil{
		return domain.Artifact{}, nil, err
	}
	body, err:=s.files.Get(ctx, artifact.ObjectKey)
	if err!=nil {
		return domain.Artifact{}, nil, err
	}
	return artifact, body, nil
}

func (s *Service) Search(ctx context.Context, filter SearchFilter) ([]domain.Artifact, error) {
	if filter.Limit <= 0 || filter.Limit > maxSearchLimit {
		filter.Limit = defaultSearchLimit
	}
	artifacts,err:=s.artifacts.Search(ctx, filter)
	if err!=nil {
		return nil, err
	}
	for i:=range artifacts {
		artifacts[i]=s.decorate(artifacts[i])
	}
	return artifacts, nil
}

func (s *Service) decorate(artifact domain.Artifact) domain.Artifact {
	artifact.ShareURL=s.baseURL+"/t/"+artifact.ShortCode
	artifact.DownloadURL=s.baseURL+"/api/artifacts/"+artifact.ShortCode+"/download"
	return artifact
}

func normalizeTags(tags []string) []string {
	seen:=map[string]struct{}{}
	result:=make([]string, 0, len(tags))
	for _,tag:=range tags {
		tag=strings.ToLower(strings.TrimSpace(tag))
		if tag=="" {
			continue
		}
		if _,ok:=seen[tag]; ok {
			continue
		}
		seen[tag]=struct{}{}
		result=append(result, tag)
	}
	return result
}

//file ka scene hai kuch
func cleanFileName(name string) string {
	name=filepath.Base(strings.TrimSpace(name))
	if name=="."||name==""{
		return "artifact.txt"
	}
	return name
}

// text preveiw
func previewText(contentType string, data []byte) string {
	if !strings.HasPrefix(contentType,"text/") &&
		!strings.Contains(contentType,"json") &&
		!strings.Contains(contentType,"xml") {
		return ""
	}
	text:=string(data)
	if len(text)>4000 {
		return text[:4000]
	}
	return text
}

func generateShortCode() (string, error) {
	raw:=make([]byte, 6)
	if _, err:=rand.Read(raw); err != nil {
		return "", err
	}
	code:=make([]byte, len(raw))
	for i,value:=range raw {
		code[i]=shortCodeAlphabet[int(value)%len(shortCodeAlphabet)]
	}
	return string(code), nil
}

func newID() string {
	raw:=make([]byte, 16)
	if _, err:=rand.Read(raw); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	const alphabet="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, len(raw))
	for i, value:=range raw {
		id[i]=alphabet[int(value)%len(alphabet)]
	}
	return string(id)
}
