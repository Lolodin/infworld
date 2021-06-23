package wmap

import (
	"fmt"
	"github.com/lolodin/infworld/chunk"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
)

const VIGILANCE = 512
const TREEDIST = 10

type Mover interface {
	chunk.Coordinater
	GetId() string
}

type WorldMap struct {
	sync.Mutex
	Chunks map[chunk.Coordinate]*chunk.Chunk
	Player map[string]*Player
}

//Возвращает чанк соответствующий координатам, возвращает ошибку если такого чанка не существует
func (w *WorldMap) GetChunk(coordinate chunk.Coordinate) (chunk.Chunk, error) {
	c, ok := w.Chunks[coordinate]
	if ok == true {
		return *c, nil
	} else {
		return *c, fmt.Errorf("Chunck is not Exist")
	}
}

// Добавлявет в мир чанк
func (w *WorldMap) AddChunk(coordinate chunk.Coordinate, chunk chunk.Chunk) {

	isExist := w.isChunkExist(coordinate)
	if isExist {
		return
	} else {
		w.Lock()
		w.Chunks[coordinate] = &chunk
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
	world.Chunks = make(map[chunk.Coordinate]*chunk.Chunk)
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
func (w *WorldMap) MovePlayer(m Mover) *chunk.Coordinate {
	id := m.GetId()

	w.Lock()
	p, ok := w.Player[id]
	w.Unlock()
	x, y := m.GetCoordinate()
	if ok {
		x = p.X + x*p.speed
		y = p.Y + y*p.speed
		t := chunk.Coordinate{x, y}
		if b := w.CheckBusyTile(t); b {
			fmt.Println(b, "busy")
			return nil
		}
		p.X = x
		p.Y = y

		return &chunk.Coordinate{x, y}

	} else {
		fmt.Println("Player is not Exile")
	}
	return nil
}

//Получаем координаты от кого пришел запрос и получаем данные для кого эти изменения актуальны
func (w *WorldMap) GetPlayers(coordinater chunk.Coordinater) Players {
	pls := Players{}
	w.Lock()
	for _, P := range w.Player {
		if ok := CalcDistance(P, coordinater, VIGILANCE); ok {
			pls.P = append(pls.P, *P)
		}
	}
	w.Unlock()
	return pls
}

// return true if Tree busy tile
func (w *WorldMap) CheckBusyTile(coordinater chunk.Coordinater) bool {
	x, y := coordinater.GetCoordinate()
	tile := CurrentTile(coordinater)
	chunkId := GetChunkID(x, y)

	if b, ok := w.Chunks[chunkId].Map[tile]; ok {
		return b.Busy
	}
	return false
}

//Возвращает Player и bool == true если игрок с таким id/name есть
func (w *WorldMap) GetPlayer(name string) (*Player, bool) {
	w.Lock()
	pl, ok := w.Player[name]
	w.Unlock()
	if ok {
		return pl, ok
	} else {
		return nil, ok
	}
}

//Врзвращает координаты тайла которому принадлежит область
func CurrentTile(coordinater chunk.Coordinater) chunk.Coordinate {
	x, y := coordinater.GetCoordinate()
	tileX := float64(x) / float64(chunk.TILE_SIZE)
	tileY := float64(y) / float64(chunk.TILE_SIZE)
	if tileX < 0 {
		x = int(math.Ceil(tileX))
	} else {
		x = int(math.Floor(tileX))
	}
	if tileY < 0 {
		y = int(math.Ceil(tileY))
	} else {
		y = int(math.Floor(tileY))
	}
	var resX, resY int
	if x == 1 || x == -1 {
		resX = morthxy(x)
	}
	if y == 1 || y == -1 {
		resY = morthxy(y)
	}
	if resX == 0 {
		if x < 0 {
			resX = x*16 - 8
		} else {
			resX = x*16 + 8
		}

	}
	if resY == 0 {
		if y < 0 {
			resY = y*16 - 8
		} else {
			resY = y*16 + 8
		}
	}

	return chunk.Coordinate{X: resX, Y: resY}
}

// изменяет координаты при int  1 и -1 для функции CurrentTile
func morthxy(x int) int {
	if x == 1 {
		x = x * 8
	}
	if x == -1 {
		x = x * 8
	}
	return x
}
func (w *WorldMap) DestroyTree(TreeCoord chunk.Coordinater, idPlayer string) bool {
	player, ok := w.GetPlayer(idPlayer)
	if !ok {
		return false
	}
	if ok := CalcDistance(player, TreeCoord, TREEDIST); !ok {
		return false
	}
	x, y := TreeCoord.GetCoordinate()
	id := GetChunkID(x, y)
	w.Lock()
	w.Chunks[id].DestroyTree(chunk.Coordinate{X: x, Y: y})
	tile := w.Chunks[id].Map[chunk.Coordinate{X: x, Y: y}]
	tile.Busy = false
	w.Unlock()

	return true
}

//Возвращает true если для данных координат подходит дистанция V; coordinater объект, targer цель между которыми высчитывается дистанция
func CalcDistance(coordinater chunk.Coordinater, target chunk.Coordinater, v int) bool {
	x, y := coordinater.GetCoordinate()
	x1, y1 := target.GetCoordinate()
	P := chunk.Coordinate{X: x, Y: y}
	calX1 := P.X - v
	calX2 := P.X + v

	calY1 := P.Y - v
	calY2 := P.Y + v
	if (x1 > calX1 && x1 < calX2) && (y1 > calY1 && y1 < calY2) {
		return true
	}
	return false
}
