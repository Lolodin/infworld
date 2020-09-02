package wmap

import (
	"github.com/lolodin/infworld/chunk"
	"sync"
)

type Player struct {
	mut      sync.Mutex
	Name     string `json:"Name"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	password string
	speed    int
	move     bool
	//	AnimKey string

}
type Players struct {
	P []Player `json:"players"`
}

func NewPlayer(n, password string) *Player {
	p := Player{}
	p.X = 16
	p.Y = 16
	p.Name = n
	p.password = password
	p.speed = 15
	p.move = false

	return &p

}
func (p *Player) GetPlayerXY() chunk.Coordinate {
	return chunk.Coordinate{X: p.X, Y: p.Y}
}

func (p *Player) SetPassword(pass string) {
	p.password = pass
}
func (p *Player) GetPassword() string {
	return p.password
}
// bool true if pass == player.password
func (p *Player) ComparePassword(pass string) bool {
	if pass == p.password {
		return true
	} else {
		return false
	}
}

