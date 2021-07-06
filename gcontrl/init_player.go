package gcontrl

import (
	"encoding/json"
	"github.com/lolodin/infworld/wmap"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type requestPlayer struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type responsePlayer struct {
	Error string `json:"error"`
	Name  string `json:"name"` // заменить на уникальный ид в будущем
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

// Точка входа в игры, юзер отправляет нам свои данные, мы отдаем данные персонажа, уникальный ид или name через которое будет совершенно socket подключение
func InitHandler(woldMap *wmap.WorldMap) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, _ := ioutil.ReadAll(request.Body)
		rp := requestPlayer{}
		err := json.Unmarshal(body, &rp)
		if err != nil {
			log.WithFields(log.Fields{
				"package": "GameController",
				"func":    "InitHandler",
				"error":   err,
				"data":    body,
			}).Error("Error get player data")
		}
		writer.Header().Set("Content-Type", "application/json")

		player, exists := woldMap.GetPlayer(rp.Name)
		if exists {
			ok := player.ComparePassword(rp.Password)
			if ok {
				resPl := responsePlayer{Error: "null", X: player.X, Y: player.Y, Name: player.Name}
				res, err := json.Marshal(resPl)
				if err != nil {
					log.WithFields(log.Fields{
						"package": "GameController",
						"func":    "InitHandler",
						"error":   err,
						"data":    resPl,
					}).Error("Error Marshal player data")
					writer.Write([]byte("{Error: error server}"))
					return
				}
				writer.Write(res)
			}
			return
		}

		player = wmap.NewPlayer(rp.Name, rp.Password)
		woldMap.AddPlayer(player)
		resPl := responsePlayer{Error: "null", X: player.X, Y: player.Y, Name: player.Name}
		res, err := json.Marshal(resPl)
		if err != nil {
			log.WithFields(log.Fields{
				"package": "GameController",
				"func":    "InitHandler",
				"error":   err,
				"data":    resPl,
			}).Error("Error Marshal player data")
			writer.Write([]byte("{Error: error server}"))
			return
		}
		writer.Write(res)
	}
}
