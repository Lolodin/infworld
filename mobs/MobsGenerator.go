package mobs

import (
     "github.com/lolodin/infworld/chunk"
	"time"
)

type MobGenerator struct {
	chunk.Coordinate
	tik        time.Duration
	ListMob    []string
	CurrentMob *Mob
}

//func NewMobGenerator() MobGenerator {
//
//}

// Запуск в конструкторе
func (g *MobGenerator) Generation() {
	for {
		if g.CurrentMob == nil {

		}
		time.Sleep(g.tik * time.Second)
	}
}

//func(g *MobGenerator) newMob() {
//	l:= len(g.ListMob)
//	el:=rand.Intn(l)
//	m:= NewMob()
//	g.CurrentMob
//}
