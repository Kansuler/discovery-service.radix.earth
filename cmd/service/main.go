package main

import (
	"os"
	"time"

	"github.com/Kansuler/radix-discovery-service/jobs"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("could not load .env file")
	}

	gocron.Every(5).Seconds().Do(jobs.Network)
	gocron.Every(15).Minutes().Do(jobs.Nodes)
	jobs.Nodes()

	<-gocron.Start()

	// This point is only reached if gocron stops
	log.Error().Msg("cron stopped unexpectedly")
}
