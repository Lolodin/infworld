package wmap

import (
	"sync"
)

type Player struct {
	mut      sync.Mutex
	Name     string `json:"Name"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	password string

	//	AnimKey string
	// Игровые характеристики
	vigilance int //стандарт 512 дальность видимости игроков
	speed    int //скорость, стандарт 1

}
type Players struct {
	Action int `json:"action"`
	P []Player `json:"players"`
}

func NewPlayer(n, password string) *Player {
	p := Player{}
	p.X = 16
	p.Y = 16
	p.Name = n
	p.password = password
	p.speed = 3
	p.vigilance = 512
	return &p

}
func (p *Player) GetCoordinate() (x, y int) {
	return p.X,  p.Y
}

func (p *Player) SetPassword(pass string) {
	p.password = pass
}
func (p Player) GetPassword() string {
	return p.password
}
func (p Player) GetId() string {
	return p.Name
}
// bool true if pass == player.password
func (p *Player) ComparePassword(pass string) bool {
	if pass == p.password {
		return true
	} else {
		return false
	}
}

