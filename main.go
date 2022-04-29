package main

import (
	"context"
	"net/http"

	"github.com/archi-dex/ingester/pkg/api"
	"github.com/archi-dex/ingester/pkg/db"
	"github.com/archi-dex/ingester/pkg/util"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var cmd = &cobra.Command{
	Use:   "ingester",
	Short: "ingests and manages storing of data",
	Run:   run,
}

func init() {
	util.InitFlags(cmd)
}

func main() {
	cmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	opts := util.InitOptions(cmd)
	logger := util.InitLogger(opts.Debug)
	ctx := context.Background()

	client, err := db.NewConnection(ctx, logger)
	if err != nil {
		logger.Fatalw("error connecting", "err", err)
	}

	defer client.Disconnect(ctx)

	if err := db.Healthcheck(ctx); err != nil {
		logger.Fatalw("error pinging database", "err", err)
	}

	client.
		Database(opts.DbName).
		Collection(db.EntityCollectionName).
		Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{Keys: bson.M{"dir": 1}},
			{Keys: bson.M{"base": 1}},
		})

	r := api.NewRouter(ctx, logger)
	logger.Infow("server ready", "url", "http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatalw("error during server runtime", "err", err)
	}
}
