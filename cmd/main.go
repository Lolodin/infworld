package main

import (
	"github.com/lolodin/infworld/gamereducer"
	"github.com/lolodin/infworld/gcontrl"
	"github.com/lolodin/infworld/playerhand"
	"github.com/lolodin/infworld/wmap"
	"github.com/sirupsen/logrus"
	"html/template"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)
var (
chEventMove = make(chan gamereducer.Eventer)
chEventGetMap = make(chan gamereducer.Eventer)
chEventTree = make(chan gamereducer.Eventer)
chEventDisconnect = make(chan gamereducer.Eventer)
)


func main() {
	log:= logrus.New()
	f, err:= os.Create("runlog")
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	World := wmap.NewCacheWorldMap()
	http.HandleFunc("/init", gcontrl.InitHandler(&World))
	http.HandleFunc("/map", gcontrl.Map_Handler(&World))
	//Сигнал действия юзера

	go gamereducer.ListnerMoveEvent(chEventMove, &World)
	go gamereducer.ListnerGetMap(chEventGetMap, &World)
	//go gamereducer.ListnerTreeEvent(chEventTree, &World)
	go gamereducer.ListnerPlayerDisconnect(chEventDisconnect, &World)
	http.HandleFunc("/player", playerhand.PlayerHandler(&World,chEventDisconnect, chEventMove, chEventGetMap, chEventTree))
	http.HandleFunc("/", indexHandler)

	//static
	http.Handle("/node_modules/phaser/dist/", http.StripPrefix("/node_modules/phaser/dist/", http.FileServer(http.Dir("./node_modules/phaser/dist/"))))
	http.Handle("/Client/", http.StripPrefix("/Client/", http.FileServer(http.Dir("./Client/"))))
	http.Handle("/Client/content/", http.StripPrefix("/Client/content/", http.FileServer(http.Dir("./Client/content/"))))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

// Обработчик для index.html, здесь мы просто отдаем клиент пользователю
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	err := t.Execute(w, "index")
	if err != nil {
		log.Fatal(err)
	}
}
