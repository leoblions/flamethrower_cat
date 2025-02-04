package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Marquee struct {
	game                *Game
	screenX             int
	screenXF            float32
	screenY             int
	screenXMax          int
	stringContent       string
	midScreenDelay      bool
	visible             bool
	speed               float32
	midScreenDelayTicks int
	midScreenLocation   int
	midScreenCounter    int
	stringImage         *ebiten.Image
	centerPaused        bool
	centerPauseDone     bool
	centerTextOffset    int
}

func NewMarquee(game *Game, startX, startY int, text string) *Marquee {

	mq := &Marquee{}
	mq.game = game
	mq.screenX = startX
	mq.screenY = startY
	mq.screenXMax = game.screenWidth
	mq.stringContent = text
	mq.visible = true
	mq.screenXF = 0.0
	mq.speed = 5.5
	mq.midScreenDelay = true
	mq.midScreenCounter = 0
	mq.centerTextOffset = 25
	mq.midScreenLocation = (game.screenWidth / 2) - mq.centerTextOffset
	mq.midScreenDelayTicks = 60
	mq.centerPaused = false
	rasterTextTemp := NewRasterString(game, text, startX, startY)
	mq.stringImage = rasterTextTemp.getRasterStringAsSingleImage()
	return mq
}

func (p *Marquee) Draw(screen *ebiten.Image) {

	if !p.visible || nil == p.stringImage {
		return
	} else {
		DrawImageAt(screen, p.stringImage, p.screenX, p.screenY)
	}

}

func (p *Marquee) centerPause() {
	wrResult := withinRange(p.midScreenLocation, p.screenX, 10)
	//fmt.Println(wrResult)
	if !p.centerPauseDone && wrResult {

		if p.midScreenCounter < p.midScreenDelayTicks {
			p.midScreenCounter += 1
			p.centerPaused = true
		} else {
			p.centerPauseDone = true
			p.centerPaused = false
			p.midScreenCounter = 0
		}
	}
}

func (p *Marquee) Update() {

	if p.screenX < p.game.screenWidth {
		//fmt.Println(p.midScreenLocation)
		p.centerPause()
		//fmt.Println(p.centerPaused)
		if p.centerPaused {

		} else {
			p.screenXF += p.speed
			p.screenX = int(p.screenXF)
		}

	} else if !p.game.centerMarqueeEndActionsComplete {
		p.game.marqueeMessageComplete()
	}
	//fmt.Println(p.game.screenWidth)

}

func (p *Marquee) Start() {

	p.screenX = 0
	p.screenXF = 0.0
	p.visible = true
	p.centerPauseDone = false
	p.centerPaused = false
	p.game.centerMarqueeEndActionsComplete = false

}

func (p *Marquee) UpdateText(text string) {

	rasterTextTemp := NewRasterString(p.game, text, p.screenX, p.screenY)
	p.stringImage = rasterTextTemp.getRasterStringAsSingleImage()
	p.midScreenLocation = (p.game.screenWidth / 2) - (len(p.stringContent) * 10)

}
