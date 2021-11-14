package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ValidatorModel struct {
	PublicKey               string    `json:"publicKey"`
	TotalDeligatedStake     string    `json:"totalDelegatedStake"`
	UptimePercentage        string    `json:"uptimePercentage"`
	ProposalsMissed         int64     `json:"proposalsMissed"`
	Address                 string    `json:"address"`
	InfoURL                 string    `json:"infoURL"`
	OwnerDelegation         string    `json:"ownerDeligation"`
	Name                    string    `json:"name"`
	ValidatorFee            string    `json:"validatorFee"`
	Registered              bool      `json:"registered"`
	OwnerAddress            string    `json:"ownerAddress"`
	IsExternalStakeAccepted bool      `json:"isExternalStakeAccepted"`
	ProposalsCompleted      int64     `json:"proposalsCompleted"`
	IP                      string    `json:"ip"`
	NodeAddress             string    `json:"nodeAddress"`
	NodeMatchFound          bool      `json:"nodeMatchFound"`
	Latitude                float64   `json:"latitude"`
	Longitude               float64   `json:"longitude"`
	Country                 string    `json:"country"`
	City                    string    `json:"city"`
	Region                  string    `json:"region"`
	ISP                     string    `json:"isp"`
	Organisation            string    `json:"org"`
	LastUpdated             time.Time `json:"lastUpdated"`
	ManualLocationData      bool      `json:"manualLocationData"`
}

func PublishValidators(validators []ValidatorModel) error {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./service-account.json")
	client, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"), sa)
	if err != nil {
		return err
	}

	defer client.Close()

	for _, validator := range validators {
		entry := client.Doc(fmt.Sprintf("Validators/%s", validator.Address))

		validatorSnapshot, err := entry.Get(ctx)
		if status.Code(err) == codes.NotFound {
			entry.Create(ctx, validator)
			continue
		}

		var validatorData ValidatorModel

		err = validatorSnapshot.DataTo(&validatorData)
		if err != nil {
			log.Error().Err(err).Msgf("could not scan validator data to model")
			continue
		}

		if validatorData.ManualLocationData {
			_, err = entry.Set(ctx, ValidatorModel{
				PublicKey:               validator.PublicKey,
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
				IP:                      validatorData.IP,
				NodeAddress:             validator.NodeAddress,
				NodeMatchFound:          validator.NodeMatchFound,
				Latitude:                validatorData.Latitude,
				Longitude:               validatorData.Longitude,
				Country:                 validatorData.Country,
				City:                    validatorData.City,
				Region:                  validatorData.Region,
				ISP:                     validatorData.ISP,
				Organisation:            validatorData.Organisation,
				LastUpdated:             validator.LastUpdated,
				ManualLocationData:      validatorData.ManualLocationData,
			})

			if err != nil {
				log.Error().Err(err).Msgf("could not update validator with partial data")
			}

			continue
		}

		_, err = entry.Set(ctx, validator)
		if err != nil {
			log.Error().Err(err).Msgf("could not save validator %s to database", validator.Name)
		}
	}

	return nil
}
