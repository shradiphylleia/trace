package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
)

type Store struct {
	client *minio.Client
	bucket string
}

func NewStore(client *minio.Client, bucket string) *Store {
	return &Store{client:client,bucket:bucket}
}

func (s *Store) EnsureBucket(ctx context.Context) error {
	exists,err:=s.client.BucketExists(ctx,s.bucket)
	if err!=nil{
		return err
	}
	if exists{
		return nil
	}
	return s.client.MakeBucket(ctx,s.bucket,minio.MakeBucketOptions{})
}

func (s *Store) Save(ctx context.Context,key string,reader io.Reader,size int64,contentType string) error {
	_, err:=s.client.PutObject(ctx,s.bucket,key, reader,size,minio.PutObjectOptions{ContentType:contentType})
	return err
}

func (s *Store) Get(ctx context.Context, key string)(io.ReadCloser,error){
	return s.client.GetObject(ctx, s.bucket,key,minio.GetObjectOptions{})
}

func (s *Store) Delete(ctx context.Context,key string) error {
	return s.client.RemoveObject(ctx,s.bucket,key,minio.RemoveObjectOptions{})
}
