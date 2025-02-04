package main

import (
	"fmt"
	"image/color"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	brickImage              = "longBrick.png"
	brickHeight             = 30
	brickWidth              = 100
	gridStartX              = 22
	gridStartY              = 50
	bricksX                 = 6
	bricksY                 = 6
	bricksTotal             = bricksX * bricksY
	brickTintMixLimit       = true
	brickTintMixLimitAmount = 0x55
)

var hexColors = []color.RGBA{
	{0xff, 0xff, 0xff, 0x3f},
	{0xff, 0xff, 0x0, 0x3f},
	{0xff, 0x0, 0xff, 0x3f},
	{0xff, 0x0, 0x0, 0x3f},
	{0x0, 0xff, 0xff, 0x3f},
	{0x0, 0xff, 0x0, 0x3f},
	{0x0, 0x0, 0xff, 0x3f},
	{0x0, 0x0, 0x0, 0x3f},
}

type BrickGrid struct {
	game               *Game
	worldX             int
	worldY             int
	screenX            int
	screenY            int
	brickGrid          [bricksTotal]*Brick
	currImage          *ebiten.Image
	coloredBrickImages []*ebiten.Image
}

type Brick struct {
	gridX   int
	gridY   int
	screenX int
	screenY int
	alive   bool
}

func (p *Brick) midpointX() int {
	return p.screenX + (brickWidth / 2)
}

func (p *Brick) midpoint() (int, int) {
	return (p.screenX + (brickWidth / 2)), p.screenY + (brickHeight / 2)
}

func (p *BrickGrid) getColliderRect() rect {
	return rect{p.screenX, p.screenY, playerWidth, playerHeight}
}

func (p *BrickGrid) midpointX() int {
	return p.screenX + (playerWidth / 2)

}

func (p *BrickGrid) Reset() {
	for _, brick := range p.brickGrid {
		brick.alive = true
	}

}

func (p *BrickGrid) checkPlayerWon() bool {
	aliveBricks := 0
	for _, v := range p.brickGrid {
		if v.alive {
			aliveBricks += 1
		}
	}
	return aliveBricks == 0

}

func (p *BrickGrid) KillAllBricks() {

	for _, v := range p.brickGrid {
		v.alive = false
	}

}

func (p *BrickGrid) init(g *Game) {
	p.game = g
	var err error
	//cwd, _ := os.Getwd()
	imageDir := path.Join(subdir, brickImage)
	//fmt.Println(imageDir)
	var rawImage *ebiten.Image
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	//stretchedImage = copyAndStretchImage(rawImage, playerWidth, playerHeight)

	if err != nil {
		log.Fatal(err)
	}
	p.currImage = rawImage
	//_ = stretchedImage

	p.brickGrid = p.createBricks()
	p.coloredBrickImages = p.createTintedBrickImages(rawImage)

	//p.velX = 2
}

func (p *BrickGrid) createBricks() [bricksTotal]*Brick {
	// returns an array of bricks
	bricksArray := [bricksTotal]*Brick{}
	arrIterator := 0
	for y := range bricksY {
		for x := range bricksX {
			screenX := x*brickWidth + gridStartX
			screenY := y*brickHeight + gridStartY
			bricksArray[arrIterator] = &Brick{x, y, screenX, screenY, true}
			arrIterator += 1
		}
	}
	return bricksArray

}

func (p *BrickGrid) createTintedBrickImages(orig *ebiten.Image) []*ebiten.Image {
	// returns an array of bricks
	imageList := []*ebiten.Image{}

	for _, color := range hexColors {
		if brickTintMixLimit {
			if color.A > brickTintMixLimitAmount {
				color.A = brickTintMixLimitAmount
			}
			if color.R > brickTintMixLimitAmount {
				color.R = brickTintMixLimitAmount
			}
			if color.G > brickTintMixLimitAmount {
				color.G = brickTintMixLimitAmount
			}
			if color.B > brickTintMixLimitAmount {
				color.B = brickTintMixLimitAmount
			}
		}
		coloredImage := copyAndTintImage(orig, &color)
		imageList = append(imageList, coloredImage)

	}
	return imageList

}

func copyAndTintImage(orig *ebiten.Image, color *color.RGBA) *ebiten.Image {

	oldX, oldY := orig.Bounds().Size().X, orig.Bounds().Size().Y
	tintImage := ebiten.NewImage(oldX, oldY)
	tintImage.Fill(color)
	outputImage := ebiten.NewImage(oldX, oldY)
	op := &ebiten.DrawImageOptions{}
	outputImage.DrawImage(orig, op)
	outputImage.DrawImage(tintImage, op)
	//outputImage.DrawImage(orig, op)
	return outputImage

}

func (p *BrickGrid) Draw(screen *ebiten.Image) {
	for iter, val := range p.brickGrid {
		_ = iter
		if val != nil && val.alive {
			brickImg := p.coloredBrickImages[val.gridY]
			DrawImageAt(screen, brickImg, val.screenX, val.screenY)
		}

	}

}

func (p *BrickGrid) bounceBallOffBrick(brick *Brick) {
	midPointBallX, midPointBallY := p.game.ball.midpoint()
	midPointBrickX, midPointBrickY := brick.midpoint()
	if midPointBallX > midPointBrickX {
		p.game.ball.velX = abs(p.game.ball.velX)
	} else if midPointBallX < midPointBrickX {
		p.game.ball.velX = -abs(p.game.ball.velX)
	}

	if midPointBallY > midPointBrickY {
		p.game.ball.velY = abs(p.game.ball.velY)
	} else if midPointBallY < midPointBrickY {
		p.game.ball.velY = -abs(p.game.ball.velY)
	}

}

func (p *BrickGrid) updateScore(val *Brick) {
	points := 1 + bricksY - val.gridY
	p.game.score += points
	scoreText := fmt.Sprintf("Score: %d", p.game.score)
	p.game.rasterStrings.scoreString.stringContent = scoreText

}

func (p *BrickGrid) Update() {
	//fmt.Println(p.game.levelCompleteScreenActive)
	if p.checkPlayerWon() && !p.game.levelCompleteScreenActive {
		p.game.ball.stop()
		p.game.levelCompleteScreenActive = true
		p.game.levelComplete()
	}

	ballRect := p.game.ball.getColliderRect()
	for iter, val := range p.brickGrid {
		_ = iter
		if val != nil && val.alive {

			brickRect := rect{val.screenX, val.screenY, brickWidth, brickHeight}

			if collideRect(brickRect, ballRect) {
				p.brickGrid[iter].alive = false
				//fmt.Println("Collide ", iter)
				p.bounceBallOffBrick(val)
				p.updateScore(val)
			}
		}

	}

}
