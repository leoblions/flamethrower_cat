package main

import (
	"fmt"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	imageSubdir      = "images"
	imageBall        = "ball.png"
	defaultSpeedBall = 4
	lowSpeedBall     = 2
	ballHeight       = 30
	ballWidth        = 30
	xMaxBall         = screenWidth - ballWidth
)

type Ball struct {
	game         *Game
	visible      bool
	stopped      bool
	frozen       bool
	screenX      int
	screenY      int
	velX         int
	velY         int
	currImage    *ebiten.Image
	followPaddle bool
}

func (p *Ball) getColliderRect() rect {
	return rect{p.screenX, p.screenY, ballWidth, ballHeight}
}

func (p *Ball) midpoint() (int, int) {
	return (p.screenX + (ballWidth / 2)), p.screenY + (ballHeight / 2)
}

func (p *Ball) init(g *Game, startX, startY, velX, velY int) {
	p.game = g
	var err error
	//cwd, _ := os.Getwd()
	imageDir := path.Join(imageSubdir, imageBall)
	//fmt.Println(imageDir)
	var rawImage, stretchedImage *ebiten.Image
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	stretchedImage = copyAndStretchImage(rawImage, ballWidth, ballHeight)

	if err != nil {
		log.Fatal(err)
	}
	p.currImage = stretchedImage
	_ = stretchedImage
	p.screenX = startX
	p.screenY = startY
	p.velX = velX
	p.velY = velY
	p.stopped = true
	p.frozen = true
	p.visible = false
	p.followPaddle = true

	//p.velX = 2
}

func (p *Ball) reset() {
	p.velX = ballInitVelX
	p.velY = ballInitVelY
	p.stopped = true
	p.screenX = startBallX
	p.screenY = startBallY
	p.visible = true
	p.followPaddle = true
}

func (p *Ball) stop() {
	p.velX = 0
	p.velY = 0
	p.stopped = true
	//p.screenX = startBallX
	//p.screenY = startBallY
}

func (p *Ball) Draw(screen *ebiten.Image) {
	if p.visible {
		DrawImageAt(screen, p.currImage, p.screenX, p.screenY)
	}

}

func (p *Ball) ballPlayerCollision() bool {
	playerColl := p.game.player.getColliderRect()
	ballColl := p.getColliderRect()
	return collideRect(playerColl, ballColl)

}

func (p *Ball) outOfBounds() bool {
	if (p.screenY + ballHeight) > screenHeight {
		return true
	} else {
		return false
	}
}

func (p *Ball) Serve() {
	p.stopped = false
	p.followPaddle = false
}

func (p *Ball) Update() {

	if p.screenX < 0 || p.screenX > xMaxBall {
		p.velX *= -1
	}
	if p.screenY < 0 {
		p.velY *= -1
	}

	if p.ballPlayerCollision() && p.screenY < p.game.player.screenY {
		p.velY = -abs(p.velY)
		p.velX = -abs(p.velX)

		playerMidpoint := p.game.player.midpointX()
		ballMPX, _ := p.midpoint()
		var speed = 0
		if abs(ballMPX-playerMidpoint) < ballWidth {
			speed = lowSpeedBall
		} else {
			speed = defaultSpeedBall
		}
		if mpx, _ := p.midpoint(); mpx > playerMidpoint {
			p.velX = abs(speed)
		} else {
			p.velX = -abs(speed)
		}
	}
	if !p.stopped && !p.frozen {
		p.screenX += p.velX
		p.screenY += p.velY

	}

	if p.outOfBounds() {
		p.reset()
		p.game.lives -= 1
		p.game.livesString.stringContent = fmt.Sprintf("Lives: %d", p.game.lives)
	}

	//p.screenY += p.velY

}
