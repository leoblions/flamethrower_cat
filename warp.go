package main

import (
	"fmt"
	"log"
	"path"
)

// data file columns: warpID mapID  gridX  gridY

type WarpDest struct {
	warpID int
	mapID  int
	gridX  int
	gridY  int
}

type WarpManager struct {
	game             *Game
	warpDestList     []*WarpDest
	lastUsedWarpZone int
}

const (
	WM_WARP_DATA_MASTER      = "warp_table_master.csv"
	WM_WRITE_FILE_IF_MISSING = true
)

func NewWarpManager(game *Game) *WarpManager {
	wm := &WarpManager{}
	wm.warpDestList = []*WarpDest{}
	wm.loadDataFromFile()
	wm.game = game
	return wm
}

func (wm *WarpManager) warpPlayerToWarpID(warpID int) {
	for _, v := range wm.warpDestList {
		if v.warpID == warpID {
			wm.game.level = v.mapID
			wm.game.loadLevel(wm.game.level)
			wm.game.player.warpPlayerToGridLocation(v.gridX, v.gridY)
			wm.lastUsedWarpZone = warpID
			wm.game.hud.levelString.updateText(fmt.Sprintf("Level: %d", v.mapID))
			return
		}

	}
	log.Println("warpPlayerToWarpID matching zone not found ", warpID)

}

func (wm *WarpManager) loadDataFromFile() error {
	fmt.Println("warp data load")
	wm.warpDestList = []*WarpDest{}
	warpDataFilePath := path.Join(GAME_LEVEL_DATA_DIR, WM_WARP_DATA_MASTER)
	numericData, err := loadDataListFromFile(warpDataFilePath)
	rows := len(numericData)

	if err != nil {
		wm.warpDestList = []*WarpDest{}
		if WM_WRITE_FILE_IF_MISSING {
			fmt.Printf("Warp data file missing %s creating it /n", warpDataFilePath)
			write2DIntListToFile(numericData, warpDataFilePath)
		}

		return err
	}
	if rows == 0 {
		wm.warpDestList = []*WarpDest{}
		fmt.Println("warp data empty")
		return nil
	}
	wm.warpDestList = []*WarpDest{}
	for i := 0; i < rows; i++ {
		data := numericData[i]
		pfTemp := &WarpDest{}
		pfTemp = &WarpDest{data[0], data[1], data[2], data[3]}
		wm.warpDestList = append(wm.warpDestList, pfTemp)

	}
	return nil
}
