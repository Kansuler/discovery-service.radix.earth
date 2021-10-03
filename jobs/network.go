package jobs

import (
	"time"

	"github.com/Kansuler/radix-discovery-service/api"
	"github.com/Kansuler/radix-discovery-service/database"
	"github.com/rs/zerolog/log"
)

func Network() {
	log.Debug().Msg("run network job")

	demandResp, err := api.GetDemand()
	if err != nil {
		log.Error().Err(err).Msg("could not get demand")
		return
	}

	throughputResp, err := api.GetThroughput()
	if err != nil {
		log.Error().Err(err).Msg("could not get throughput")
		return
	}

	err = database.PublishNetwork(database.NetworkModel{
		ThroughputTPS: throughputResp.Result.TPS,
		DemandTPS:     demandResp.Result.TPS,
		LastUpdated:   time.Now(),
	})
	if err != nil {
		log.Error().Err(err).Msg("could not save network activity to database")
		return
	}

	log.Info().Msg("saved network activity to database")
}
