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

	"github.com/Thibaut-gauvin/kie/internal/k8s"
	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/go-co-op/gocron/v2"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
)

const (
	// Go cronjob task name.
	clusterImageScrapJobName = "images metrics scrapper"
)

// ServeOpts hold "kie serve" command options.
type ServeOpts struct {
	ListenPort      string `mapstructure:"listen_port"`
	Kubeconfig      string `mapstructure:"kubeconfig"`
	RefreshInterval int    `mapstructure:"refresh_interval"`
}

// StartKie handle main logic of starting kie server.
func StartKie(opts ServeOpts) error {
	k8sClient, err := k8s.NewClient(opts.Kubeconfig)
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	// Run our server in a goroutine so that it doesn't block.
	srv := getHTTPServer(opts.ListenPort)
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

	err = startCronScheduler(opts.RefreshInterval, *k8sClient)
	if err != nil {
		return err
	}

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

// getHTTPServer initialize a new router and declare our app routes then return a http.Server instance.
func getHTTPServer(listenPort string) *http.Server {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	router.HandleFunc("/readiness", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	router.Handle("/metrics", promhttp.Handler())

	timeoutDuration := 15 * time.Second
	return &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%s", listenPort),
		Handler:           handlers.CombinedLoggingHandler(os.Stdout, router), // Nginx access logs style
		ReadHeaderTimeout: timeoutDuration,
		ReadTimeout:       timeoutDuration,
		WriteTimeout:      timeoutDuration,
		IdleTimeout:       timeoutDuration,
	}
}

// startCronScheduler create & start a new cronjob scheduler instance in a non-blocking way.
// It also registers our "images metrics scrapper" task.
func startCronScheduler(refreshInterval int, client kubernetes.Clientset) error {
	scheduler, err := gocron.NewScheduler(
		gocron.WithGlobalJobOptions(
			gocron.WithStartAt(gocron.WithStartImmediately()),
			gocron.WithSingletonMode(gocron.LimitModeReschedule),
		),
	)
	if err != nil {
		return fmt.Errorf("unable to create cron scheduler, %w", err)
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(time.Duration(refreshInterval)*time.Second),
		gocron.NewTask(scrapClusterImages, client),
		gocron.WithName(clusterImageScrapJobName),
	)
	if err != nil {
		return fmt.Errorf("unable to create job %s, %w", clusterImageScrapJobName, err)
	}

	scheduler.Start()
	return nil
}
