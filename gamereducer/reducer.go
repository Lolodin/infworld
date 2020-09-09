package gamereducer

import (
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/lolodin/infworld/chunk"
	"github.com/lolodin/infworld/wmap"
	"log"
	"net"
)

var clients = make(map[chunk.Coordinate]*PlayerConn)

type PlayerConn struct {
	wr *wsutil.Writer
	encoder *json.Encoder
}
//Добавляет соединение в чанк
func NewPlayerConn(conn net.Conn, x,y int) {
	chunkID:= wmap.GetChunkID(x,y)
	pc:= PlayerConn{}
	wr:= wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
	encoder := json.NewEncoder(wr)
	pc.encoder = encoder
	pc.wr = wr
	clients[chunkID] =&pc

}

func ListnerMoveEvent(ch <-chan struct{}, w *wmap.WorldMap) {
	for  {
		select {
		case <-ch:
			pls:= w.GetPlayers()
			for _,conn := range clients {
				e:=conn.encoder.Encode(pls)
				if e!= nil {
					log.Println(e)
				}
				e =conn.wr.Flush()
				if e!= nil {
					log.Println(e)
				}
			}
		}
	}

}