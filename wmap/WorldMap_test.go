package wmap

import (
	"fmt"
	"github.com/lolodin/infworld/chunk"
	"testing"
)

type testCoord struct {
	X      int
	Y      int
	Result chunk.Coordinate
}

var tests = []testCoord{
	{X: 320, Y: 320, Result: chunk.Coordinate{X: 2, Y: 2}},
	{X: 2560, Y: 2560, Result: chunk.Coordinate{X: 3, Y: 3}},
}

func TestGetChankID(t *testing.T) {
	for _, testValue := range tests {
		t := GetChunkID(testValue.X, testValue.Y)
		fmt.Println(t == testValue.Result, t)
	}

}
func TestCurrentTile(t *testing.T) {
	c:=chunk.Coordinate{X:86, Y:35}
	x:=CurrentTile(c)
	if x.X != 88 || x.Y != 40 {
		t.Error("Coordinate not correct", x)
		return
	}
	t.Log("Test positive ok", x)
	c=chunk.Coordinate{X:-107, Y:-1352} // -104
	x=CurrentTile(c) //-104;-1352
	if x.X != -104 || x.Y != -1352 {
		t.Error("Coordinate not correct", x)
		return
	}
	t.Log("Test negative ok", x)


}
