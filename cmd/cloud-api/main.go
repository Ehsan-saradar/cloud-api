package main

import (
	"api.cloud.io/internal/db/pg"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"api.cloud.io/config"
	"api.cloud.io/internal/api"
	"api.cloud.io/internal/util/jobs"

	"github.com/pascaldekloe/metrics/gostat"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var signals chan os.Signal

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	log.Info().Msgf("Daemon launch as %s", strings.Join(os.Args, " "))
	sec:=rand.Intn(10) + 2
	log.Info().Msgf("wait for %d sec and then try to connect to pg",sec)
	time.Sleep(time.Duration(sec)*time.Second)
	signals = make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	// include Go runtime metrics
	gostat.CaptureEvery(5 * time.Second)
	for i:=0;i<1000;i++{
		_, err := pg.NewClient(config.PgConfig{
			Host:           "pg",
			Port:           5432,
			UserName:       "cloudapi",
			Password:       "password",
			Database:       "cloudapi",
			Sslmode:        "disable",
			MaxOpenConns:   10,
			MigrationsDir:  "file://./db/migration/",
			MigrateVersion: 1,
			JsonDir:        "./json/",
		})
		if err != nil {
			log.Fatal().Str("Failed to start pg:", err.Error())
			if i<999{
				log.Info().Msgf("Failed to connect to pg, wait for 10 sec and try again (err:%s)",err.Error())
				time.Sleep(10*time.Second)
			}else{
				log.Info().Msgf("Failed to connect to pg, exit (err:%s)",err.Error())
				return
			}

		}else{
			break
		}
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
	api.InitHandler(8080, 8090)
	mainSrv := &http.Server{
		Handler:      cors.AllowAll().Handler(api.Handler),
		Addr:         fmt.Sprintf(":%d", c.ListenPort),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// launch HTTP server
	go func() {
		err := mainSrv.ListenAndServe()
		log.Error().Err(err).Msg("HTTP stopped")
		signals <- syscall.SIGABRT
	}()

	analysisSrv := &http.Server{
		Handler:      cors.AllowAll().Handler(api.Handler),
		Addr:         fmt.Sprintf(":%d", 8090),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	// launch HTTP server
	go func() {
		err := analysisSrv.ListenAndServe()
		log.Error().Err(err).Msg("HTTP stopped")
		signals <- syscall.SIGABRT
	}()

	ret := jobs.Start("HTTPserver", func() {
		<-ctx.Done()
		if err := mainSrv.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("HTTP failed shutdown")
		}
		if err := analysisSrv.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("HTTP failed shutdown")
		}
	})
	return &ret
}
