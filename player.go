package main

import (
	"fmt"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	subdir = "images"
	//paddleImage  = "longPaddle.png"
	playerImageL          = "ffcatL.png"
	playerImageR          = "ffcatR.png"
	PL_SPRITE_SHEET       = "cattoSS.png"
	defaultSpeed          = 5
	playerHeight          = 100
	playerWidth           = 100
	PL_SPRITE_W           = 100
	PL_SPRITE_H           = 100
	xMax                  = screenWidth - playerWidth
	PL_JUMP_HEIGHT        = 20
	PL_GRAVITY_AMOUNT     = 1
	PL_RUN_BOOST          = 5
	PL_FRAME_CHANGE_TICKS = 20
)

type Player struct {
	game             *Game
	worldX           int
	worldY           int
	screenX          int
	screenY          int
	state            rune // s=stand w=walk f=fall/jump
	currentFrame     int
	currentTickCount int
	velX             float32
	velY             float32
	currImage        *ebiten.Image
	imageL           *ebiten.Image
	imageR           *ebiten.Image
	imageWalkL       []*ebiten.Image
	imageWalkR       []*ebiten.Image
	imageFallL       *ebiten.Image
	imageFallR       *ebiten.Image
	frozen           bool
	hoverMode        bool
	faceLeft         bool
	run              bool
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
	p.initImages()
	p.worldX = startX
	p.worldY = startY
	p.hoverMode = false
	p.run = false
	p.state = 's'
	p.currentFrame = 0
	//p.velX = 2
	return p
}

func (pl *Player) initImages() error {

	var err error
	//cwd, _ := os.Getwd()
	var rawImage *ebiten.Image

	imageDir := path.Join(subdir, playerImageL)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	pl.imageL = copyAndStretchImage(rawImage, playerWidth, playerHeight)

	imageDir = path.Join(subdir, playerImageR)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	pl.imageR = copyAndStretchImage(rawImage, playerWidth, playerHeight)

	if err != nil {
		log.Fatal(err)
	}
	pl.currImage = pl.imageR
	pl.imageWalkL = []*ebiten.Image{}
	//walk images
	pl.imageWalkL = grabImagesRowToListFromFilename(PL_SPRITE_SHEET, 100, 1, 4)
	pl.imageWalkR = copyAndReverseListOfImages(pl.imageWalkL)
	fmt.Println("tempFallImages ", len(pl.imageWalkR))
	//fall, jump images
	tempFallImages := grabImagesRowToListFromFilename(PL_SPRITE_SHEET, 100, 2, 1)
	fmt.Println("tempFallImages ", len(tempFallImages))
	pl.imageFallL = tempFallImages[0]
	pl.imageFallR = copyAndReverseListOfImages(tempFallImages)[0]

	return err

}

func (p *Player) Draw(screen *ebiten.Image) {
	DrawImageAt(screen, p.currImage, p.screenX, p.screenY)

}

func (p *Player) walkCycleNumber() int {
	if p.currentTickCount < PL_FRAME_CHANGE_TICKS {
		p.currentTickCount += 1
	} else {
		if p.currentFrame < 3 {
			p.currentFrame += 1
		} else {
			p.currentFrame = 0
		}
		p.currentTickCount = 0
	}
	return p.currentFrame
}

func (p *Player) selectImage() {
	switch p.state {
	case 'w':
		wsIndex := p.walkCycleNumber()
		if p.faceLeft {
			p.currImage = p.imageWalkL[wsIndex]
		} else {
			p.currImage = p.imageWalkR[wsIndex]
		}

	case 'f':
		if p.faceLeft {
			p.currImage = p.imageFallL
		} else {
			p.currImage = p.imageFallR
		}
	case 's':
		if p.faceLeft {
			p.currImage = p.imageL
		} else {
			p.currImage = p.imageR
		}
	default:
		if p.faceLeft {
			p.currImage = p.imageL
		} else {
			p.currImage = p.imageR
		}
	}

}

func (p *Player) playerMotion() {
	dflags := p.game.input.dflags

	solidBelowPlayer := p.game.tileMap.solidUnderPlayer(0)
	var jump bool = false
	plat := p.game.platformManager.playerStandingOnPlatform

	if plat {
		//fmt.Println("platform force up")
	}
	//playerFootHeight := p.worldX + playerHeight

	if dflags.up && !dflags.down && (solidBelowPlayer || plat || p.hoverMode) {
		jump = true
		// jumps can work
		if p.hoverMode {
			p.velY -= defaultSpeed
		}
	} else if !dflags.up && dflags.down {
		if p.hoverMode || plat {
			p.velY += defaultSpeed
			//plat = false //press down to fall thru platform
		}

	} else if p.hoverMode {
		p.velY = attenuate(p.velY, 1.0)
		//plat = false
	} else if plat && !dflags.up && !dflags.down {
		p.velY = float32(p.game.platformManager.touchingPlatformVelY)
		p.velX = float32(p.game.platformManager.touchingPlatformVelX)
		if p.velY > 0 {

			p.velY = -1

		}
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
	} else if !plat {
		p.velX = attenuate(p.velX, 1.0)
	}
	p.run = false
	testRect := rect{p.worldX + int(p.velX), p.worldY + int(p.velY), playerWidth, playerHeight}
	sideCollisions := p.game.tileMap.getSideCollisionData(testRect)

	dflags.reset()

	if !p.hoverMode && !plat {
		//fall
		p.velY += PL_GRAVITY_AMOUNT
	}

	if jump && (solidBelowPlayer || plat) && !p.hoverMode {
		//jump
		p.velY -= PL_JUMP_HEIGHT
		_ = playSound(p.game.soundEffectPlayers["jump"])
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
			//p.currImage = p.imageL
			p.faceLeft = true

		}
		if !sideCollisions.right && p.velX > 0 {
			p.worldX += int(p.velX)
			//p.currImage = p.imageR
			p.faceLeft = false

		}
		p.screenX = p.worldX - worldOffsetX
		p.screenY = p.worldY - worldOffsetY
	}
}

func (p *Player) updateState() {

	if p.velX != 0 {
		p.state = 'w'
	} else {
		p.state = 's'
	}
	if p.velY != 0 {
		p.state = 'f'
	}
}

func (p *Player) Update() {

	p.playerMotion()
	p.updateState()
	p.selectImage()
}

func (p *Player) warpPlayerToGridLocation(gridX, gridY int) {

	p.worldX = gridX * p.game.tileMap.tileSize
	p.worldY = gridY * p.game.tileMap.tileSize

	p.screenX = p.worldX - worldOffsetX
	p.screenY = p.worldY - worldOffsetY

}
