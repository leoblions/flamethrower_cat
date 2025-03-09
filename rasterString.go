package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// var rasterStringStatic = struct {
// 	ssCharsWide int
// 	ssCharsTall int

// 	ssimageSubdir              string
// 	lettersSpritesheetFilename string
// 	lettersCharmapFile         string
// 	runeImageMap               map[rune]*ebiten.Image
// }{
// 	ssCharsWide: 10,
// 	ssCharsTall: 10,

// 	ssimageSubdir:              "images",
// 	lettersSpritesheetFilename: "letterSpritesW.png",
// 	lettersCharmapFile:         "charmap_letters.txt",
// 	runeImageMap:               nil,
// }

var (
	rast_ss_CharsWide int = 10
	ssCharsTall       int = 10
	// resource files
	rast_ss_ssimageSubdir              string = "images"
	rast_ss_lettersSpritesheetFilename string = "letterSpritesW.png"
	//rast_ss_data_subdir                string                 = "leveldata"
	rast_ss_lettersCharmapFile      string                 = "charmap_letters.txt"
	rast_ss_lettersCharmapFileTitle string                 = "charmap_letters_title.txt"
	rast_ss_runeImageMap            map[rune]*ebiten.Image = nil
	rast_ss_runeImageMapTitle       map[rune]*ebiten.Image = nil
)

const (
	ssimageSubdir              = "images"
	lettersSpritesheetFilename = "letterSpritesT.png"
	characterOrdinalOffset     = 48
)
const (
	RAST_SS_CHAR_PIX_H           = 15
	RAST_SS_CHAR_PIX_W           = 10
	RAST_SS_CHAR_PIX_TITLE_H     = 25
	RAST_SS_CHAR_PIX_TITLE_W     = 25
	RAST_SS_CHAR_KERNING         = 10
	RAST_SS_CHAR_KERNING_TITLE   = 25
	RAST_SS_BLINK_RATE           = 30
	RAST_SS_BG_FILL              = true
	RAST_SS_BG_ALPHA             = 0x8f
	RAST_SS_IMAGE_FILENAME       = "letterSpritesT.png"
	RAST_SS_IMAGE_TITLE_FILENAME = "titleLetter.png"
	IMAGES_SUBDIR                = "images"
	LETTERS_LOCATION_FILE        = "letters.json"
	RAST_SS_DATA_SUBDIR          = "leveldata"
)

type RasterString struct {
	game              *Game
	stringContent     string
	screenX           int
	screenY           int
	spritesheet       *ebiten.Image
	letterSprites     []*ebiten.Image
	visible           bool
	blink             bool
	currImage         *ebiten.Image
	stringImage       *ebiten.Image
	runeImageMap      map[rune]*ebiten.Image
	runeImageMapTitle map[rune]*ebiten.Image
	letterHeight      int
	letterWidth       int
	letterKerning     int
	blinkCounter      func() bool
	backgroundColor   *color.RGBA
	//updateImageMethod func() *ebiten.Image
}

type characterRecord struct {
	Col    int
	Row    int
	Letter rune
}

var (
	bg_color = color.RGBA{0x0f, 0x0f, 0x0f, RAST_SS_BG_ALPHA}
)

func (p *RasterString) getColliderRect() rect {
	return rect{p.screenX, p.screenY, PL_COLLRECT_W, PL_COLLRECT_H}
}

func (p *RasterString) midpointX() int {
	return p.screenX + (PL_COLLRECT_W / 2)

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
	p.blinkCounter = counterClosureTF(RAST_SS_BLINK_RATE)
	//p.updateImageMethod = p.getRasterStringAsSingleImage
	//srite sheet
	var err error
	spriteSheetFullPath := path.Join(rast_ss_ssimageSubdir, lettersSpritesheetFilename)
	charmapFilePath := path.Join(RAST_SS_DATA_SUBDIR, rast_ss_lettersCharmapFile)
	p.visible = true

	if err != nil {
		log.Fatal(nil)
	}

	if rast_ss_runeImageMap == nil {
		p.runeImageMap = p.initializeLetterSprites(spriteSheetFullPath, charmapFilePath)
		rast_ss_runeImageMap = p.runeImageMap
	} else {
		p.runeImageMap = rast_ss_runeImageMap
	}
	p.backgroundColor = &bg_color
	//p.runeImageMap = p.initializeLetterSprites(spriteSheetFullPath)
	p.stringImage = p.getRasterStringAsSingleImage()

	return p

}

func NewRasterTitleString(g *Game, content string, startX, startY int) *RasterString {
	p := &RasterString{}
	p.game = g
	p.stringContent = content
	p.screenX = startX
	p.screenY = startY
	p.letterHeight = RAST_SS_CHAR_PIX_TITLE_H
	p.letterWidth = RAST_SS_CHAR_PIX_TITLE_W
	p.letterKerning = RAST_SS_CHAR_KERNING_TITLE
	p.blinkCounter = counterClosureTF(RAST_SS_BLINK_RATE)
	//p.blink = true
	//srite sheet
	var err error
	spriteSheetFullPath := path.Join(rast_ss_ssimageSubdir, RAST_SS_IMAGE_TITLE_FILENAME)
	charmapFilePath := path.Join(RAST_SS_DATA_SUBDIR, rast_ss_lettersCharmapFileTitle)
	p.visible = true

	if err != nil {
		log.Fatal(nil)
	}

	if rast_ss_runeImageMapTitle == nil {
		p.runeImageMap = p.initializeLetterSprites(spriteSheetFullPath, charmapFilePath)
		rast_ss_runeImageMapTitle = p.runeImageMap
	} else {
		p.runeImageMap = rast_ss_runeImageMapTitle
	}
	p.backgroundColor = &bg_color
	//p.runeImageMap = p.initializeLetterSprites(spriteSheetFullPath)
	p.stringImage = p.getRasterStringAsSingleImage()

	return p

}

func (p *RasterString) setBackgroundColor(newColor *color.RGBA) {
	p.backgroundColor = newColor
	p.stringImage = p.getRasterStringAsSingleImage()
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

	p.stringImage = p.getRasterStringAsSingleImage()

}

func (p *RasterString) Draw(screen *ebiten.Image) {
	if !p.visible {
		return
	}

	DrawImageAt(screen, p.stringImage, p.screenX, p.screenY)

}

func (p *RasterString) getRasterStringAsSingleImage() *ebiten.Image {
	width := len(p.stringContent) * p.letterWidth
	if width < 1 {
		width = 1
	}
	stringImage := ebiten.NewImage(width, p.letterHeight)
	xOffsetTotal := 0
	if RAST_SS_BG_FILL {
		stringImage.Fill(*p.backgroundColor)
	}
	for _, letter := range p.stringContent {
		letterImage := p.runeImageMap[letter] // get image from map using rune as key

		if letterImage != nil {

			DrawImageAt(stringImage, letterImage, xOffsetTotal, 0)
			xOffsetTotal += p.letterKerning
		} else if letter == ' ' {
			xOffsetTotal += p.letterKerning
		}

	}

	return stringImage
}

func (p *RasterString) Update() {
	if p.blink && p.blinkCounter() {
		p.visible = !p.visible
	}

}

func (p *RasterString) readSpriteLocationTableFile(charmapFilePath string) []*characterRecord {
	linesList := getListOfLinesFromFile(charmapFilePath)
	//p.exportStringSliceToJSON(linesList)
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

func (p *RasterString) encodeCRlist(crList []*characterRecord) error {
	var err error
	crListV := []characterRecord{}
	for _, v := range crList {
		crListV = append(crListV, *v)
	}
	// serialize to json
	dataFilePath := path.Join(GAME_LEVEL_DATA_DIR, LETTERS_LOCATION_FILE)
	//var marshalledData []byte
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(crListV)

	// write file
	file, err := os.Create(dataFilePath)
	writer := bufio.NewWriter(file)
	defer file.Close()
	//mdList := []string{}
	//fmt.Println(" serialized bytes ", buffer.Bytes())

	for b := buffer.Next(10); err == io.EOF; {

		_, err = writer.Write(b)
		log.Println(b)

	}

	//marshalledData, err = json.Marshal(mdList)

	//fmt.Println(marshalledData)

	if err != nil {
		log.Fatal(err)
		log.Fatal(fmt.Sprintf("charmap failed to write file: %s/n", dataFilePath))
		return err
	}

	//write data to file

	return err

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func (p *RasterString) encodeCRlistP(crList []*characterRecord) error {
	var err error

	// serialize to json
	dataFilePath := path.Join(GAME_LEVEL_DATA_DIR, LETTERS_LOCATION_FILE)

	jsonData, err := json.Marshal(crList)
	check(err)

	//fmt.Println(string(jsonData))

	file, err := os.Create(dataFilePath)

	defer file.Close()
	_, err = file.Write(jsonData)
	file.Sync()

	if err != nil {
		log.Fatal(err)
		log.Fatal(fmt.Sprintf("charmap failed to write file: %s/n", dataFilePath))
		return err
	}

	//write data to file

	return err

}

func (p *RasterString) initializeLetterSprites(spriteSheetFullPath, charmapFilePath string) map[rune]*ebiten.Image {
	recordsList := p.readSpriteLocationTableFile(charmapFilePath)
	var err error

	if !fileExists(spriteSheetFullPath) {
		log.Fatal("initializeLetterSprites can't find the file ", spriteSheetFullPath)
	}
	if !fileExists(charmapFilePath) {
		log.Fatal("initializeLetterSprites can't find the file ", charmapFilePath)
	}

	//imageDir := path.Join(IMAGES_SUBDIR, RAST_SS_IMAGE_FILENAME)
	p.spritesheet, _, err = ebitenutil.NewImageFromFile(spriteSheetFullPath)
	rawImage, _, err := ebitenutil.NewImageFromFile(spriteSheetFullPath)
	if err != nil {
		log.Fatal(err)
	}

	//letterSpriteList := []*ebiten.Image{}
	runeImageMap := map[rune]*ebiten.Image{}
	for _, record := range recordsList {
		x := record.Col * p.letterWidth
		y := record.Row * p.letterHeight
		w := p.letterWidth
		h := p.letterHeight
		thisRune := record.Letter
		letterImage := getSubImage(rawImage, x, y, w, h)
		runeImageMap[thisRune] = letterImage
	}
	return runeImageMap
}
