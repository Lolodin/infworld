package chunk

import (
	"github.com/lolodin/infworld/perlinNoise"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var Mutex = sync.Mutex{}

const CHUNKIDSIZE = 16
const TILE_SIZE = 16
const CHUNK_SIZE = 16 * 16
const PERLIN_SEED float32 = 2300

type Coordinater interface {
	GetCoordinate() (x, y int)
}

// Чанк который хранит тайтлы и другие игровые объекты
type Chunk struct {
	ChunkID [2]int
	Map     map[Coordinate]*Tile
	Tree    map[Coordinate]*Tree
}

/*
Тайтл игрового мира
*/
type Tile struct {
	Key  string `json:"key"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Busy bool
}

// Освобождаем чанк для движения

/*
Универсальная структура для хранения координат
*/
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (ch *Chunk) DestroyTree(coordinate Coordinate) {
	Mutex.Lock()
	delete(ch.Tree, coordinate)
	Mutex.Unlock()
}
func (t *Tile) TileClear() {
	t.Busy = false
}

func (c Coordinate) GetCoordinate() (x, y int) {
	return c.X, c.Y
}

func (c Coordinate) MarshalText() ([]byte, error) {
	return []byte("[" + strconv.Itoa(c.X) + "," + strconv.Itoa(c.Y) + "]"), nil
}

func fillChunk(chunkMap map[Coordinate]*Tile, treeMap map[Coordinate]*Tree, posX float32, posY float32, x int, y int) {
	tile := Tile{}
	tile.X = int(posX)
	tile.Y = int(posY)
	var tree *Tree

	perlinValue := perlinNoise.Noise(posX/PERLIN_SEED, posY/PERLIN_SEED)
	switch {
	case perlinValue < -0.012:
		tile.Key = "Water"
		tile.Busy = true
	case perlinValue >= -0.012 && perlinValue < 0:
		tile.Key = "Sand"
	case perlinValue >= 0 && perlinValue <= 0.5:
		tile.Key = "Ground"
		rand.Seed(int64(time.Now().Nanosecond() + x - y))
		randomTree := rand.Float32()
		if randomTree > 0.95 {
			tree = NewTree(Coordinate{X: tile.X, Y: tile.Y})
		}
	case perlinValue > 0.5:
		tile.Key = "Mount"
	}

	if tree != nil {
		treeMap[Coordinate{X: tree.X, Y: tree.Y}] = tree
		tile.Busy = true
		tree = nil
	}

	chunkMap[Coordinate{X: tile.X, Y: tile.Y}] = &tile
}

/*
Создает карту чанка из тайтлов, генерирует карту на основе координаты чанка
Например [1,1]
*/
func NewChunk(idChunk Coordinate) Chunk {
	log.WithFields(log.Fields{
		"package": "Chunk",
		"func":    "NewChunk",
		"idChunk": idChunk,
	}).Info("Create new Chunk")

	chunk := Chunk{ChunkID: [2]int{idChunk.X, idChunk.Y}}
	var chunkXMax, chunkYMax int
	var chunkMap map[Coordinate]*Tile
	var treeMap map[Coordinate]*Tree
	chunkMap = make(map[Coordinate]*Tile)
	treeMap = make(map[Coordinate]*Tree)
	chunkXMax = idChunk.X * CHUNK_SIZE
	chunkYMax = idChunk.Y * CHUNK_SIZE

	switch {
	case chunkXMax < 0 && chunkYMax < 0:
		{
			for x := chunkXMax + CHUNK_SIZE; x > chunkXMax; x -= TILE_SIZE {
				for y := chunkYMax + CHUNK_SIZE; y > chunkYMax; y -= TILE_SIZE {
					posX := float32(x - (TILE_SIZE / 2))
					posY := float32(y + (TILE_SIZE / 2))
					fillChunk(chunkMap, treeMap, posX, posY, x, y)
				}
			}
		}
	case chunkXMax < 0:
		{
			for x := chunkXMax + CHUNK_SIZE; x > chunkXMax; x -= TILE_SIZE {
				for y := chunkYMax - CHUNK_SIZE; y < chunkYMax; y += TILE_SIZE {
					posX := float32(x - (TILE_SIZE / 2))
					posY := float32(y + (TILE_SIZE / 2))
					fillChunk(chunkMap, treeMap, posX, posY, x, y)
				}
			}
		}
	case chunkYMax < 0:
		{
			for x := chunkXMax - CHUNK_SIZE; x < chunkXMax; x += TILE_SIZE {
				for y := chunkYMax + CHUNK_SIZE; y > chunkYMax; y -= TILE_SIZE {
					posX := float32(x + (TILE_SIZE / 2))
					posY := float32(y - (TILE_SIZE / 2))
					fillChunk(chunkMap, treeMap, posX, posY, x, y)
				}
			}
		}
	default:
		{
			for x := chunkXMax - CHUNK_SIZE; x < chunkXMax; x += TILE_SIZE {
				for y := chunkYMax - CHUNK_SIZE; y < chunkYMax; y += TILE_SIZE {
					posX := float32(x + (TILE_SIZE / 2))
					posY := float32(y + (TILE_SIZE / 2))
					fillChunk(chunkMap, treeMap, posX, posY, x, y)
				}
			}
		}

	}

	chunk.Map = chunkMap
	chunk.Tree = treeMap

	return chunk
}
