package main

import (
	"api.cloud.io/config"
	"api.cloud.io/internal/api"
	"api.cloud.io/internal/db/pg"
	"api.cloud.io/internal/util/jobs"
	"context"
	"fmt"
	"github.com/pascaldekloe/metrics/gostat"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)
var signals chan os.Signal

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	log.Info().Msgf("Daemon launch as %s", strings.Join(os.Args, " "))
	signals = make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// include Go runtime metrics
	gostat.CaptureEvery(5 * time.Second)
	_, err := pg.NewClient(config.PgConfig{
		Host:           "127.0.0.1",
		Port:           5432,
		UserName:       "cloudapi",
		Password:       "password",
		Database:       "cloudapi",
		Sslmode:        "disable",
		MaxOpenConns:   10,
		MigrationsDir:  "file://./db/migration/",
		MigrateVersion: 1,
		JsonDir: "./json/",
	})
	if err != nil {
		log.Fatal().Str("Failed to start pg:",err.Error())
		return
	}
	var c config.Config
	mainContext, mainCancel := context.WithCancel(context.Background())
	httpServerJob := startHTTPServer(mainContext, &c)
	signal := <-signals
	timeout := c.ShutdownTimeout.WithDefault(5 * time.Second)
	log.Info().Msgf("Shutting down services initiated with timeout in %s", timeout)
	mainCancel()
	finishCTX, finishCancel := context.WithTimeout(context.Background(), timeout)
	defer finishCancel()
	jobs.WaitAll(finishCTX,
		httpServerJob,
	)
	log.Fatal().Msgf("Exit on signal %s", signal)
}
func startHTTPServer(ctx context.Context, c *config.Config) *jobs.Job {
	if c.ListenPort == 0 {
		c.ListenPort = 8080
		log.Info().Msgf("Default HTTP server listen port to %d", c.ListenPort)
	}
	api.InitHandler()
	srv := &http.Server{
		Handler:      cors.AllowAll().Handler(api.Handler),
		Addr:         fmt.Sprintf(":%d", c.ListenPort),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// launch HTTP server
	go func() {
		err := srv.ListenAndServe()
		log.Error().Err(err).Msg("HTTP stopped")
		signals <- syscall.SIGABRT
	}()

	ret := jobs.Start("HTTPserver", func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("HTTP failed shutdown")
		}
	})
	return &ret

}


