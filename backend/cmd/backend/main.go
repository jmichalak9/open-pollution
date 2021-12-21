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
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/storage/localfs"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
	sentinelhead "github.com/areknoster/public-distributed-commit-log/thead/sentinel"
	"github.com/ipfs/go-cid"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/jmichalak9/open-pollution/cmd/pdcl"
	"github.com/jmichalak9/open-pollution/cmd/pdcl/pb"
	"github.com/jmichalak9/open-pollution/server"
	"github.com/jmichalak9/open-pollution/server/measurement"
)

type Config struct {
	Address  string `envconfig:"ADDRESS" required:"true"`
	PDCLHost string `envconfig:"PDCL_HOST" required:"true"`
	PDCLPort string `envconfig:"PDCL_PORT" required:"true"`
	pdcl.LocalStorageConfig
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}
	measurementCache := measurement.NewInMemoryCache(measurement.ExampleMeasurements)
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
	conn, err := grpc.Dial(
		net.JoinHostPort(config.PDCLHost, config.PDCLPort),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	sentinelHeadReader := sentinelhead.New(sentinelClient)
	consumerOffsetManager := memory.NewHeadManager(cid.Undef)
	fsStorage, err := localfs.NewStorage(config.Directory)
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
