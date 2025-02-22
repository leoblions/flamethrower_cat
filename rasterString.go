package main

import (
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var rasterStringStatic = struct {
	ssCharsWide int
	ssCharsTall int

	ssimageSubdir              string
	lettersSpritesheetFilename string
	lettersCharmapFile         string
	runeImageMap               map[rune]*ebiten.Image
}{
	ssCharsWide: 10,
	ssCharsTall: 10,

	ssimageSubdir:              "images",
	lettersSpritesheetFilename: "letterSpritesW.png",
	lettersCharmapFile:         "charmap_letters.txt",
	runeImageMap:               nil,
}

const (
	ssimageSubdir              = "images"
	lettersSpritesheetFilename = "letterSpritesT.png"
	characterOrdinalOffset     = 48
)
const (
	RAST_SS_CHAR_PIX_H     = 15
	RAST_SS_CHAR_PIX_W     = 10
	RAST_SS_CHAR_KERNING   = 10
	RAST_SS_IMAGE_FILENAME = "letterSpritesT.png"
	IMAGES_SUBDIR          = "images"
)

type RasterString struct {
	game          *Game
	stringContent string
	screenX       int
	screenY       int
	spritesheet   *ebiten.Image
	letterSprites []*ebiten.Image
	visible       bool
	currImage     *ebiten.Image
	runeImageMap  map[rune]*ebiten.Image
	letterHeight  int
	letterWidth   int
	letterKerning int
}

type characterRecord struct {
	col    int
	row    int
	letter rune
}

func (p *RasterString) getColliderRect() rect {
	return rect{p.screenX, p.screenY, playerWidth, playerHeight}
}

func (p *RasterString) midpointX() int {
	return p.screenX + (playerWidth / 2)

}

func NewRasterString(g *Game, content string, startX, startY int) *RasterString {
	p := &RasterString{}
	p.game = g
	p.stringContent = content
	p.screenX = startX
	p.screenY = startY
	p.letterHeight = RAST_SS_CHAR_PIX_H
	p.letterWidth = RAST_SS_CHAR_PIX_W
	p.letterKerning = RAST_SS_CHAR_KERNING
	//srite sheet
	var err error
	imageDir := path.Join(ssimageSubdir, lettersSpritesheetFilename)
	p.spritesheet, _, err = ebitenutil.NewImageFromFile(imageDir)
	p.currImage = getSubImage(p.spritesheet, 0, 0, p.letterWidth, p.letterHeight)
	p.visible = true

	if err != nil {
		log.Fatal(nil)
	}

	if rasterStringStatic.runeImageMap == nil {
		p.runeImageMap = p.initializeLetterSprites()
		rasterStringStatic.runeImageMap = p.runeImageMap
	} else {
		p.runeImageMap = rasterStringStatic.runeImageMap
	}
	p.runeImageMap = p.initializeLetterSprites()
	return p

}

func (p *RasterString) updateText(content string) {

	p.stringContent = content

	var err error
	imageDir := path.Join(ssimageSubdir, lettersSpritesheetFilename)
	p.spritesheet, _, err = ebitenutil.NewImageFromFile(imageDir)
	p.currImage = getSubImage(p.spritesheet, 0, 0, p.letterWidth, p.letterHeight)
	p.visible = true

	if err != nil {
		log.Fatal(nil)
	}

}

func (p *RasterString) Draw(screen *ebiten.Image) {
	//runeList := []rune(p.stringContent)
	//DrawImageAt(screen, p.currImage, p.screenX, p.screenY)
	if !p.visible {
		return
	}
	xOffsetTotal := 0
	for _, letter := range p.stringContent {
		letterImage := p.runeImageMap[letter]

		if letterImage != nil {
			//fmt.Println("draw ", c)
			DrawImageAt(screen, letterImage, p.screenX+xOffsetTotal, p.screenY)
			xOffsetTotal += p.letterKerning
		} else if letter == ' ' {
			xOffsetTotal += p.letterKerning
		}

	}

}

func (p *RasterString) getRasterStringAsSingleImage() *ebiten.Image {
	width := len(p.stringContent) * p.letterWidth
	stringImage := ebiten.NewImage(width, p.letterHeight)
	xOffsetTotal := 0
	for _, letter := range p.stringContent {
		letterImage := p.runeImageMap[letter]

		if letterImage != nil {
			//fmt.Println("draw ", c)
			DrawImageAt(stringImage, letterImage, xOffsetTotal, 0)
			xOffsetTotal += p.letterKerning
		} else if letter == ' ' {
			xOffsetTotal += p.letterKerning
		}

	}

	return stringImage
}

func (p *RasterString) Update() {

}

func (p *RasterString) readSpriteLocationTableFile() []*characterRecord {
	linesList := getListOfLinesFromFile(rasterStringStatic.lettersCharmapFile)
	recordsList := []*characterRecord{}
	for _, line := range linesList {
		runeSlice := []rune(*line)
		row := int(runeSlice[1]) - characterOrdinalOffset
		col := int(runeSlice[0]) - characterOrdinalOffset
		cr := characterRecord{col, row, runeSlice[2]}
		//fmt.Println(cr)
		recordsList = append(recordsList, &cr)
	}

	return recordsList

}

func (p *RasterString) initializeLetterSprites() map[rune]*ebiten.Image {
	recordsList := p.readSpriteLocationTableFile()
	imageDir := path.Join(IMAGES_SUBDIR, RAST_SS_IMAGE_FILENAME)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	if err != nil {
		log.Fatal(err)
	}
	//letterSpriteList := []*ebiten.Image{}
	runeImageMap := map[rune]*ebiten.Image{}
	for _, record := range recordsList {
		x := record.col * p.letterWidth
		y := record.row * p.letterHeight
		w := p.letterWidth
		h := p.letterHeight
		thisRune := record.letter
		letterImage := getSubImage(rawImage, x, y, w, h)
		runeImageMap[thisRune] = letterImage
	}
	return runeImageMap
}
