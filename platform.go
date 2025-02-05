package main

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	PF_DEFAULT_WIDTH    = 150
	PF_DEFAULT_HEIGHT   = 30
	PF_DEFAULT_SPEED    = 2
	PF_MOTION_TICK_RATE = 4
	PF_TICK_COUNT_MAX   = 4
	PF_IMAGE_FILENAME   = "longBrick.png"
	PF_FILENAME_BASE    = "platform"
	PF_FILENAME_END     = ".csv"
	PF_MAX_PLATFORMS    = 10
	PF_MOD_RECT_H       = 3
)

type Platform struct {
	startX    int
	startY    int
	moveGridX int
	moveGridY int
	width     int
	height    int
	endX      int
	endY      int
	currX     int
	currY     int
	gridX     int
	gridY     int
	visable   bool
	active    bool
	speed     int
	kind      int
}

func (pf *Platform) getRect() *rect {
	r := &rect{}
	r.x = pf.currX
	r.y = pf.currY
	r.height = pf.height
	r.width = pf.width
	return r
}

func NewPlatform(gridX, gridY, moveGridX, moveGridY, kind int) *Platform {
	pf := &Platform{}
	pf.gridX = gridX
	pf.gridY = gridY
	pf.moveGridX = moveGridX
	pf.moveGridY = moveGridY
	pf.startX = gridX * GAME_TILE_SIZE
	pf.startY = gridY * GAME_TILE_SIZE
	pf.width = PF_DEFAULT_WIDTH
	pf.height = PF_DEFAULT_HEIGHT
	pf.endX = (pf.gridX + moveGridX) * GAME_TILE_SIZE
	pf.endY = (pf.gridY + moveGridY) * GAME_TILE_SIZE
	pf.currX = pf.startX
	pf.currY = pf.startY
	return pf
}

type PlatformManager struct {
	game           *Game
	currentTick    int
	platformsArray [PF_MAX_PLATFORMS]*Platform

	defImage                 *ebiten.Image
	images                   []*ebiten.Image
	testRect                 rect
	filename_base            string
	assetID                  int
	playerStandingOnPlatform bool
	modifiedRect             *rect
}

func NewPlatformManager(game *Game) *PlatformManager {
	pfm := &PlatformManager{}
	pfm.filename_base = PF_FILENAME_BASE
	pfm.game = game
	pfm.platformsArray = [PF_MAX_PLATFORMS]*Platform{}
	pfm.initImages()
	pfm.assetID = 0
	//pfm.AddInstanceToGrid(3, 3, 0)
	pfm.modifiedRect = &rect{}
	pfm.modifiedRect.height = PF_MOD_RECT_H
	pfm.modifiedRect.width = playerWidth
	pfm.loadDataFromFile()
	return pfm

}

func (pfm *PlatformManager) initImages() {
	var err error
	//cwd, _ := os.Getwd()
	imageDir := path.Join(imageSubdir, PF_IMAGE_FILENAME)
	//fmt.Println(imageDir)
	var rawImage, stretchedImage *ebiten.Image
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	stretchedImage = copyAndStretchImage(rawImage, PF_DEFAULT_WIDTH, PF_DEFAULT_HEIGHT)

	if err != nil {
		log.Fatal(err)
	}
	pfm.defImage = stretchedImage
	pfm.images = []*ebiten.Image{}
	pfm.images = append(pfm.images, pfm.defImage)

}

func (pfm *PlatformManager) Draw(screen *ebiten.Image) {

	for _, v := range pfm.platformsArray {
		if nil != v {
			screenX := (v.currX) - worldOffsetX
			screenY := (v.currY) - worldOffsetY
			DrawImageAt(screen, pfm.images[v.kind], screenX, screenY)
		}
	}

}

func (pfm *PlatformManager) Update() {
	pfm.checkPlatformTouchedPlayer()
	pfm.platformMotion()

}

func (pfm *PlatformManager) platformMotion() {

}

func (pfm *PlatformManager) checkPlatformTouchedPlayer() {

	playerRect := pfm.game.player.getWorldColliderRect()
	rectTop := playerRect.y + playerRect.height - 3
	pfm.modifiedRect.x = playerRect.x
	pfm.modifiedRect.y = rectTop
	for _, v := range pfm.platformsArray {
		if nil != v {
			pfm.testRect.height = v.height
			pfm.testRect.width = v.width
			pfm.testRect.x = v.currX
			pfm.testRect.y = v.currY

			if collideRect(*pfm.modifiedRect, pfm.testRect) {
				//fmt.Println("player touched platform")
				pfm.playerStandingOnPlatform = true
			} else {
				pfm.playerStandingOnPlatform = false
			}
		}
	}
}

func (pfm *PlatformManager) saveDataToFile() {
	name := pfm.getDataFileURL()
	numericData := [][]int{}
	rows := len(pfm.platformsArray)
	for i := 0; i < rows; i++ {
		PlatformObj := pfm.platformsArray[i]
		if PlatformObj != nil {
			//gridX, gridY, moveGridX, moveGridY, kind
			record := []int{
				PlatformObj.gridX,
				PlatformObj.gridY,
				PlatformObj.moveGridX,
				PlatformObj.moveGridY,
				PlatformObj.kind}
			numericData = append(numericData, record)
		}

	}
	if rows != 0 {

		write2DIntListToFile(numericData, name)
	} else {
		log.Println("Platforms: no data to write, ", name)
	}
}

func (pfm *PlatformManager) getDataFileURL() string {
	filename := pfm.filename_base + strconv.Itoa(pfm.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}

func (pfm *PlatformManager) loadDataFromFile() error {
	//fmt.Println("platform load")
	//writeMapToFile(pfm.tileData, pfm_DEFAULT_MAP_FILENAME)
	name := pfm.getDataFileURL()
	numericData, err := loadDataListFromFile(name)
	rows := len(numericData)
	if rows == 0 {
		pfm.platformsArray = [PF_MAX_PLATFORMS]*Platform{}
		return nil
	}
	if err != nil {
		pfm.platformsArray = [PF_MAX_PLATFORMS]*Platform{}
		return err
	}
	pfm.platformsArray = [PF_MAX_PLATFORMS]*Platform{}
	for i := 0; i < PF_MAX_PLATFORMS && i < rows; i++ {
		data := numericData[i]
		pfTemp := &Platform{}
		//gridX, gridY, moveGridX, moveGridY, kind
		pfTemp = NewPlatform(data[0], data[1], data[2], data[3], data[4])
		pfm.platformsArray[i] = pfTemp

	}
	return nil
}

func (pfm *PlatformManager) validatAssetID(kind int) bool {
	if kind < len(pfm.images) && kind > -1 {
		return true
	} else {
		return false
	}

}

func (pfm *PlatformManager) CycleAssetKind(direction int) {
	propAssetID := pfm.assetID + direction
	isValid := pfm.validatAssetID(propAssetID)
	if isValid {
		pfm.assetID = propAssetID
	}

	fmt.Println("Selected Platform ", pfm.assetID)

}

func (pfm *PlatformManager) inactiveSlot() int {
	// find usable slot in pickups array, or -1 if there is none
	for i := 0; i < len(pfm.platformsArray); i++ {
		if nil == pfm.platformsArray[i] {
			return i
		}
	}
	return -1
}

func (pfm *PlatformManager) AddInstanceToGrid(gridX, gridY, kind int) {
	emptySlot := pfm.inactiveSlot()
	if emptySlot != -1 {

		pu := NewPlatform(gridX, gridY, 0, 0, kind)
		pfm.platformsArray[emptySlot] = pu
		log.Println("Added Platform ", kind)
	} else {
		log.Println("Failed to add Platform, no open slots")
	}
}
