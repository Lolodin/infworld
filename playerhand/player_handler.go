package playerhand

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lolodin/infworld/action"
	"github.com/lolodin/infworld/chunk"
	"github.com/lolodin/infworld/gamereducer"
	"github.com/lolodin/infworld/wmap"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

/*
{action: Move, id = name, x:+1, y:+1}
*/
type PlayerResponseMOVE struct {
	Id string `json:"id"`
	chunk.Coordinate
}
type PlayerResponseTREE struct {
	Id string `json:"id"`
	chunk.Coordinate
}
type PlayerResponseGETMAP struct {
	Id string `json:"id"`
	chunk.Coordinate
}

func(p PlayerResponseMOVE) GetId() string {
	return p.Id
}
func(p PlayerResponseGETMAP) GetId() string {
	return p.Id
}



func PlayerHandler(W *wmap.WorldMap, eventMove chan<- chunk.Coordinater, EventGetMap chan <-gamereducer.Eventer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, e :=ws.UpgradeHTTP(r,w)
		if e!=nil {
			fmt.Println(e)
		}
		msg, _, err:=wsutil.ReadClientData(conn)
		if err != nil {
			fmt.Println(err)
		}
		getId := PlayerResponseMOVE{}
		json.Unmarshal(msg, &getId)

		//Добавить реальные x, y
		gamereducer.NewPlayerConn(conn, 16, 16, getId.Id)
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(log.Fields{
					"package": "GameController",
					"func":    "PlayerHandler",
					"error":   err,
				}).Error("Error ws")
			}

		}()





		//Game Loop
		log.WithFields(log.Fields{
			"package": "GameController",
			"func":    "PlayerHandler",
		}).Info("Connect player")


		//ws handler
go func() {

	for {
		
		msg, _, err := wsutil.ReadClientData(conn)
		if err != nil {
			panic("conn cancel")
		}
		str:=string(msg[10:13])
		a, e := strconv.Atoi(str)
		if e != nil {
			fmt.Println("Error conv json")
		}
		fmt.Println(str)
		
		switch a {
		case action.MOVE:
				req := PlayerResponseMOVE{}
				json.Unmarshal(msg, &req)
				coord:=W.MovePlayer(req)
				if coord == nil {
					continue
				}
				eventMove <- coord
			log.WithFields(log.Fields{
				"func":    "PlayerHandler",
				"Player": req,
			}).Info("MOVE")
		case action.TREE:
				req:= PlayerResponseTREE{}
				json.Unmarshal(msg, &req)
				fmt.Println(req)
		case action.GET_MAP:
			req:= PlayerResponseGETMAP{}
			json.Unmarshal(msg, &req)
			EventGetMap<-req


		}

		/* ответ сервера - положение игроков*/



	}
}()

	}
}