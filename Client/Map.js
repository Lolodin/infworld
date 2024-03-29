import {TREE} from "./Action.js";

export {GameMap}

class GameMap {
    constructor(scene) {
        this.scene = scene
    }

    drawMapController(requestMapServer) {
        for (let i = 0; i < 9; i++) {
            this.drawTileChunk(requestMapServer[i].Map, requestMapServer[i].ChunkID)
            this.drawTree(requestMapServer[i].Tree, requestMapServer[i].ChunkID)
        }
    }

    drawTileChunk(chunk, chunkID) {
        // Check chunk is Load
        if (this.scene.LoadChunks[chunkID] == true) {
            return
        }
        // add chunk Group for tiles
        // Load chunk true
        this.scene.CurrentMap[chunkID] = this.scene.add.group()
        this.scene.LoadChunks[chunkID] = true
        for (let coordTile in chunk) {
            let tile

            if (chunk[coordTile].key == "Water") {
                tile = this.scene.add.sprite(chunk[coordTile].x, chunk[coordTile].y, chunk[coordTile].key).play('water', true);
            } else {
                tile = this.scene.add.image(chunk[coordTile].x, chunk[coordTile].y, chunk[coordTile].key)
            }
            tile.setDepth(tile.y - 30)
            tile.setInteractive()
            this.scene.CurrentMap[chunkID].add(tile)

        }
        console.log(this.scene.CurrentMap)
    }

    drawTree(chunk, chunkID) {
        if (this.scene.LoadChunksTree[chunkID] === true) {
            return
        }
        this.scene.LoadChunksTree[chunkID] = true
        for (let coordTile in chunk) {
            let tree
            //Поднимаем дерево на закрытый тайл
            let yNumber
            if (chunk[coordTile].y < 0) {
                yNumber = -12
            } else {
                yNumber = 0
            }
            tree = this.scene.add.image(chunk[coordTile].x, chunk[coordTile].y + yNumber, chunk[coordTile].tree)
            tree.setDepth(chunk[coordTile].y + 5);
            tree.setRotation(chunk[coordTile].age / 5);
            tree.setInteractive();
            tree.on('clicked', this.clickTree, this);
            tree.ServerX = chunk[coordTile].x
            tree.ServerY = chunk[coordTile].y
            this.scene.CurrentMap[chunkID].add(tree)
        }
    }

    clearMap(newCoordinate) {
        for (let i = 0; i < this.scene.coordinate.length; i++) {
            let chunkIsNotExist = true
            newCoordinate.forEach((v) => {
                if (this.scene.coordinate[i][0] == v[0] && this.scene.coordinate[i][1] == v[1]) {
                    chunkIsNotExist = false
                }
            })

            if (chunkIsNotExist) {
                let c = this.scene.coordinate[i][0] + "," + this.scene.coordinate[i][1]
                delete this.scene.LoadChunks[c]
                delete this.scene.LoadChunksTree[c]
                console.log(this.scene.CurrentMap, this.scene.coordinate, c, "TEST")
                try {
                    this.scene.CurrentMap[c].clear(true, true)
                    delete this.scene.CurrentMap[c];
                } catch (e) {
                    console.log("Error clear map", e)
                }
            }
        }
    }

    clickTree(objTree) {
        let playerRequest = {action: TREE, id: this.scene.ID.Name, x: objTree.ServerX, y: objTree.ServerY}
        this.scene.websocket.send(JSON.stringify(playerRequest))
    }

    destroyTree(x, y) {
        const chunkId = this.scene.getChunkID(x, y);
        this.scene.LoadChunksTree[chunkId] = false;
        for (let one of this.scene.CurrentMap[chunkId].children.entries) {
            if (one.x !== x || one.y !== y) {
                continue;
            }

            if (["Oak", "Spruce"].includes(one.texture.key)) {
                one.destroy();
                break;
            }
        }
    }
}