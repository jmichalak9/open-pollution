package main

import (
	"context"
	"crypto"
	"fmt"
	"math/rand"
	"net"
	"time"

	pdclcrypto "github.com/areknoster/public-distributed-commit-log/crypto"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jmichalak9/open-pollution/cmd/pdcl/pb"
)

type Config struct {
	Host                     string `envconfig:"SENTINEL_HOST" required:"true"`
	Port                     string `envconfig:"SENTINEL_PORT" required:"true"`
	IPFSHost                 string `envconfig:"IPFS_HOST" default:"localhost"`
	IPFSPort                 string `envconfig:"IPFS_PORT" default:"5001"`
	ConcurrentProducerConfig producer.BasicConcurrentProducerConfig
	SignerID                 string `envconfig:"SIGNER_ID" required:"true"`
	PrivKeyPath              string `envconfig:"PRODUCER_KEY_PATH" required:"true"`
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

	codec := pbcodec.Json{}

	shell := shell.NewShell(fmt.Sprintf("%s:%s", config.IPFSHost, config.IPFSPort))
	writer := ipfsstorage.NewStorage(shell, codec)

	privKey, err := pdclcrypto.LoadFromPKCSFromPEMFile(config.PrivKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("get privKey")
	}

	signer, ok := privKey.(crypto.Signer)
	if !ok {
		log.Fatal().Msgf("key is not private crypto.Signer type but %T", privKey)
	}

	signedWriter := pdclcrypto.NewSignedMessageWriter(writer, codec, config.SignerID, signer)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	blockingProducer := producer.NewBlockingProducer(signedWriter, sentinelClient)
	concurrentProducer := producer.StartBasicConcurrentProducer(ctx, blockingProducer, config.ConcurrentProducerConfig)

	go queueMessages(ctx, concurrentProducer.Messages())
	handleErrors(concurrentProducer.Errors())
}

const (
	warsawLat   = 52
	warsawLong  = 21
	maxLevel    = 100
	levelChance = 0.7
	tempMin     = -20
	tempMax     = 30
)

func handleErrors(errors <-chan producer.Error) {
	for err := range errors {
		log.Error().
			RawJSON("message", []byte(protojson.Format(err.Message))).
			Err(err.Err).
			Msg("message production failed")
	}
}

func queueMessages(ctx context.Context, messages chan<- proto.Message) {
	defer close(messages)
	for i := int64(0); i < 10; i++ {
		select {
		case <-ctx.Done():
			log.Error().Int64("message_index", i).Msg("context done, production stopped before all messages were sent")
			return
		default:
			o3 := rand.Int63n(maxLevel)
			temp := rand.Int63n(tempMax-tempMin) - tempMin
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
			messages <- message
		}
	}
	log.Info().Msg("messages queued")
}
