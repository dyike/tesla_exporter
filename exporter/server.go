package exporter

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, r *prometheus.Registry) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		r, promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
	))
	httpServer := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	s := &Server{
		httpServer: httpServer,
	}
	return s
}

// ListenAndServe method serves HTTP requests.
func (s *Server) ListenAndServe() {
	log.Println("Starting HTTP server")

	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Printf("Failed to start serving HTTP requests: %v", err)
	}
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.httpServer.Shutdown(ctx)
}
