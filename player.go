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
	playerImageL         = "ffcatL.png"
	playerImageR         = "ffcatR.png"
	PL_SPRITE_SHEET      = "cattoSS.png"
	defaultSpeed         = 5
	PL_SPRITE_W          = 100
	PL_SPRITE_H          = 100
	PL_COLLRECT_W        = 50
	PL_COLLRECT_H        = 100
	PL_COLLRECT_OFFSET_X = 25

	PL_JUMP_HEIGHT                = 20
	PL_JUMP_HEIGHT_W              = 1
	PL_GRAVITY_AMOUNT             = 1.0
	PL_GRAVITY_AMOUNT_W           = 0.2
	PL_RUN_BOOST                  = 5
	PL_FRAME_CHANGE_TICKS         = 20
	PL_HEALTH_MAX                 = 100
	PL_HALF_W                     = 50
	PL_LAVA_DAMAGE_AMOUNT         = 1
	PL_FRICTION_X                 = 1.0
	PL_FRICTION_X_W               = 1.0
	PL_MAX_VEL            float32 = 20.0
	PL_MAX_VEL_W          float32 = 2.0
	PL_MAX_VEL_LADDER     float32 = 5.0
	PL_DAMAGE_DEBOUNCE            = 100

	PLAYER_PLATFORM_FOOT_POS_X = 50
	PLAYER_PLATFORM_FOOT_POS_Y = 120

	PL_BUBBLE_PERIOD = 100
)

type Player struct {
	game                                         *Game
	xMax                                         int
	screenX                                      int
	screenY                                      int
	state                                        rune // s=stand w=walk f=fall/jump
	currentFrame                                 int
	currentTickCount                             int
	health                                       int
	bubbleTicks                                  int8
	velX                                         float32
	velY                                         float32
	currImage                                    *ebiten.Image
	imageL                                       *ebiten.Image
	imageR                                       *ebiten.Image
	imageWalkL, imageWalkR, imageDieR, imageDieL []*ebiten.Image
	imageFallL                                   *ebiten.Image
	imageFallR                                   *ebiten.Image
	imageIndex                                   int
	frozen                                       bool
	hoverMode                                    bool
	faceLeft                                     bool
	run                                          bool
	footUnderwater                               bool
	headUnderwater                               bool
	treadingWater                                bool
	touchingLadder                               bool
	standingOnPlatform                           bool
	damageDebounce                               int

	collRect     *rect
	collRectTile *rect
	testRect     *rect
	MobileObject
}

func (p *Player) getWorldColliderRect() rect {

	return *p.collRect
}
func (p *Player) getScreenCollrect() rect {
	return rect{p.screenX, p.screenY, PL_COLLRECT_W, PL_COLLRECT_H}
}

func (p *Player) midpointX() int {
	return p.screenX + (PL_COLLRECT_W / 2)

}

func NewPlayer(g *Game, startX, startY int) *Player {
	p := &Player{}
	p.game = g
	p.initImages()
	p.worldX = startX
	p.worldY = startY
	p.hoverMode = false
	p.run = false
	p.xMax = panelWidth - PL_COLLRECT_W
	p.health = 100
	p.state = 's'
	p.currentFrame = 0
	p.width = PL_COLLRECT_W
	p.height = PL_COLLRECT_H
	p.testRect = &rect{p.worldX, p.worldY, PL_COLLRECT_W, PL_COLLRECT_H}
	p.collRect = &rect{p.worldX, p.worldY, p.width, p.height}
	p.collRectTile = &rect{p.worldX + PL_COLLRECT_OFFSET_X,
		p.worldY, PL_COLLRECT_W, PL_COLLRECT_H}
	//p.velX = 2
	return p
}

func (pl *Player) initImages() error {

	var err error
	//cwd, _ := os.Getwd()
	var rawImage *ebiten.Image

	imageDir := path.Join(subdir, playerImageL)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	pl.imageL = copyAndStretchImage(rawImage, PL_SPRITE_W, PL_SPRITE_H)

	imageDir = path.Join(subdir, playerImageR)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	pl.imageR = copyAndStretchImage(rawImage, PL_SPRITE_W, PL_SPRITE_H)

	if err != nil {
		log.Fatal(err)
	}
	pl.currImage = pl.imageR
	pl.imageWalkL = []*ebiten.Image{}
	//walk images
	pl.imageWalkL = grabImagesRowToListFromFilename(PL_SPRITE_SHEET, 100, 1, 4)
	pl.imageWalkR = copyAndReverseListOfImages(pl.imageWalkL)
	//fmt.Println("tempFallImages ", len(pl.imageWalkR))
	//fall, jump images
	tempFallImages := grabImagesRowToListFromFilename(PL_SPRITE_SHEET, 100, 2, 1)
	//fmt.Println("tempFallImages ", len(tempFallImages))
	pl.imageFallL = tempFallImages[0]
	pl.imageFallR = copyAndReverseListOfImages(tempFallImages)[0]

	tempImages := grabImagesRowToListFromFilename(PL_SPRITE_SHEET, 100, 3, 4)
	//fmt.Println("tempFallImages ", len(tempFallImages))
	pl.imageDieL = tempImages
	pl.imageDieR = copyAndReverseListOfImages(tempImages)

	return err

}

func (p *Player) Draw(screen *ebiten.Image) {
	DrawImageAt(screen, p.currImage, p.screenX, p.screenY)

}

func (p *Player) motionWaterPhysics() (float32, float32, float32) {
	var gravityAmount, friction, jumpHeight float32
	if p.footUnderwater && p.headUnderwater {
		gravityAmount = PL_GRAVITY_AMOUNT_W
		friction = PL_FRICTION_X_W
		jumpHeight = PL_JUMP_HEIGHT_W

	} else if p.footUnderwater && !p.headUnderwater {
		gravityAmount = 0
		friction = PL_FRICTION_X
		jumpHeight = PL_JUMP_HEIGHT
	} else if p.hoverMode {
		gravityAmount = 0
		friction = PL_FRICTION_X
		jumpHeight = PL_JUMP_HEIGHT
	} else {
		gravityAmount = PL_GRAVITY_AMOUNT
		friction = PL_FRICTION_X
		jumpHeight = PL_JUMP_HEIGHT
	}
	return gravityAmount, friction, jumpHeight
}

func (p *Player) motionTileCollision() {

}

func (p *Player) motionSpeedLimit() {
	if p.footUnderwater && !p.treadingWater {
		p.velY = clamp(-PL_MAX_VEL_W, PL_MAX_VEL_W, p.velY)
		p.velX = clamp(-PL_MAX_VEL_W, PL_MAX_VEL_W, p.velX)
	} else if p.touchingLadder {
		p.velY = clamp(-PL_MAX_VEL_LADDER, PL_MAX_VEL_LADDER, p.velY)
		p.velX = clamp(-PL_MAX_VEL, PL_MAX_VEL, p.velX)
	} else {
		p.velY = clamp(-PL_MAX_VEL, PL_MAX_VEL, p.velY)
		p.velX = clamp(-PL_MAX_VEL, PL_MAX_VEL, p.velX)
	}

}

func (p *Player) playerMotion() {
	dflags := p.game.input.dflags
	// underwater physics
	p.headUnderwater = p.checkHeadUnderWater()

	gravityAmount, friction, jumpHeight := p.motionWaterPhysics()
	p.treadingWater = !p.headUnderwater && p.footUnderwater
	solidBelowPlayer := p.game.tileMap.solidUnderPlayer(0)
	p.standingOnPlatform = p.game.platformManager.playerStandingOnPlatform
	canJump := (solidBelowPlayer || p.standingOnPlatform || p.hoverMode || p.footUnderwater || p.treadingWater)

	var jump bool = false

	if p.standingOnPlatform {
	}

	if dflags.up && !dflags.down {
		// jump or float up
		if canJump {
			jump = true
		} else if p.touchingLadder || p.hoverMode {
			p.velY -= defaultSpeed

		}

	} else if !dflags.up && dflags.down {
		if p.hoverMode || p.standingOnPlatform || p.footUnderwater || p.touchingLadder {
			p.velY += defaultSpeed

		}

	} else if !dflags.down && !dflags.up && (p.treadingWater) && !p.touchingLadder {
		// don't sink if treading water
		p.velY = 0
	} else if p.hoverMode || p.touchingLadder {
		p.velY = attenuate(p.velY, friction)
	} else if p.standingOnPlatform && !dflags.up && !dflags.down {
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
	} else if !p.standingOnPlatform {
		p.velX = attenuate(p.velX, friction)
	}
	p.run = false

	p.updateTestRect(p.velX, p.velY)

	sideCollisions := p.game.tileMap.getSideCollisionData(*p.testRect)

	dflags.reset()

	if !p.hoverMode && !p.standingOnPlatform && !p.treadingWater && !p.touchingLadder {
		// fall gravity

		p.velY += gravityAmount

	} else {
	}

	if jump && !p.hoverMode {
		//jump
		p.velY -= jumpHeight
		if !p.footUnderwater {
			p.game.audioPlayer.playSoundByID("jump")
		}

		//_ = playSound(p.game.audioPlayer.soundEffectPlayers["jump"])
	} else if jump && p.footUnderwater {
		p.velY -= jumpHeight
		//p.game.audioPlayer.playSoundByID("jump")
	} else if jump && p.hoverMode {
		p.velY -= jumpHeight
	}

	p.motionSpeedLimit()

	//don't move if dead
	if p.state == 'd' {
		p.velX = 0
		p.velY = clamp(-10, 0, p.velY)
	}
	//apply velocity to position
	if !p.frozen {

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

			p.faceLeft = true

		}
		if !sideCollisions.right && p.velX > 0 {
			p.worldX += int(p.velX)

			p.faceLeft = false

		}

	}

	// cleanup
	p.screenX = p.worldX - worldOffsetX
	p.screenY = p.worldY - worldOffsetY
	p.updateCollrect()
	p.touchingLadder = false

}

func (p *Player) updateTestRect(velX, velY float32) {
	p.testRect.x = p.worldX + int(p.velX) + PL_COLLRECT_OFFSET_X
	p.testRect.y = p.worldY + int(p.velY)
	p.testRect.width = PL_COLLRECT_W
	p.testRect.height = PL_COLLRECT_H

}

func (p *Player) updateCollrect() {
	p.collRect.x = p.worldX + 30
	p.collRect.width = PL_COLLRECT_W
	p.collRect.height = PL_COLLRECT_H
	p.collRect.y = p.worldY
	p.collRectTile.x = p.worldX + PL_COLLRECT_OFFSET_X
	p.collRectTile.y = p.worldY

}

func (p *Player) checkPlayerUnderwater() bool {
	pX := p.worldX + PL_HALF_W
	pY := p.worldY + 50
	return p.game.tileMap.pointCollidedWithGivenTileKind(pX, pY, TM_WATER_TILE_ID)

}

func (p *Player) checkFootUnderwater() bool {
	pX := p.worldX + PL_HALF_W
	pY := p.worldY + PL_SPRITE_H
	return p.game.tileMap.pointCollidedWithGivenTileKind(pX, pY, TM_WATER_TILE_ID)

}

func (p *Player) checkPlayerStandingOnLava() bool {
	pX := p.worldX + PL_HALF_W
	pY := p.worldY + PL_SPRITE_H
	return p.game.tileMap.pointCollidedWithGivenTileKind(pX, pY, TM_LAVA_TILE_ID)

}

func (p *Player) checkHeadUnderWater() bool {
	pX := p.worldX + PL_HALF_W
	pY := p.worldY
	return p.game.tileMap.pointCollidedWithGivenTileKind(pX, pY, TM_WATER_TILE_ID)

}

func (p *Player) Update() {
	if p.damageDebounce > 0 {
		p.damageDebounce -= 1
	}
	p.playerMotion()
	p.updateState()
	p.selectImage()

	if p.checkPlayerStandingOnLava() {
		p.changeHealth(-PL_LAVA_DAMAGE_AMOUNT)
		p.game.audioPlayer.playSoundByID("lavahiss")
	}

	p.footUnderwater = p.checkPlayerUnderwater()
	if p.headUnderwater {
		p.bubbleEmitter()
	}
}

func (p *Player) warpPlayerToGridLocation(gridX, gridY int) {

	p.worldX = gridX * p.game.tileMap.tileSize
	p.worldY = gridY * p.game.tileMap.tileSize

	p.screenX = p.worldX - worldOffsetX
	p.screenY = p.worldY - worldOffsetY

}

func (p *Player) changeHealth(hitpoints int) {

	healthTemp := p.health
	healthTemp += hitpoints
	healthTemp = clamp(0, PL_HEALTH_MAX, healthTemp)
	p.health = healthTemp
	p.game.healthBar.FillPercent(p.health)

}

func (p *Player) takeDamageWithDebounce(damage int) {
	if p.damageDebounce <= 0 {
		p.damageDebounce = PL_DAMAGE_DEBOUNCE
		healthTemp := p.health
		healthTemp -= damage
		healthTemp = clamp(0, PL_HEALTH_MAX, healthTemp)
		p.health = healthTemp
		p.game.healthBar.FillPercent(p.health)
	}

}
