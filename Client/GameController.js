

export default class Controller {
    constructor(scene, map, players) {
        this.websocket = scene.websocket
        this.Map = map
        this.Players = players
    }
    listner() {
        this.websocket.onmessage = (e)=> {
            console.log(e)
            let players = e.data
            players = JSON.parse(players)
            console.log(players, "получаем данные от сервера")
            this.Players.DrawPlayer(players.players)
        }
    }
}