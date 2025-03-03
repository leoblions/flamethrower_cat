package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BOSS_ENT_KIND = 10
	BOSS_PART_W   = 100
	BOSS_PART_H   = 100
)

type BossManager struct {
	game *Game
	kind int
	MobileObject
	bodyPartEntities []*Entity
}

type BossPart struct {
	MobileObject
	currentImage *ebiten.Image
	imageID      int
	health       int

	frame int
}

func (b *BossManager) initBodyPartEntities() []*Entity {
	parts := []*Entity{}
	return parts
}

func (b *BossManager) initImages() []*ebiten.Image {
	images := []*ebiten.Image{}
	return images
}

func (b *BossManager) removeAllBossParts() {
	for i, v := range b.game.entityManager.entityList {
		if v.kind == BOSS_ENT_KIND {
			b.game.entityManager.entityList[i] = nil
		}
	}
}

func (em *EntityManager) addLevelBoss() {
	currentLevel := em.game.level
	if currentLevel == 20 {
		fmt.Println("ADD BARNACLE BOSS")
	}

}
