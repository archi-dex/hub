package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/archi-dex/ingester/pkg/db"
	"github.com/archi-dex/ingester/pkg/util"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type inputPayload struct {
	Filepath string                 `json:"filepath"`
	Data     map[string]interface{} `json:"data"`
}

func ingest(ctx context.Context, logger util.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Infow("serving ingest ws", "origin", r.RemoteAddr)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Errorw("failed to upgrade connection", "err", err)
			return
		}

		count := 0
		seconds := util.Timer()
		defer func() {
			logger.Infow("closing ingest ws", "count", count, "duration", seconds())
			if err := conn.Close(); err != nil {
				logger.Errorw("failed to close ws", "err", err)
			}
		}()

		for {
			_, buff, err := conn.ReadMessage()
			if _, ok := err.(*websocket.CloseError); ok {
				break
			}

			if err != nil {
				logger.Errorw("error reading message", "err", err)

				resp := map[string]string{"message": "error reading input: " + err.Error()}
				if err := conn.WriteJSON(resp); err != nil {
					logger.Errorw("error sending responding", "err", err)
				}

				break
			}

			data := string(buff)
			count += 1

			var input inputPayload
			if err = json.Unmarshal(buff, &input); err != nil {
				logger.Errorw("failed to parse message", "data", data)
				continue
			}

			logger.Infow("parsed data", "parsed", input)

			if _, err := db.CreateEntity(ctx, input.Filepath, input.Data); err != nil {
				logger.Errorw("failed to create entity", "input", input, "err", err)
			}
		}
	}
}
