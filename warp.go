package main

import "log"

type WarpDest struct {
	warpID int
	mapID  int
	gridX  int
	gridY  int
}

type WarpManager struct {
	game         *Game
	warpDestList []*WarpDest
}

func NewWarpManager(game *Game) *WarpManager {
	wm := &WarpManager{}
	wm.warpDestList = []*WarpDest{}
	wm.warpDestList = append(wm.warpDestList, &WarpDest{0, 1, 3, 3})
	wm.warpDestList = append(wm.warpDestList, &WarpDest{1, 2, 3, 3})
	wm.warpDestList = append(wm.warpDestList, &WarpDest{2, 3, 3, 3})
	wm.warpDestList = append(wm.warpDestList, &WarpDest{3, 4, 3, 3})
	wm.warpDestList = append(wm.warpDestList, &WarpDest{4, 5, 3, 3})
	wm.game = game
	return wm
}

func (wm *WarpManager) warpPlayerToWarpID(warpID int) {
	for _, v := range wm.warpDestList {
		if v.warpID == warpID {
			wm.game.level = v.mapID
			//wm.game.tileMap.loadCurrentLevelMapFromFile()
			//wm.game.pickupManager.loadDataFromFile()
			//wm.game.fidgetManager.loadDataFromFile()
			wm.game.loadLevel(wm.game.level)
			wm.game.player.warpPlayerToGridLocation(v.gridX, v.gridY)
			return
		}

	}
	log.Println("warpPlayerToWarpID matching zone not found ", warpID)

}
