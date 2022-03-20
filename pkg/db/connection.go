package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/archi-dex/ingester/pkg/util"
	trace "github.com/hans-m-song/go-stacktrace"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client            *mongo.Client
	ErrorConnecting   = trace.New("ERROR_CONNECTING")
	ErrorNotConnected = trace.New("ERROR_NOT_CONNECTED")
	ErrorPinging      = trace.New("ERROR_PINGING")
)

func NewConnection(ctx context.Context, logger util.Logger) (*mongo.Client, error) {
	if client != nil {
		return client, nil
	}

	opts := util.GetOptions()
	server := fmt.Sprintf("%s:%s", opts.DbHost, opts.DbPort)
	obscured := strings.Map(func(r rune) rune { return 'x' }, opts.DbPass)
	uri := fmt.Sprintf("mongodb://%s:%s@%s", opts.DbUser, opts.DbPass, server)
	obscuredUri := fmt.Sprintf("mongodb://%s:%s@%s", opts.DbUser, obscured, server)

	logger.Debugw("connecting to db", "uri", obscuredUri)

	newClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, ErrorConnecting.Trace(err).Add("server", obscuredUri)
	}

	client = newClient
	return client, nil
}

func Healthcheck(ctx context.Context) error {
	if client == nil {
		return ErrorNotConnected.Tracef("not currently connected to a database")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return ErrorPinging.Trace(err)
	}

	return nil
}
