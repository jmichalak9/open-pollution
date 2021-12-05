package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/areknoster/public-distributed-commit-log/consumer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/storage/localfs"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
	"github.com/areknoster/public-distributed-commit-log/thead/sentinel_reader"
	"github.com/ipfs/go-cid"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/jmichalak9/open-pollution/server"
	"github.com/jmichalak9/open-pollution/server/measurement"
)

type Config struct {
	Address  string `envconfig:"ADDRESS" required:"true"`
	PDCLHost string `envconfig:"PDCL_HOST" required:"true"`
	PDCLPort string `envconfig:"PDCL_PORT" required:"true"`
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}
	measurementCache := measurement.NewInMemoryCache(measurement.ExampleMeasurements)
	srv := server.NewServer(config.Address, measurementCache)
	go setupPDCL(measurementCache, config)
	err := srv.Run()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func setupPDCL(cache measurement.Cache, config Config) {
	conn, err := grpc.Dial(
		net.JoinHostPort(config.PDCLHost, config.PDCLPort),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	sentinelHeadReader := sentinel_reader.NewSentinelHeadReader(sentinelClient)
	consumerOffsetManager := memory.NewHeadManager(cid.Undef)
	// TODO: make this configurable
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("reading user home directory")
	}
	fsStorage, err := localfs.NewStorage(dirname + "/.local/share/pdcl/storage")
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize storage")
	}
	messageStorage := storage.NewProtoMessageStorage(fsStorage)

	firstToLastConsumer := consumer.NewFirstToLastConsumer(
		sentinelHeadReader,
		consumerOffsetManager,
		messageStorage,
		consumer.FirstToLastConsumerConfig{
			// TODO: these should be configurable
			PollInterval: 10 * time.Second,
			PollTimeout:  100 * time.Second,
		})

	c := make(chan os.Signal, 1)
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signal.Notify(c, os.Interrupt)
	// TODO: graceful shutdown
	defer signal.Stop(c)
	go func() {
		for range c {
			cancel()
		}
	}()
	err = firstToLastConsumer.Consume(globalCtx, consumer.MessageFandlerFunc(
		func(ctx context.Context, unmarshallable storage.ProtoUnmarshallable) error {
			message := &pb.Message{}
			if err := unmarshallable.Unmarshall(message); err != nil {
				return fmt.Errorf("unmarshall message: %w", err)
			}
			type pdclMeasurement struct {
				PollutionLevel float32 `json:"pollutionLevel"`
			}
			mes := measurement.Measurement{
				O3: int(message.PollutionLevel),
			}
			cache.AppendMeasurements([]measurement.Measurement{mes})
			log.Info().Msgf("received %+v", mes)
			return nil
		}))
	if err != nil {
		log.Fatal().Err(err).Msg("consuming messages")
	}
}
