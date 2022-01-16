package main

import (
	"context"
	"fmt"
	"github.com/zookzook/diesel-importer/pkg/config"
	"github.com/zookzook/diesel-importer/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

// fetches the raw data from our provider
func fetchData(cfg *config.Config) ([]utils.Station, error) {

	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, cfg.CoreAPIURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: got %v", res.Status)
	}

	return utils.Parse(res.Body)
}

// saves the new data to the stations collections
func saveData(client *mongo.Client, stations []utils.Station) error {

	fuelDatabase := client.Database("fuel")
	stationsCollection := fuelDatabase.Collection("stations")

	for _, station := range stations {
		_, err := stationsCollection.InsertOne(context.TODO(), station)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("can't get config: %s", err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	s := utils.NewStopWatch()

	for {
		s.Start("import data")

		stations, err := fetchData(cfg)

		if err != nil {
			log.Fatal(err)
		}

		err = saveData(client, stations)

		if err != nil {
			log.Fatal(err)
		}

		s.Stop()

		log.Println(s.String())

		time.Sleep(5 * time.Minute)
	}
}
