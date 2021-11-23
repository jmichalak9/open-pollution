package server

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/jmichalak9/open-pollution/server/measurement"
)

// Server serves open pollution API.
type Server struct {
	httpServer       *http.Server
	measurementCache measurement.Cache
}

// NewServer returns a new instance of Server.
func NewServer(address string, measurementCache measurement.Cache) *Server {
	mux := http.NewServeMux()
	s := &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: mux,
		},
		measurementCache: measurementCache,
	}
	mux.HandleFunc("/measurements", s.handleMeasurements())

	return s
}

func (s *Server) handleMeasurements() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := json.NewEncoder(w).Encode(s.measurementCache.GetMeasurements())
		if err != nil {
			log.Info().Err(err).Msg("encoding response")
		}
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
