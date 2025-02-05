package main

import (
	"fmt"
	"log"
	"path"
)

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

const (
	WM_WARP_DATA_MASTER      = "warp_table_master.csv"
	WM_WRITE_FILE_IF_MISSING = true
)

func NewWarpManager(game *Game) *WarpManager {
	wm := &WarpManager{}
	wm.warpDestList = []*WarpDest{}
	//wm.warpDestList = append(wm.warpDestList, &WarpDest{0, 1, 3, 3})
	//wm.warpDestList = append(wm.warpDestList, &WarpDest{1, 2, 3, 3})
	//wm.warpDestList = append(wm.warpDestList, &WarpDest{2, 3, 3, 3})
	//wm.warpDestList = append(wm.warpDestList, &WarpDest{3, 4, 3, 3})
	//wm.warpDestList = append(wm.warpDestList, &WarpDest{4, 5, 3, 3})
	wm.loadDataFromFile()
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

func (wm *WarpManager) loadDataFromFile() error {
	fmt.Println("warp data load")
	wm.warpDestList = []*WarpDest{}
	//writeMapToFile(wm.tileData, wm_DEFAULT_MAP_FILENAME)
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
		//gridX, gridY, moveGridX, moveGridY, kind
		pfTemp = &WarpDest{data[0], data[1], data[2], data[3]}
		wm.warpDestList = append(wm.warpDestList, pfTemp)

	}
	return nil
}
