package main

import (
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	subdir = "images"
	//paddleImage  = "longPaddle.png"
	playerImageL      = "ffcatL.png"
	playerImageR      = "ffcatR.png"
	defaultSpeed      = 5
	playerHeight      = 100
	playerWidth       = 100
	xMax              = screenWidth - playerWidth
	PL_JUMP_HEIGHT    = 20
	PL_GRAVITY_AMOUNT = 1
	PL_RUN_BOOST      = 5
)

type Player struct {
	game    *Game
	worldX  int
	worldY  int
	screenX int
	screenY int

	velX      float32
	velY      float32
	currImage *ebiten.Image
	imageL    *ebiten.Image
	imageR    *ebiten.Image
	frozen    bool
	hoverMode bool
	faceLeft  bool
	run       bool
}

func (p *Player) getWorldColliderRect() rect {
	return rect{p.worldX, p.worldY, playerWidth, playerHeight}
}
func (p *Player) getColliderRect() rect {
	return rect{p.screenX, p.screenY, playerWidth, playerHeight}
}

func (p *Player) midpointX() int {
	return p.screenX + (playerWidth / 2)

}

func NewPlayer(g *Game, startX, startY int) *Player {
	p := &Player{}
	p.game = g
	var err error
	//cwd, _ := os.Getwd()
	var rawImage, stretchedImage *ebiten.Image

	imageDir := path.Join(subdir, playerImageL)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	p.imageL = copyAndStretchImage(rawImage, playerWidth, playerHeight)

	imageDir = path.Join(subdir, playerImageR)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	p.imageR = copyAndStretchImage(rawImage, playerWidth, playerHeight)

	if err != nil {
		log.Fatal(err)
	}
	p.currImage = p.imageR
	_ = stretchedImage
	p.worldX = startX
	p.worldY = startY
	p.hoverMode = false
	p.run = false

	//p.velX = 2
	return p
}

func (p *Player) Draw(screen *ebiten.Image) {
	DrawImageAt(screen, p.currImage, p.screenX, p.screenY)

}

func (p *Player) Update() {

	dflags := p.game.input.dflags

	solidBelowPlayer := p.game.tileMap.solidUnderPlayer(0)
	// if !solidBelowPlayer {
	// 	p.velY += PL_GRAVITY_AMOUNT
	// } else {
	// 	if p.velY > 0 {
	// 		p.velY = 0
	// 	}
	// }
	var jump bool = false
	//var fall bool = false

	if dflags.up && !dflags.down && (solidBelowPlayer || p.hoverMode) {
		jump = true
		if p.hoverMode {
			p.velY -= defaultSpeed
		}
	} else if (!dflags.up && dflags.down) && p.hoverMode {
		p.velY += defaultSpeed
	} else if p.hoverMode {
		p.velY = attenuate(p.velY, 1.0)
	}

	if dflags.left && !dflags.right {
		p.velX = -defaultSpeed
		if p.run {
			p.velX -= PL_RUN_BOOST
		}
	} else if !dflags.left && dflags.right {
		p.velX = defaultSpeed
		if p.run {
			p.velX += PL_RUN_BOOST
		}
	} else {
		p.velX = attenuate(p.velX, 1.0)
	}
	p.run = false
	testRect := rect{p.worldX + int(p.velX), p.worldY + int(p.velY), playerWidth, playerHeight}
	sideCollisions := p.game.tileMap.getSideCollisionData(testRect)

	dflags.reset()

	if !p.hoverMode {
		p.velY += PL_GRAVITY_AMOUNT
	}

	if jump && solidBelowPlayer && !p.hoverMode {
		p.velY -= PL_JUMP_HEIGHT
	}

	p.velY = clamp(-20.0, 20.0, p.velY)

	if !p.frozen {
		//fmt.Println("move player")
		if !sideCollisions.up && p.velY < 0 {
			p.worldY += int(p.velY)

		}
		if !sideCollisions.down && p.velY > 0 {

			p.worldY += int(p.velY)

		} else if sideCollisions.down && p.velY > 0 {
			p.velY = 0
		}

		if !sideCollisions.left && p.velX < 0 {
			p.worldX += int(p.velX)
			p.currImage = p.imageL
			p.faceLeft = true

		}
		if !sideCollisions.right && p.velX > 0 {
			p.worldX += int(p.velX)
			p.currImage = p.imageR
			p.faceLeft = false

		}
		p.screenX = p.worldX - worldOffsetX
		p.screenY = p.worldY - worldOffsetY
	}

}

func (p *Player) warpPlayerToGridLocation(gridX, gridY int) {

	p.worldX = gridX * p.game.tileMap.tileSize
	p.worldY = gridY * p.game.tileMap.tileSize

	p.screenX = p.worldX - worldOffsetX
	p.screenY = p.worldY - worldOffsetY

}
