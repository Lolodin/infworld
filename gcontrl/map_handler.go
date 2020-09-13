package gcontrl

import (
	"encoding/json"
	"fmt"
	"github.com/lolodin/infworld/wmap"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)
type requestMap struct {
	X        int
	Y        int
	PlayerID string
}

func Map_Handler(W *wmap.WorldMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("MAP HANDLER")
		body, _ := ioutil.ReadAll(r.Body)

		rm := requestMap{}

		err := json.Unmarshal(body, &rm)
		if err != nil {
			log.WithFields(log.Fields{
				"package": "GameController",
				"func":    "InitHandler",
				"error":   err,
				"data":    body,
			}).Error("Error Marshal data")
		}
		fmt.Println(rm.X, rm.Y)

		c := wmap.GetChunkID(rm.X, rm.Y)
		d := wmap.GetCurrentPlayerMap(c)
		x := wmap.GetPlayerDrawChunkMap(d, W)
		playerMap := wmap.MapToJSON(x, rm.PlayerID)
		w.Header().Set("Content-Type", "application/json")
		w.Write(playerMap)

	}

}



