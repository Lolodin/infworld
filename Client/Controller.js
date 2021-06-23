import {GETMAP, MOVE} from './Action.js'

export default class Controller {
    constructor(scene, map, players) {
        this.ID = scene.ID
        this.websocket = scene.websocket
        this.Map = map
        this.Players = players
    }

    listener() {
        this.websocket.onmessage = (e) => {
            let serverAnswer = e.data
            let data = JSON.parse(serverAnswer)
            console.log(serverAnswer, data, "debug")
            switch (data.action) {
                case MOVE:
                    console.log("MOVE", data)
                    this.Players.DrawPlayer(data.players)
                    break;
                case GETMAP:
                    console.log("GETMAP", data)
                    let gmap = data.gamemap
                    this.Map.drawMapController(gmap)
                    break;
                default:
                    console.log("action not found")
            }
        }
    }
}