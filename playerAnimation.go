package main

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

func (p *Player) fallCycleNumber() int {
	if p.currentTickCount < PL_FRAME_CHANGE_TICKS*2 {
		p.currentTickCount += 1
	} else {
		if p.currentFrame < 3 {
			p.currentFrame += 1
		} else {
			//p.currentFrame = 0
		}
		//p.currentTickCount = 0
	}
	return p.currentFrame
}

func (p *Player) selectImage() {
	switch p.state {
	case 'w':
		p.imageIndex = p.walkCycleNumber()
		if p.faceLeft {
			p.currImage = p.imageWalkL[p.imageIndex]
		} else {
			p.currImage = p.imageWalkR[p.imageIndex]
		}
	case 'd':
		p.imageIndex = p.fallCycleNumber()
		if p.faceLeft {
			p.currImage = p.imageDieL[p.imageIndex]
		} else {
			p.currImage = p.imageDieR[p.imageIndex]
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

func (p *Player) updateState() {
	footX := p.collRect.x + PLAYER_PLATFORM_FOOT_POS_X
	footY := p.collRect.y + PLAYER_PLATFORM_FOOT_POS_Y
	if p.game.platformManager.pointCollidesWithPlatform(footX, footY) {
		p.state = 's'
	} else {
		if p.velX != 0 {
			p.state = 'w'
		} else {
			p.state = 's'
		}
		if p.velY != 0 {
			p.state = 'f'
		}
		if p.health <= 0 && !p.game.godMode && p.state != 'd' {
			p.state = 'd'
			p.imageIndex = 0
		}
	}

}

func (p *Player) bubbleEmitter() {

	if p.bubbleTicks < PL_BUBBLE_PERIOD {
		p.bubbleTicks++
	} else {
		p.bubbleTicks = 0
		p.game.particleManager.AddParticle(p.worldX+50, p.worldY, 1)
	}
}
