package main

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/jmichalak9/open-pollution/server"
)

const (
	addressEnv = "ADDRESS"
)

func main() {
	address := mustGetEnv(addressEnv)
	srv := server.NewServer(address)
	err := srv.Run()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server failed")
	}
}

func mustGetEnv(envName string) string {
	envVal := os.Getenv(envName)
	if len(envVal) == 0 {
		log.Fatal().Msgf("env %s not set", envName)
	}
	return envVal
}
