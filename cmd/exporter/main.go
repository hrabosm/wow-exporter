package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/hrabosm/wow-exporter/internal/collector"
	"github.com/hrabosm/wow-exporter/internal/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	log.Printf("Loaded characters: %+v", cfg.Characters)

	wowCollector := collector.NewWowCollector(cfg.Characters)

	prometheus.MustRegister(wowCollector)

	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr: ":9105",
	}

	go func() {
		fmt.Println("Exporter running on :9105")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	fmt.Println("Shutting down exporter...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Shutdown error:", err)
	}

	fmt.Println("Exporter stopped")
}
