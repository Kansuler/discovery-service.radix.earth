package jobs

import (
	"fmt"
	"time"

	"github.com/Kansuler/radix-discovery-service/api"
	"github.com/Kansuler/radix-discovery-service/database"
	"github.com/Kansuler/radix-discovery-service/lookup"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/rs/zerolog/log"
)

func Nodes() {
	log.Debug().Msg("run nodes job")

	peers, err := api.GetPeers()
	if err != nil {
		log.Error().Err(err).Msg("could not get peers")
		return
	}

	log.Debug().Int("peers length", len(peers.Result)).Send()

	validators, err := api.GetValidatorSet()
	if err != nil {
		log.Error().Err(err).Msg("could not get validators")
		return
	}

	var validatorModels []database.ValidatorModel

	log.Debug().Int("validators length", len(validators.Result.Validators)).Send()

	for _, validator := range validators.Result.Validators {
		_, decodedBytes, err := bech32.Decode(validator.Address)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}

		convertedBytes, err := bech32.ConvertBits(decodedBytes, 5, 8, false)

		validatorModels = append(validatorModels, database.ValidatorModel{
			PublicKey:               fmt.Sprintf("%x", convertedBytes),
			TotalDeligatedStake:     validator.TotalDeligatedStake,
			UptimePercentage:        validator.UptimePercentage,
			ProposalsMissed:         validator.ProposalsMissed,
			Address:                 validator.Address,
			InfoURL:                 validator.InfoURL,
			OwnerDelegation:         validator.OwnerDelegation,
			Name:                    validator.Name,
			ValidatorFee:            validator.ValidatorFee,
			Registered:              validator.Registered,
			OwnerAddress:            validator.OwnerAddress,
			IsExternalStakeAccepted: validator.IsExternalStakeAccepted,
			ProposalsCompleted:      validator.ProposalsCompleted,
			LastUpdated:             time.Now(),
			DisplayNode:             validator.UptimePercentage != "0.00" && validator.UptimePercentage != "0",
		})
	}

	var matches int

	for _, peer := range peers.Result {
		if len(peer.Channels) == 0 {
			log.Debug().Str("peerAddress", peer.Address).Msg("no channels exist for peer")
			continue
		}

		_, decodedBytes, err := bech32.Decode(peer.Address)
		if err != nil {
			log.Error().Err(err).Send()
			continue
		}

		convertedBytes, err := bech32.ConvertBits(decodedBytes, 5, 8, false)
		pubKey := fmt.Sprintf("%x", convertedBytes)

		var found bool
		for index, validator := range validatorModels {
			if validator.PublicKey == pubKey {
				validatorModels[index].IP = peer.Channels[0].IP
				validatorModels[index].NodeAddress = peer.Address
				validatorModels[index].NodeMatchFound = true
				matches++
				break
			}
		}

		if !found {
			log.Debug().Str("pubKey", pubKey).Msg("no match found...")
		}
	}

	log.Debug().Int("matched peers", matches).Int("unmatched peers", len(peers.Result)-matches).Send()

	var batches [][]string

	for _, validator := range validatorModels {
		if !validator.NodeMatchFound {
			log.Debug().Str("name", validator.Name).Msg("no peer match")
			continue
		}

		index := 0
		var assigned bool
		for !assigned {
			if len(batches) < index+1 {
				batches = append(batches, []string{})
			}

			if len(batches[index]) == 100 {
				index++
				continue
			}

			batches[index] = append(batches[index], validator.IP)
			assigned = true
		}
	}

	for _, batch := range batches {
		locations, err := lookup.IP(batch)
		if err != nil {
			log.Error().Err(err).Msg("could not get information about the ip numbers")
			return
		}

		for _, location := range *locations {
			for index, validator := range validatorModels {
				if validator.IP == location.Query {
					validatorModels[index].Latitude = location.Latitude
					validatorModels[index].Longitude = location.Longitude
					validatorModels[index].Country = location.Country
					validatorModels[index].City = location.City
					validatorModels[index].Region = location.Region
					validatorModels[index].ISP = location.ISP
					validatorModels[index].Organisation = location.Organisation

					break
				}
			}
		}
	}

	err = database.PublishValidators(validatorModels)
	if err != nil {
		log.Error().Err(err).Msg("could not save dataset to database")
		return
	}

	log.Info().Msg("successfully saved all nodes")
}
