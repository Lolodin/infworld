package gamereducer

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lolodin/infworld/action"
	"github.com/lolodin/infworld/chunk"
	"github.com/lolodin/infworld/wmap"

	"net"
	"sync"
)

// TODO дальность видимости
const VIGILANCE = 512
// Структура хранения подключений
var clients = clientmap{}
//
var eventsMutex = sync.Mutex{}
// Для хранения подключений по чанкам
type clientmap map[chunk.Coordinate]map[string]*PlayerConn

// Каждый эвент должен реализовать данный интерфейс
type Eventer interface {
	GetId() string
	chunk.Coordinater
}
// Подключение с игроком
type PlayerConn struct {
	wr      *wsutil.Writer
	encoder *json.Encoder
}
// Структура для отправки данных по Websocket
func (conn PlayerConn) sendData(i interface{}) {
	eventsMutex.Lock()
	e := conn.encoder.Encode(i)
	if e != nil {

	}
	e = conn.wr.Flush()
	if e != nil {

	}
	eventsMutex.Unlock()
}

//Добавляет соединение в чанк
func NewPlayerConn(conn net.Conn,  coordinater Eventer) {

   // x,y := coordinater.GetCoordinate()
	idPlayer := coordinater.GetId()
	chunkID := wmap.GetChunkID(coordinater)
	pc := PlayerConn{}
	wr := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(wr)
	pc.encoder = encoder
	pc.wr = wr
	eventsMutex.Lock()
	if _, ok := clients[chunkID]; !ok {
		clients[chunkID] = map[string]*PlayerConn{}
	}

	clients[chunkID][idPlayer] = &pc
	eventsMutex.Unlock()

}

// Слушает события движения и отправляет данные нужным подключения
func ListnerMoveEvent(chEventMove <-chan Eventer, w *wmap.WorldMap) {
	for {
		select {
		case coord := <-chEventMove:

			pls := w.GetPlayers(coord)
			pls.Action = action.MOVE
			clients.sendDataToChunck(coord, pls)


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
			id := getMapEvent.GetId()
			c := wmap.GetChunkID(getMapEvent)
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
func ListnerPlayerDisconnect(chEventDisconnect <-chan Eventer, w *wmap.WorldMap) {
	for {
		select {
		case event:=<-chEventDisconnect:
	    clients.deleteConn(event)



		}
	}
}

// Для функций изменяющих данные,
func(m clientmap) sendDataToChunck(e chunk.Coordinater, data interface{}) {
	currentChunk := wmap.GetChunkID(e)
	mp:=wmap.GetCurrentPlayerMap(currentChunk)
	eventsMutex.Lock()
	for _, chunkID := range mp {
		for _, conn := range clients[chunkID] {
			conn.sendData(data)
		}
	}
	eventsMutex.Unlock()
}

func(m clientmap) deleteConn(eventer Eventer) {
	playerID := eventer.GetId()
	chunCoord:=wmap.GetChunkID(eventer)
	eventsMutex.Lock()
	delete(m[chunCoord],playerID)
	eventsMutex.Unlock()
}