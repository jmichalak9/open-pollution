package main

import (
	"fmt"
	"time"

	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/ratelimiting"
	"github.com/areknoster/public-distributed-commit-log/sentinel/commiter"
	"github.com/areknoster/public-distributed-commit-log/sentinel/pinner"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel/service"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/storage/localfs"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/ipfs/go-cid"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-limiter/memorystore"

	"github.com/jmichalak9/open-pollution/cmd/pdcl"
	"github.com/jmichalak9/open-pollution/cmd/pdcl/sentinel/validator"
)

type Config struct {
	GRPC grpc.ServerConfig
	pdcl.LocalStorageConfig
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	contentStorage, err := localfs.NewStorage(config.Directory)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize storage")
	}
	messageStorage := storage.NewProtoMessageStorage(contentStorage)

	schemaValidator := validator.NewSchemaValidator(messageStorage)
	memoryPinner := pinner.NewMemoryPinner()
	headManager := memory.NewHeadManager(cid.Undef) // initialize it as if it was initializing topic for the first time
	instantCommiter := commiter.NewInstant(headManager, messageStorage, memoryPinner)

	sentinelService := service.New(schemaValidator, memoryPinner, instantCommiter, headManager)

	ratelimiter, err := setupRateLimiter(config.GRPC.RPS)
	if err != nil {
		log.Fatal().Err(err).Msg("can't setup rate limiter")
	}

	grpcServer, err := grpc.NewServer(config.GRPC, ratelimiter)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize grpc server")
	}
	sentinelpb.RegisterSentinelServer(grpcServer, sentinelService)
	log.Fatal().Err(grpcServer.ListenAndServe()).Msg("error running grpc server")
}

func setupRateLimiter(rps int) (ratelimit.Limiter, error) {
	var limiter ratelimit.Limiter
	if rps > 0 {
		store, err := memorystore.New(&memorystore.Config{
			Tokens:   uint64(rps),
			Interval: time.Second,
		})
		if err != nil {
			return nil, fmt.Errorf("create rate limiter: %w", err)
		}
		limiter = ratelimiting.NewTokenBucketLimiter(store)
	} else {
		limiter = ratelimiting.NewAlwaysAllowLimiter()
	}

	return limiter, nil
}
