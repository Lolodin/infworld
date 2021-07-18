package gcontrl

import (
	"encoding/json"
	"github.com/lolodin/infworld/wmap"
	"io/ioutil"
	"net/http"
)

type requestPlayer struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type responsePlayer struct {
	Error string `json:"error"`
	Name  string `json:"name"` // заменить на уникальный ид в будущем
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

// Точка входа в игрy, юзер отправляет нам свои данные, мы отдаем данные персонажа, уникальный ид или name через которое будет совершенно socket подключение
func InitHandler(W *wmap.WorldMap) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		body, _ := ioutil.ReadAll(r.Body)
		rp := requestPlayer{}
		err := json.Unmarshal(body, &rp)
		if err != nil {

		}
		w.Header().Set("Content-Type", "application/json")
		p, exile := W.GetPlayer(rp.Name)
		if exile {
			ok := p.ComparePassword(rp.Password)
			if ok {
				resPl := responsePlayer{Error: "null", X: p.X, Y: p.Y, Name: p.Name}
				res, err := json.Marshal(resPl)
				if err != nil {
					w.Write([]byte("{Error: error server}"))
					return
				}
				w.Write(res)
				return
			}
		} else {
			p := wmap.NewPlayer(rp.Name, rp.Password)
			W.AddPlayer(p)
			resPl := responsePlayer{Error: "null", X: p.X, Y: p.Y, Name: p.Name}
			res, err := json.Marshal(resPl)
			if err != nil {
				w.Write([]byte("{Error: error server}"))
				return
			}
			w.Write(res)
			return

		}

	}
}
