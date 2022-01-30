package main

import (
	"context"
	"math/rand"
	"net"
	"time"

	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/storage/localfs"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jmichalak9/open-pollution/cmd/pdcl"
	"github.com/jmichalak9/open-pollution/cmd/pdcl/pb"
)

type Config struct {
	Host string `envconfig:"SENTINEL_SERVICE_HOST" default:"localhost"`
	Port string `envconfig:"SENTINEL_SERVICE_PORT" default:"8000"`
	pdcl.LocalStorageConfig
	GRPCConfig
}

type GRPCConfig struct {
	MaxRetries     uint `envconfig:"GRPC_MAX_RETRIES" default:"10"`
	BackoffSeconds int  `envconfig:"GRPC_BACKOFF_SEC" default:"5"`
	Exponential    bool `envconfig:"GRPC_BACKOFF_EXPONENTIAL" default:"true"`
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	contentStorage, err := localfs.NewStorage(config.Directory)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize contentStorage")
	}
	messageStorage := storage.NewProtoMessageStorage(contentStorage)

	var backoffFunc grpc_retry.BackoffFunc
	if config.GRPCConfig.Exponential {
		backoffFunc = grpc_retry.BackoffExponential(
			time.Duration(config.GRPCConfig.BackoffSeconds) * time.Second)
	} else {
		backoffFunc = grpc_retry.BackoffLinear(
			time.Duration(config.GRPCConfig.BackoffSeconds) * time.Second)
	}
	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(config.GRPCConfig.MaxRetries),
		grpc_retry.WithBackoff(backoffFunc),
		grpc_retry.WithPerRetryTimeout(20 * time.Second),
		grpc_retry.WithCodes(codes.DeadlineExceeded, codes.ResourceExhausted, codes.Unavailable),
	}
	conn, err := grpc.Dial(
		net.JoinHostPort(config.Host, config.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	messageProducer := producer.NewMessageProducer(messageStorage, sentinelClient)
	r := randomOPMessageProducer{producer: messageProducer}
	r.run()
}

type randomOPMessageProducer struct {
	producer producer.Producer
}

const (
	warsawLat   = 52
	warsawLong  = 21
	maxLevel    = 100
	levelChance = 0.7
	tempMin     = -20
	tempMax     = 30
)

func (r *randomOPMessageProducer) run() {
	for {
		o3 := rand.Int63n(maxLevel)
		temp := rand.Int63n(tempMax-tempMin) - tempMin
		time.Sleep(1 * time.Second)
		message := &pb.Message{
			MeasureTime: timestamppb.Now(),
			Location: &pb.Location{
				Latitude:   rand.NormFloat64() + warsawLat,
				Longtitude: rand.NormFloat64() + warsawLong,
			},
			O3Level:     &o3,
			Temperature: &temp,
		}
		if rand.Float64() < levelChance {
			so2 := rand.Int63n(maxLevel)
			message.SO2Level = &so2
		}
		if rand.Float64() < levelChance {
			pm10 := rand.Int63n(maxLevel)
			message.PM10Level = &pm10
		}
		if rand.Float64() < levelChance {
			pm25 := rand.Int63n(maxLevel)
			message.PM25Level = &pm25
		}
		log.Info().Time("measure_time", message.MeasureTime.AsTime()).Msg("produced message")
		// Request context deadline is managed by gRPC client.
		ctx := context.Background()
		if err := r.producer.Produce(ctx, message); err != nil {
			log.Fatal().Err(err).Msg("error producing message")
		}
	}
}
