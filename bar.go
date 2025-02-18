package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	BAR_DEFAULT_FG = &color.RGBA{0xff, 0x0f, 0x0f, 0x90}
	BAR_DEFAULT_BG = &color.RGBA{0x0f, 0x0f, 0x0f, 0x90}
	BAR_DEFAULT_SH = &color.RGBA{0x30, 0x10, 0x10, 0x90}
)

type Bar struct {
	rect
	fillX int
	//fillXlast int
	colorFG *color.RGBA
	colorBG *color.RGBA
	shadow  *color.RGBA
	image   *ebiten.Image
}

func NewBar(screenX, screenY, width, height int) *Bar {
	bar := &Bar{}
	bar.x = screenX
	bar.y = screenY
	bar.width = width
	bar.height = height
	bar.colorFG = BAR_DEFAULT_FG
	bar.colorBG = BAR_DEFAULT_BG
	bar.shadow = BAR_DEFAULT_SH
	bar.fillX = bar.width
	//bar.fillX = 50
	bar.updateImage()

	return bar

}

func (bar *Bar) updateImage() {
	// background
	barImage := ebiten.NewImage(bar.width, bar.height)
	barImage.Fill(bar.colorBG)
	// filled in part
	filledPart := ebiten.NewImage(bar.fillX, bar.height)
	filledPart.Fill(bar.colorFG)
	op := &ebiten.DrawImageOptions{}
	barImage.DrawImage(filledPart, op)
	// shadow
	shadowTop := bar.height * 2 / 3
	shadow := ebiten.NewImage(bar.fillX, shadowTop)
	shadow.Fill(bar.shadow)
	op.GeoM.Translate(float64(0), float64(shadowTop))
	barImage.DrawImage(shadow, op)
	bar.image = barImage

}

func (bar *Bar) Update() {

}

func (bar *Bar) UpdateFilledAreaX(filledX int) {
	filledXclamped := clamp(1, bar.width, filledX)
	bar.fillX = filledXclamped
	bar.updateImage()

}
func (bar *Bar) FillPercent(percent int) {
	pillPct := clamp(0, 100, percent)
	fillX := int(float32(bar.width) * (float32(pillPct) / float32(100)))
	bar.UpdateFilledAreaX(fillX)

}
func (bar *Bar) Draw(screen *ebiten.Image) {
	DrawImageAt(screen, bar.image, bar.x, bar.y)

}
