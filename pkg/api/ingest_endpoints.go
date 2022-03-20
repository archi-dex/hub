package api

import (
	"context"
	"net/http"

	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func ingest(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infow("serving ingest ws")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorw("failed to upgrade connection", "err", err)
			return
		}

		count := 0
		seconds := util.Timer()
		var msg interface{}
		defer conn.Close()
		for {
			if err := conn.ReadJSON(&msg); err != nil {
				if _, ok := err.(*websocket.CloseError); !ok {
					logger.Errorw("error reading message", "err", err)
				}
				break
			}

			count += 1
			logger.Infow("received message", "msg", msg)
			// TODO switch mt, create entity
		}

		logger.Infow("closed ingest ws", "count", count, "duration", seconds())
	}
}
