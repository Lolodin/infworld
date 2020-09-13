package gamereducer

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lolodin/infworld/action"
	"github.com/lolodin/infworld/chunk"
	"github.com/lolodin/infworld/wmap"
	"log"
	"net"
)


var clients = make(map[chunk.Coordinate]map[string]*PlayerConn)

type Eventer interface {
	GetId() string
	chunk.Coordinater
}
type PlayerConn struct {
	wr *wsutil.Writer
	encoder *json.Encoder
}


func(conn PlayerConn) sendData(i interface{}) {
	e:=conn.encoder.Encode(i)
	if e!= nil {
		log.Println(e)
	}
	e =conn.wr.Flush()
	if e!= nil {
		log.Println(e)
	}
}
//Добавляет соединение в чанк
func NewPlayerConn(conn net.Conn, x,y int, idPlayer string) {
	chunkID:= wmap.GetChunkID(x,y)
	pc:= PlayerConn{}
	wr:= wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(wr)
	pc.encoder = encoder
	pc.wr = wr
	clients[chunkID] = map[string]*PlayerConn{}
	clients[chunkID][idPlayer] =&pc

}
// Слушает события движения и отправляет данные нужным подключения 
func ListnerMoveEvent(chEventMove <-chan chunk.Coordinater, w *wmap.WorldMap) {
	for  {
		select {
		case coord:=<-chEventMove:
			pls:= w.GetPlayers(coord)
			pls.Action = action.MOVE
			for _,arrayconn := range clients {
				for _, conn := range arrayconn {
						conn.sendData(pls)
					}
				}

			}

		}
	}



type ResMap struct {
	Action int `json:"action"`
	Map [9]chunk.Chunk `json:"gamemap"`
}
func ListnerGetMap(chGetMapEvent<- chan Eventer,  w *wmap.WorldMap)  {
for {
	select {
	case getMapEvent:=<-chGetMapEvent:
		rs := ResMap{}
		x, y := getMapEvent.GetCoordinate()
		id := getMapEvent.GetId()
		c := wmap.GetChunkID(x, y)
		d := wmap.GetCurrentPlayerMap(c)
		m := wmap.GetPlayerDrawChunkMap(d, w)
		rs.Map = m
		rs.Action = action.GET_MAP
		for _, sl := range clients {
			if v, ok := sl[id]; ok {
				delete(sl,id)
				v.sendData(rs)
				if _, ok:=clients[c]; ok {
					clients[c][id] = v
				}
				clients[c] = map[string]*PlayerConn{}
				clients[c][id] = v
			}
		}
	}
}
}