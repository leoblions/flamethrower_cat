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
	RAST_SS_BG_FILL        = true
	RAST_SS_BG_ALPHA       = 0x8f
	RAST_SS_IMAGE_FILENAME = "letterSpritesT.png"
	IMAGES_SUBDIR          = "images"
	LETTERS_LOCATION_FILE  = "letters.json"
)

type RasterString struct {
	game            *Game
	stringContent   string
	screenX         int
	screenY         int
	spritesheet     *ebiten.Image
	letterSprites   []*ebiten.Image
	visible         bool
	currImage       *ebiten.Image
	stringImage     *ebiten.Image
	runeImageMap    map[rune]*ebiten.Image
	letterHeight    int
	letterWidth     int
	letterKerning   int
	backgroundColor *color.RGBA
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
	p.backgroundColor = &bg_color
	p.runeImageMap = p.initializeLetterSprites()
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
		letterImage := p.runeImageMap[letter]

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

}

func (p *RasterString) readSpriteLocationTableFile() []*characterRecord {
	linesList := getListOfLinesFromFile(rasterStringStatic.lettersCharmapFile)
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

func (p *RasterString) initializeLetterSprites() map[rune]*ebiten.Image {
	recordsList := p.readSpriteLocationTableFile()
	//p.encodeCRlist(recordsList)
	//p.encodeCRlistP(recordsList)

	imageDir := path.Join(IMAGES_SUBDIR, RAST_SS_IMAGE_FILENAME)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
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
