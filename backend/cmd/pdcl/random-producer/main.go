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
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jmichalak9/open-pollution/cmd/pdcl"
	"github.com/jmichalak9/open-pollution/cmd/pdcl/pb"
)

type Config struct {
	Host string `envconfig:"SENTINEL_SERVICE_HOST" default:"localhost"`
	Port string `envconfig:"SENTINEL_SERVICE_PORT" default:"8000"`
	pdcl.LocalStorageConfig
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

	conn, err := grpc.Dial(
		net.JoinHostPort(config.Host, config.Port),
		grpc.WithInsecure(),
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

func (r *randomOPMessageProducer) run() {
	for {
		o3 := rand.Int63n(100)
		temp := rand.Int63n(50) - 20
		time.Sleep(1 * time.Second)
		message := &pb.Message{
			MeasureTime: timestamppb.Now(),
			Location: &pb.Location{
				Latitude:   rand.NormFloat64() + 52,
				Longtitude: rand.NormFloat64() + 21,
			},
			O3Level:     &o3,
			Temperature: &temp,
		}
		if rand.Float64() > 0.7 {
			so2 := rand.Int63n(100)
			message.SO2Level = &so2
		}
		if rand.Float64() > 0.7 {
			pm10 := rand.Int63n(100)
			message.PM10Level = &pm10
		}
		if rand.Float64() > 0.7 {
			pm25 := rand.Int63n(100)
			message.PM25Level = &pm25
		}
		log.Info().Time("measure_time", message.MeasureTime.AsTime()).Msg("produced message")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := r.producer.Produce(ctx, message); err != nil {
			log.Fatal().Err(err).Msg("error producing message")
		}
	}
}