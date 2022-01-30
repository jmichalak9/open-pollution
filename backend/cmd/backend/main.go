package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/areknoster/public-distributed-commit-log/consumer"
	pdclcrypto "github.com/areknoster/public-distributed-commit-log/crypto"
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/jmichalak9/open-pollution/cmd/pdcl/pb"
	"github.com/jmichalak9/open-pollution/server"
	"github.com/jmichalak9/open-pollution/server/measurement"
)

type Config struct {
	Address  string `envconfig:"ADDRESS" required:"true"`
	PDCLHost string `envconfig:"PDCL_HOST" required:"true"`
	PDCLPort string `envconfig:"PDCL_PORT" required:"true"`
	IPFSPort string `envconfig:"IPFS_PORT" required:"true"`
	IPFSHost string `envconfig:"IPFS_HOST" required:"true"`
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
	measurementCache := measurement.NewInMemoryCache([]measurement.Measurement{})
	srv := server.NewServer(config.Address, measurementCache)
	ctx, cancelPDCL := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		setupPDCL(ctx, measurementCache, config)
	}()

	go func() {
		wg.Add(1)
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()
	waitForShutdown()
	cancelPDCL()
	wg.Done()
	if err := srv.Shutdown(); err != nil {
		log.Error().Err(err).Msg("failed to shut down")
	}
	wg.Done()
	wg.Wait()
	log.Debug().Msg("server stopped")
}

func waitForShutdown() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	log.Debug().Msg("interruption signal received")
}

func setupPDCL(ctx context.Context, cache measurement.Cache, config Config) {
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
		net.JoinHostPort(config.PDCLHost, config.PDCLPort),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to sentinel")
	}
	resp, err := sentinelClient.GetHeadIPNS(context.Background(), &sentinelpb.GetHeadIPNSRequest{})
	if err != nil {
		log.Fatal().Err(err).Msg("getting ipns address")
	}
	log.Info().Msgf("IPNS head address is %s", resp.IpnsAddr)
	if resp.IpnsAddr == "" {
		log.Fatal().Msg("could not get valid IPNS address from sentinel")
	}
	consumerOffsetManager := memory.NewHeadManager(cid.Undef)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize storage")
	}

	shell := shell.NewShell(fmt.Sprintf("%s:%s", config.IPFSHost, config.IPFSPort))
	ipnsResolver := ipns.NewIPNSResolver(shell)
	reader := ipfsstorage.NewStorage(shell, pbcodec.Json{})

	firstToLastConsumer := consumer.NewFirstToLastConsumer(
		consumerOffsetManager,
		reader,
		pdclcrypto.NewSignedMessageUnwrapper(reader, pbcodec.Json{}),
		consumer.FirstToLastConsumerConfig{
			PollInterval: 20 * time.Second,
			PollTimeout:  20 * time.Second,
			IPNSAddr:     resp.IpnsAddr,
		},
		ipnsResolver,
	)

	err = firstToLastConsumer.Consume(ctx, consumer.MessageHandlerFunc(
		func(ctx context.Context, unmarshallable storage.ProtoUnmarshallable) error {
			message := &pb.Message{}
			if err := unmarshallable.Unmarshall(message); err != nil {
				return fmt.Errorf("unmarshall message: %w", err)
			}

			mes := measurement.Measurement{
				Position: measurement.Position{
					Lat:  message.Location.Latitude,
					Long: message.Location.Longtitude,
				},
				Timestamp: message.MeasureTime.AsTime(),
			}
			if message.O3Level != nil {
				mes.O3 = int(*message.O3Level)
			}
			if message.SO2Level != nil {
				mes.SO2 = int(*message.SO2Level)
			}
			if message.PM10Level != nil {
				mes.PM10 = int(*message.PM10Level)
			}
			if message.PM25Level != nil {
				mes.PM25 = int(*message.PM25Level)
			}
			if message.Temperature != nil {
				mes.Temperature = int(*message.Temperature)
			}
			cache.UpdateMeasurements([]measurement.Measurement{mes})
			log.Info().Msgf("received %+v", mes)
			return nil
		}))
	log.Debug().Msg("server stopped consuming messages")
	if err != nil {
		log.Fatal().Err(err).Msg("consuming messages")
	}
}
