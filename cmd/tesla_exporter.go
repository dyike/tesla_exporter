package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tesla_exporter/exporter"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var s *exporter.Server

func main() {
	var email = flag.String("email", "", "tesla email address.")
	var password = flag.String("password", "", "tesla account password.")
	var internal = flag.Duration("expire", 30*time.Second, "expire cache metrics.")

	flag.Parse()
	// init collector
	collector := exporter.NewCollector(*email, *password, *internal)
	go collector.Refresh()
	r := prometheus.NewRegistry()
	if err := r.Register(collector); err != nil {
		log.Fatal("Register collector failed with %w", err)
	}
	s = exporter.NewServer("9610", r)
	go s.ListenAndServe()

	// handle exit signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	s.Stop()
}
