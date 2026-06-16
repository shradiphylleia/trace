package main
import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"traceshare/internal/app"
	"traceshare/internal/cache"
	"traceshare/internal/config"
	"traceshare/internal/db"
	"traceshare/internal/httpapi"
	"traceshare/internal/storage"
	"traceshare/internal/worker"
)

func main() {
	zerolog.TimeFieldFormat=time.RFC3339
	cfg:=config.Load()
	ctx,stop:=signal.NotifyContext(context.Background(),syscall.SIGINT,syscall.SIGTERM)
	defer stop()

	pool,err:=pgxpool.New(ctx, cfg.DatabaseURL)
	if err!=nil{
		log.Fatal().Err(err).Msg("connect postgres")
	}
	defer pool.Close()

	if err:=pool.Ping(ctx);err!=nil{
		log.Fatal().Err(err).Msg("ping postgres")
	}

	redisClient:=redis.NewClient(&redis.Options{Addr:cfg.RedisAddr})
	defer redisClient.Close()

	minioClient,err:=minio.New(cfg.MinIOEndpoint,&minio.Options{
		Creds:credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOSecretKey, ""),
		Secure:cfg.MinIOUseSSL,
	})
	if err!=nil {
		log.Fatal().Err(err).Msg("connect minio")
	}
	objectStore:=storage.NewStore(minioClient, cfg.MinIOBucket)
	if err:=objectStore.EnsureBucket(ctx); err != nil {
		log.Fatal().Err(err).Msg("ensure minio bucket")
	}

	artifactStore:=db.NewArtifactStore(pool)
	artifactCache:=cache.NewStore(redisClient)
	service:=app.NewService(artifactStore, objectStore, artifactCache, cfg.BaseURL)

	cleaner:=worker.NewCleaner(artifactStore, objectStore, cfg.CleanupInterval)
	go cleaner.Start(ctx)

	router:=chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:[]string{"*"},
		AllowedMethods:[]string{"GET","POST","OPTIONS"},
		AllowedHeaders:[]string{"Accept","Authorization","Content-Type"},
		MaxAge:300,
	}))

	httpapi.NewHandler(service).Register(router)
	server:=&http.Server{
		Addr:":"+cfg.Port,
		Handler:router,
		ReadHeaderTimeout:5*time.Second,
	}

	go func() {
		log.Info().Str("addr", server.Addr).Msg("starting application")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("http server")
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err:=server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("shutdown server")
		os.Exit(1)
	}
}
