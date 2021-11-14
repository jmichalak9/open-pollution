package server

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/jmichalak9/open-pollution/server/measurement"
)

// Server serves open pollution API.
type Server struct {
	httpServer *http.Server
}

// NewServer returns a new instance of Server.
func NewServer(address string) *Server {
	mux := http.NewServeMux()
	s := &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: mux,
		},
	}
	mux.HandleFunc("/measurements", s.handleMeasurements())

	return s
}

func (s *Server) handleMeasurements() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// measurements := []measurement.Measurement{
		// 	{
		// 		Position: measurement.Position{
		// 			Lat:  21.37,
		// 			Long: 42.69,
		// 		},
		// 	},
		// }
		// bytes, err := json.Marshal(m)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(measurement.ExampleMeasurements)
		if err != nil {
			log.Info().Err(err).Msg("encoding response")
		}
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
