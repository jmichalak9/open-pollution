package main

import (
	"fmt"
	"time"

	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/ratelimiting"
	"github.com/areknoster/public-distributed-commit-log/sentinel/commiter"
	"github.com/areknoster/public-distributed-commit-log/sentinel/pinner"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel/service"
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	memoryhead "github.com/areknoster/public-distributed-commit-log/thead/memory"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-limiter/memorystore"

	"github.com/jmichalak9/open-pollution/cmd/pdcl/sentinel/internal/validator"
)

type Config struct {
	DaemonStorage ipfsstorage.Config
	Validator     validator.Config
	GRPC          grpc.ServerConfig
	Commiter      commiter.MaxBufferCommiterConfig
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	codec := pbcodec.Json{}

	ipfsShell := ipfsstorage.NewShell(config.DaemonStorage)
	storage := ipfsstorage.NewStorage(ipfsShell, codec)

	messageValidator, err := validator.New(storage, codec, config.Validator)
	if err != nil {
		log.Fatal().Err(err).Msg("initialize message validator")
	}
	memPinner := pinner.NewMemoryPinner()
	headManager := memoryhead.NewHeadManager(cid.Undef)
	ipnsManager, err := setupIPNSManager(config, ipfsShell)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't set up ipns manager")
	}
	instantCommiter := commiter.NewMaxBufferCommitter(
		headManager,
		storage,
		memPinner,
		ipnsManager,
		config.Commiter)
	sentinel := service.New(messageValidator, memPinner, instantCommiter, headManager, ipnsManager)
	ratelimiter, err := setupRateLimiter(config.GRPC.RPS)
	if err != nil {
		log.Fatal().Err(err).Msg("can't setup rate limiter")
	}

	grpcServer, err := grpc.NewServer(config.GRPC, ratelimiter)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize grpc server")
	}
	sentinelpb.RegisterSentinelServer(grpcServer, sentinel)
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

func setupIPNSManager(config Config, shell *shell.Shell) (ipns.Manager, error) {
	return ipns.NewIPNSManager(shell)
}
