package gcontrl

import (
	"fmt"
	"github.com/lolodin/infworld/chunk"
	"github.com/lolodin/infworld/wmap"
	"golang.org/x/net/websocket"
	log "github.com/sirupsen/logrus"
)

/*
{action: Move, id = name, x:+1, y:+1}
*/
type PlayerResponse struct {
	Id string
	chunk.Coordinate
}
type checkAction struct {
	Action int `json:"action"`
}


func PlayerHandler(W *wmap.WorldMap) func(ws *websocket.Conn) {
	return func(ws *websocket.Conn) {

		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"package": "GameController",
					"func":    "PlayerHandler",
					"error":   err,
				}).Error("Error ws")
			}

		}()

		player := PlayerResponse{}
		websocket.JSON.Receive(ws, &player)
		//Game Loop
		log.WithFields(log.Fields{
			"package": "GameController",
			"func":    "PlayerHandler",
			"player":  player,
		}).Info("Connect player")
fmt.Println(player)

		//ws handler
		for {



		}

	}
}