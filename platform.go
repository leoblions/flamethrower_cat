package main

const (
	PF_DEFAULT_WIDTH    = 150
	PF_DEFAULT_HEIGHT   = 30
	PF_DEFAULT_SPEED    = 2
	PF_MOTION_TICK_RATE = 4
	PF_TICK_COUNT_MAX   = 4
	PF_IMAGE_FILENAME   = "longBrick.png"
)

type Platform struct {
	startX int
	startY int
	width  int
	height int
	endX   int
	endY   int
	currX  int
	currY  int

	speed int
	kind  int
}

func (pf *Platform) getRect() *rect {
	r := &rect{}
	r.x = pf.currX
	r.y = pf.currY
	r.height = pf.height
	r.width = pf.width
	return r
}

func NewPlatform(startX, startY, kind int) *Platform {
	pf := &Platform{}
	pf.startX = startX
	pf.startY = startY
	pf.width = PF_DEFAULT_WIDTH
	pf.height = PF_DEFAULT_HEIGHT
	pf.endX = startX
	pf.endY = startY
	pf.currX = startX
	pf.currY = startY
	return pf
}

type PlatformManager struct {
	game                *Game
	currentTick         int
	currentPlatformList []*Platform
}

func NewPlatformManager(game *Game) *PlatformManager {
	pfm := &PlatformManager{}
	pfm.game = game
	pfm.currentPlatformList = []*Platform{}
	pfm.initImages()

	return pfm

}

func (pfm *PlatformManager) initImages() {

}
