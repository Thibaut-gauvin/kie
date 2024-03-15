package kie

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ServeOpts struct {
	ListenPort string `mapstructure:"listen_port"`
}

func StartServer(opts ServeOpts) error {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("/readiness", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle("/metrics", promhttp.Handler())

	timeoutDuration := 15 * time.Second
	srv := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%s", opts.ListenPort),
		Handler:           handlers.CombinedLoggingHandler(os.Stdout, router), // Nginx access logs style
		ReadHeaderTimeout: timeoutDuration,
		ReadTimeout:       timeoutDuration,
		WriteTimeout:      timeoutDuration,
		IdleTimeout:       timeoutDuration,
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// Run our server in a goroutine so that it doesn't block.
	logger.Infof("Sarting http server on port %s", opts.ListenPort)
	go func() {
		defer waitGroup.Done()
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				logger.Infof("Server closed")
			} else {
				logger.Fatalf("Failed to start server %v", err)
			}
		}
	}()

	// Graceful shutdown, inspired by https://github.com/gorilla/mux#graceful-shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	<-c

	time.Sleep(15 * time.Second) // We give time to the readiness probe to be down
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	waitGroup.Wait()
	logger.Infof("Shutting down...")

	return nil
}
