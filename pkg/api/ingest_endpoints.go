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
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorw("failed to upgrade connection", "err", err)
			return
		}

		defer conn.Close()
		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Errorw("failed to read message", "err", err)
				break
			}

			logger.Infow("received message", "type", mt, "msg", msg)
			// TODO switch mt, create entity
		}
	}
}
