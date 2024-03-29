package main

import (
	"github.com/lolodin/infworld/gamereducer"
	"github.com/lolodin/infworld/gcontrl"
	"github.com/lolodin/infworld/playerhand"
	"github.com/lolodin/infworld/wmap"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var (
	chEventMove   = make(chan gamereducer.Eventer)
	chEventGetMap = make(chan gamereducer.Eventer)
	chEventTree   = make(chan gamereducer.Eventer)
)

func init() {
	filelog, e := os.Create("log")
	if e != nil {
		panic("error create log file")
	}
	log.SetOutput(filelog)
}
func main() {
	log.WithFields(log.Fields{
		"package": "main",
		"func":    "main",
	}).Info("Server start")
	World := wmap.NewCacheWorldMap()
	http.HandleFunc("/init", gcontrl.InitHandler(&World))
	http.HandleFunc("/map", gcontrl.MapHandler(&World))
	//Сигнал действия юзера

	go gamereducer.OnMove(chEventMove, &World)
	go gamereducer.OnGetMap(chEventGetMap, &World)
	go gamereducer.OnTree(chEventTree, &World)
	http.HandleFunc("/player", playerhand.PlayerHandler(&World, chEventMove, chEventGetMap, chEventTree))
	http.HandleFunc("/", indexHandler)

	//static
	http.Handle("/node_modules/phaser/dist/", http.StripPrefix("/node_modules/phaser/dist/", http.FileServer(http.Dir("./node_modules/phaser/dist/"))))
	http.Handle("/Client/", http.StripPrefix("/Client/", http.FileServer(http.Dir("./Client/"))))
	http.Handle("/Client/content/", http.StripPrefix("/Client/content/", http.FileServer(http.Dir("./Client/content/"))))

	port := "8080"
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "main",
			"func":    "main",
			"error":   err,
		}).Fatal("Error start server")
	}

}

// Обработчик для index.html, здесь мы просто отдаем клиент пользователю
func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	err := t.Execute(w, "index")
	if err != nil {
		log.WithFields(log.Fields{
			"package": "main",
			"func":    "indexHandler",
			"error":   err,
		}).Error("Error get index.html")
	}
}
