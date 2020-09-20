package gamereducer

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lolodin/infworld/action"
	"github.com/lolodin/infworld/chunk"
	"github.com/lolodin/infworld/wmap"
	"log"
	"net"
	"sync"
)
const VIGILANCE = 512
var clients = make(map[chunk.Coordinate]map[string]*PlayerConn)
var eventsMutex = sync.Mutex{}

type Eventer interface {
	GetId() string
	chunk.Coordinater
}
type PlayerConn struct {
	wr      *wsutil.Writer
	encoder *json.Encoder
}

func (conn PlayerConn) sendData(i interface{}) {
	eventsMutex.Lock()
	e := conn.encoder.Encode(i)
	if e != nil {
		log.Println(e)
	}
	e = conn.wr.Flush()
	if e != nil {
		log.Println(e)
	}
	eventsMutex.Unlock()
}

//Добавляет соединение в чанк
func NewPlayerConn(conn net.Conn,  coordinater Eventer) {
    x,y := coordinater.GetCoordinate()
	idPlayer := coordinater.GetId()
	chunkID := wmap.GetChunkID(x, y)
	pc := PlayerConn{}
	wr := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(wr)
	pc.encoder = encoder
	pc.wr = wr
	if _, ok := clients[chunkID]; !ok {
		clients[chunkID] = map[string]*PlayerConn{}
	}

	clients[chunkID][idPlayer] = &pc

}

// Слушает события движения и отправляет данные нужным подключения
func ListnerMoveEvent(chEventMove <-chan Eventer, w *wmap.WorldMap) {
	for {
		select {
		case coord := <-chEventMove:
			fmt.Println(clients)
			pls := w.GetPlayers(coord)
			pls.Action = action.MOVE
			sendDataToChunck(coord, pls)

		}

	}
}

type ResMap struct {
	Action int            `json:"action"`
	Map    [9]chunk.Chunk `json:"gamemap"`
}

func ListnerGetMap(chGetMapEvent <-chan Eventer, w *wmap.WorldMap) {
	for {
		select {
		case getMapEvent := <-chGetMapEvent:
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
					eventsMutex.Lock()
					delete(sl, id)
					eventsMutex.Unlock()
					v.sendData(rs)
					if _, ok := clients[c]; ok {
						clients[c][id] = v
						continue
					}
					clients[c] = map[string]*PlayerConn{}
					clients[c][id] = v

				}
			}
		}
	}
}
func ListnerTreeEvent(chGetMapEvent <-chan Eventer, w *wmap.WorldMap) {
	for {
		select {
		case event :=<-chGetMapEvent:
			id:= event.GetId()
			w.Treehandler(event, id)
		}
	}


}

// Для функций изменяющих данные,
func sendDataToChunck(e chunk.Coordinater, data interface{}) {
	x,y := e.GetCoordinate()
	currentChunk := wmap.GetChunkID(x,y)
	m:=wmap.GetCurrentPlayerMap(currentChunk)
	for _, chunkID := range m {
		for _, conn := range clients[chunkID] {
			conn.sendData(data)
		}
	}
}