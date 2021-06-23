import {GETMAP, MOVE, TREE} from './Action.js'

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
            switch (data.action) {
                case MOVE:
                    this.Players.DrawPlayer(data.players)
                    break;
                case GETMAP:
                    let gmap = data.gamemap
                    this.Map.drawMapController(gmap)
                    break;
                case TREE:
                    if (!data.result) {
                        return;
                    }
                    this.Map.destroyTree(data.x, data.y);
                    break;
                default:
                    console.log("action not found")
            }
        }
    }
}