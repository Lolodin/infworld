package wmap

import (
	"github.com/lolodin/infworld/chunk"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
)

type Mover interface {
	chunk.Coordinater
	GetId() string
}

type WorldMap struct {
	sync.Mutex
	Chunks map[chunk.Coordinate]chunk.Chunk
	Player map[string]*Player

}

//Возвращает чанк соответствующий координатам, возвращает ошибку если такого чанка не существует
func (w *WorldMap) GetChunk(coordinate chunk.Coordinate) (chunk.Chunk, error) {
	c, ok := w.Chunks[coordinate]
	if ok == true {
		return c, nil
	} else {
		return c, fmt.Errorf("Chunck is not Exist")
	}

}

// Добавлявет в мир чанк
func (w *WorldMap) AddChunk(coordinate chunk.Coordinate, chunk chunk.Chunk) {

	isExist := w.isChunkExist(coordinate)
	if isExist {
		return
	} else {
		w.Lock()
		w.Chunks[coordinate] = chunk
		log.WithFields(log.Fields{
			"package":  "WorldMap",
			"func":     "AddChunk",
			"Chunk":    chunk,
			"map Tree": chunk.Tree,
		}).Info("Create new Chunk")
		w.Unlock()
	}

}

//Проверяет, существует чанк в мире или нет
func (w *WorldMap) isChunkExist(coordinate chunk.Coordinate) bool {
	_, ok := w.Chunks[coordinate]

	return ok
}

func NewCacheWorldMap() WorldMap {
	world := WorldMap{}
	world.Chunks = make(map[chunk.Coordinate]chunk.Chunk)
	world.Player = make(map[string]*Player)
	return world
}

//Добавляем нового игрока в карту
func (w *WorldMap) AddPlayer(player *Player) {

	_, ok := w.Player[player.Name]
	if !ok {
		fmt.Println(player.Name)
		w.Lock()
		w.Player[player.Name] = player
		w.Unlock()
	} else {
		fmt.Println("Relogin: " + player.Name)
	}
}

// Обновляем данные персонажа в мире
func (w *WorldMap) MovePlayer(m Mover) {
	id:= m.GetId()
	w.Lock()
	p, ok := w.Player[id]
	w.Unlock()
	x,y :=m.GetCoordinate()
	if ok {
		x = p.X + x*p.speed
		y = p.Y + y*p.speed
		t:=chunk.Coordinate{x,y}
		if b :=w.CheckBusyTile(t); b {
			fmt.Println(b, "busy")
			return
		}
		p.X = x
		p.Y = y



	} else {
		fmt.Println("Player is not Exile")
	}

}

//map players
func (w *WorldMap) GetPlayers() Players {
	pls := Players{}
	w.Lock()
	for _, P := range w.Player {
		pls.P = append(pls.P, *P)
	}
	w.Unlock()
	return pls
}

// return true if Tree busy tile
func (w *WorldMap) CheckBusyTile(coordinater chunk.Coordinater) bool {
	x,y:= coordinater.GetCoordinate()
	tile := CurrentTile(coordinater)
	fmt.Println(tile, "debug", x,y)
	chunkId := GetChunkID(x,y)
	w.Lock()
	defer w.Unlock()
	b := w.Chunks[chunkId].Map[tile].Busy
	return b

}

//Получить player
func (w *WorldMap) GetPlayer(name string) (*Player, bool) {
	pl, ok := w.Player[name]
	if ok {
		return pl, ok
	} else {
		return &Player{}, ok
	}

}
//Врзвращает координаты тайла которому принадлежит область
func CurrentTile(coordinater chunk.Coordinater) (chunk.Coordinate) {
	x,y := coordinater.GetCoordinate()
	tileX := float64(x)/float64(chunk.TILE_SIZE)
	tileY := float64(y)/float64(chunk.TILE_SIZE)
	if tileX<0 {
		x= int(math.Ceil(tileX))
	} else {
		x= int(math.Floor(tileX))
	}
	if tileY<0 {
		y= int(math.Ceil(tileY))
	} else {
		y= int(math.Floor(tileY))
	}
	var resX, resY int
	if x == 1 || x ==-1 {
		resX = morthxy(x)
	}
	if y == 1 || y ==-1 {
		resY = morthxy(y)
	}
	if resX == 0 {
		if x<0 {
			resX = x*16-8
		} else {
			resX = x*16+8
		}

	}
	if resY == 0 {
		if y<0 {
			resY = y*16-8
		} else {
			resY = y*16+8
		}
	}

	return chunk.Coordinate{X:resX, Y: resY}


}
// изменяет координаты при int  1 и -1 для функции CurrentTile
func morthxy(x int) int {
	if x == 1 {
		x = x*8
	}
	if x == -1 {
		x = x*8
	}
	return x

}
func (w *WorldMap) Treehandler(coordinater chunk.Coordinater)  {
	x,y := coordinater.GetCoordinate()
	id:=GetChunkID(x,y)
	w.Lock()
	tree := w.Chunks[id].Tree[chunk.Coordinate{X:x, Y:y}]
	tile := w.Chunks[id].Map[chunk.Coordinate{X:x, Y:y}]
	w.Unlock()
	tile.Busy = false



}