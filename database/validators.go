package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
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
}

func PublishValidators(validators []ValidatorModel) error {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./service-account.json")
	client, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"), sa)
	if err != nil {
		return err
	}

	for _, validator := range validators {
		entry := client.Doc(fmt.Sprintf("Validators/%s", validator.Address))

		var err error
		if validator.NodeMatchFound {
			_, err = entry.Set(ctx, validator)
		} else {
			entry.Create(ctx, validator)
		}

		if err != nil {
			log.Error().Err(err).Msgf("could not save validator %s to database", validator.Name)
		}
	}

	return nil
}
