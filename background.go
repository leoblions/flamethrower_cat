package main

import (
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	BGI_MAX_PICKUPS = 10
	BGI_CLOUDS      = "backgroundClouds.png"
	BGI_CAVE        = "backgroundCaves.png"
	BGI_ISLAND      = "backgroundIsland.png"
	BGI_LIGHTHOUSE  = "backgroundLighthouse.png"
	BGI_IMAGE_DIR   = "images"

	BGI_PICKUP_H = 50
	BGI_PICKUP_W = 50
)

type Background struct {
	game   *Game
	images []*ebiten.Image

	filename_base string
	assetID       int
	currBGIIndex  int
}

func NewBackground(game *Game) *Background {

	BGI := &Background{}
	BGI.game = game
	BGI.assetID = 0
	BGI.initImages()
	BGI.currBGIIndex = 0

	return BGI
}

func (BGI *Background) Draw(screen *ebiten.Image) {
	screenX := 0
	screenY := 0
	DrawImageAt(screen, BGI.images[BGI.currBGIIndex], screenX, screenY)

}

func (BGI *Background) Update() {
	switch BGI.game.level {
	case 3:
		//cave
		BGI.currBGIIndex = 1
	case 12, 13, 14, 15:
		//island
		BGI.currBGIIndex = 2
	case 19, 20, 21:
		//lighthouse
		BGI.currBGIIndex = 3
	default:
		//cloud
		BGI.currBGIIndex = 0

	}

}

func (BGI *Background) initImages() error {
	BGI.images = []*ebiten.Image{}
	var tempImg *ebiten.Image

	// clouds
	imageDir := path.Join(BGI_IMAGE_DIR, BGI_CLOUDS)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	BGI.images = append(BGI.images, rawImage)

	// cave
	imageDir = path.Join(BGI_IMAGE_DIR, BGI_CAVE)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	BGI.images = append(BGI.images, rawImage)

	// island
	imageDir = path.Join(BGI_IMAGE_DIR, BGI_ISLAND)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	BGI.images = append(BGI.images, rawImage)

	// lighthouse
	imageDir = path.Join(BGI_IMAGE_DIR, BGI_LIGHTHOUSE)
	tempImg, _, err = ebitenutil.NewImageFromFile(imageDir)
	//tempImg = copyAndStretchImage(rawImage, BGI_PICKUP_W, BGI_PICKUP_H)
	BGI.images = append(BGI.images, tempImg)

	imageDir = path.Join(BGI_IMAGE_DIR, BGI_CLOUDS)
	tempImg, _, err = ebitenutil.NewImageFromFile(imageDir)
	//tempImg = copyAndStretchImage(rawImage, BGI_PICKUP_W, BGI_PICKUP_H)
	BGI.images = append(BGI.images, tempImg)
	return err

}
