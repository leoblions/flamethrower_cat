package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Hud struct {
	game *Game
	//components
	scoreString   *RasterString
	actString     *RasterString
	livesString   *RasterString
	levelString   *RasterString
	centerText    *RasterString
	t1Text        *RasterString
	centerMarquee *Marquee
}

func NewHud(game *Game) *Hud {
	hud := &Hud{}
	hud.game = game
	hud.initRasterStrings()
	return hud
}

func (hud *Hud) Update() {
	hud.updateRasterStrings()
	hud.centerMarquee.Update()

}
func (hud *Hud) Draw(screen *ebiten.Image) {
	hud.drawRasterStrings(screen)
	hud.centerMarquee.Draw(screen)

}

func (hud *Hud) drawRasterStrings(screen *ebiten.Image) {
	hud.scoreString.Draw(screen)
	hud.actString.Draw(screen)
	hud.centerText.Draw(screen)
	hud.levelString.Draw(screen)
	hud.livesString.Draw(screen)
	hud.t1Text.Draw(screen)
}
func (hud *Hud) updateRasterStrings() {
	hud.scoreString.Update()
	hud.centerText.Update()
	hud.levelString.Update()
	hud.livesString.Update()
	hud.t1Text.Update()
}

func (hud *Hud) initRasterStrings() {
	g := hud.game
	hud.scoreString = NewRasterString(g, "Score: 0", 10, 10)
	hud.centerMarquee = NewMarquee(g, 0, 200, "Cooking with gas!")
	hud.centerText = NewRasterString(g, "You won", centerTextX, centerTextY)
	hud.centerText.visible = false
	levelStr := fmt.Sprintf("Level: %d", g.level)
	hud.levelString = NewRasterString(g, levelStr, g.screenWidth-90, 10)
	hud.livesString = NewRasterString(g, "Lives: 3", (g.screenWidth/2)-50, 10)
	hud.actString = NewRasterString(g, "Press E", (g.screenWidth)-70, (g.screenHeight)-50)
	hud.actString.setBackgroundColor(&color.RGBA{0xff, 0x00, 0x00, 0xff})
	hud.t1Text = NewRasterTitleString(g, "dead cat", 200, 200)
	hud.t1Text.visible = false

}
