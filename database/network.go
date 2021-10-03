package database

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type NetworkModel struct {
	DemandTPS     int64
	ThroughputTPS int64
	LastUpdated   time.Time
}

func PublishNetwork(network NetworkModel) error {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./service-account.json")
	client, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"), sa)
	if err != nil {
		return err
	}

	entry := client.Doc("Network/Activity")
	_, err = entry.Set(ctx, network)
	if err != nil {
		return err
	}

	return nil
}
