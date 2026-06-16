package worker
import (
	"context"
	"time"
	"github.com/rs/zerolog/log"
	"traceshare/internal/app"
)
type Cleaner struct {
	artifacts app.ArtifactStore
	files app.FileStore
	interval time.Duration
}
func NewCleaner(artifacts app.ArtifactStore, files app.FileStore, interval time.Duration) *Cleaner {
	return &Cleaner{artifacts: artifacts, files: files, interval: interval}
}
func (c *Cleaner) Start(ctx context.Context) {
	if c.interval<=0{
		c.interval=15*time.Minute
	}
	c.clean(ctx)
	ticker:=time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select{
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.clean(ctx)
		}
	}
}
func (c *Cleaner) clean(ctx context.Context) {
	artifacts,err:=c.artifacts.ListExpired(ctx, time.Now().UTC(), 100)
	if err!=nil{
		log.Error().Err(err).Msg("list expired artifacts")
		return
	}
	if len(artifacts)==0{
		return
	}
	ids:=make([]string, 0, len(artifacts))
	for _,artifact:=range artifacts {
		if err:=c.files.Delete(ctx, artifact.ObjectKey); err != nil {
			log.Warn().Err(err).Str("object_key", artifact.ObjectKey).Msg("delete expired object")
		}
		ids=append(ids,artifact.ID)
	}
	if err:=c.artifacts.DeleteExpired(ctx,ids);err!=nil {
		log.Error().Err(err).Msg("delete expired artifact rows")
		return
	}
	log.Info().Int("count",len(ids)).Msg("cleaned expired artifacts")
}
