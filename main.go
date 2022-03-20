package main

import (
	"context"
	"net/http"

	"github.com/archi-dex/ingester/pkg/api"
	"github.com/archi-dex/ingester/pkg/db"
	"github.com/archi-dex/ingester/pkg/util"
	"github.com/spf13/cobra"
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
	ctx := context.TODO()

	client, err := db.NewConnection(ctx, logger)
	if err != nil {
		logger.Fatalw("error connecting", "err", err)
	}

	defer client.Disconnect(ctx)

	if err := db.Healthcheck(ctx); err != nil {
		logger.Fatalw("error pinging database", "err", err)
	}

	r := api.NewRouter(ctx, logger)
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Fatalw("error during server runtime", "err", err)
	}
}
