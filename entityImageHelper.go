package main

import (
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EntityManagerImageCollections_0 struct {
	jackieImages  []*ebiten.Image
	jackieImagesA []*ebiten.Image
	jackieImagesS []*ebiten.Image
	dogImages     []*ebiten.Image
	dogImagesA    []*ebiten.Image
	blobImages    []*ebiten.Image
	blobImagesA   []*ebiten.Image
	golemImages   []*ebiten.Image
	golemImagesA  []*ebiten.Image
}

type AllEntitySpriteCollections struct {
	jackieSC *EntitySpriteCollection
	dogSC    *EntitySpriteCollection
	blobSC   *EntitySpriteCollection
	golemSC  *EntitySpriteCollection
	antSC    *EntitySpriteCollection
}

type EntitySpriteCollection struct {
	walk   []*ebiten.Image
	attack []*ebiten.Image
	stand  []*ebiten.Image
}

func (em *EntityManager) initEntityImages() error {
	var err error
	// jackie
	em.jackieSC = &EntitySpriteCollection{}
	em.jackieSC.walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_JACKIE, 1)
	em.jackieSC.attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_JACKIE, 2)
	em.jackieSC.stand = grabImagesRowToListFromFilename(IMAGES_JACKIE, 100, 3, 2)
	//fmt.Println("Jackie images", len(em.jackieImages))
	// robo dog
	em.dogSC = &EntitySpriteCollection{}
	em.dogSC.walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)
	em.dogSC.attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)
	em.dogSC.stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)

	// worm blob
	em.blobSC = &EntitySpriteCollection{}
	em.blobSC.walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 1)
	em.blobSC.attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 2)
	em.blobSC.attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 2)

	// golem
	em.golemSC = &EntitySpriteCollection{}
	em.golemSC.walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 3)
	em.golemSC.attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 4)
	em.golemSC.stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 4)

	// ant
	em.antSC = &EntitySpriteCollection{}
	em.antSC.walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_ANT, 1)
	em.antSC.attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_ANT, 2)
	em.antSC.stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_ANT, 2)
	return err

}

func (em *EntityManager) updateFrame_0(ent *Entity) {
	rightIndexStart := 1 + EN_FRAME_MAX_VAL
	rightIndexEnd := 1 + EN_FRAME_MAX_VAL*2

	// low values left
	if ent.direction == 'l' && ent.frame >= EN_FRAME_MAX_VAL {
		ent.frame = 0
	} else if ent.direction == 'l' && ent.frame < EN_FRAME_MAX_VAL {
		ent.frame++
	} else if ent.direction == 'r' && ent.frame >= rightIndexEnd {
		ent.frame = rightIndexStart
	} else if ent.direction == 'r' && ent.frame < rightIndexEnd {
		ent.frame++
	}

}

func (em *EntityManager) updateFrame(ent *Entity) {

	// left: 0 1 2 3
	// right: 4 5 6 7

	var leftRune, rightRune rune
	// swap image directions if not robo dog
	// if ent.kind != 0 {
	// 	leftRune = 'l'
	// 	rightRune = 'r'
	// } else {
	// 	leftRune = 'r'
	// 	rightRune = 'l'
	// }
	leftRune = 'r'
	rightRune = 'l'
	changeDirection := false
	if ent.direction != ent.prevDirection {
		changeDirection = true
	}
	if ent.direction == leftRune {
		if ent.frame > 2 || changeDirection || ent.velX == 0 {
			ent.frame = 0
		} else {
			ent.frame++
		}
	} else if ent.direction == rightRune {
		if ent.frame > 6 || changeDirection || ent.velX == 0 {
			ent.frame = 4
		} else {
			ent.frame++
		}
	} else {
	}

}

func getForwardAndReverseSpriteRowFromFile(filename string, gridStartRow int) ([]*ebiten.Image, error) {
	tempImagesList := []*ebiten.Image{}
	imageDir := path.Join(subdir, filename)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	startYInPixels := gridStartRow * EN_SPRITE_SIZE
	tempImagesList = append(tempImagesList, grabImagesRowToList(rawImage, EN_SPRITE_SIZE, startYInPixels, EN_SPRITES_PER_ROW)...)
	//reverse
	tempImagesList = append(tempImagesList, copyAndReverseListOfImages(tempImagesList)...)

	return tempImagesList, err
}

func copyAndReverseListOfImages(inputList []*ebiten.Image) []*ebiten.Image {
	outputList := []*ebiten.Image{}
	for _, v := range inputList {
		tempImage := FlipHorizontal(v)
		outputList = append(outputList, tempImage)
	}
	return outputList
}

func grabImagesRowToList(inputImage *ebiten.Image, imageSize, startYInPixels, amountX int) []*ebiten.Image {

	outputList := []*ebiten.Image{}
	for i := range amountX {
		x := i * imageSize
		y := startYInPixels
		tempImage := getSubImage(inputImage, x, y, imageSize, imageSize)
		outputList = append(outputList, tempImage)
	}
	return outputList
}

func grabImagesRowToListFromFilename(filename string, imageSize, gridY, amountX int) []*ebiten.Image {
	imageDir := path.Join(subdir, filename)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	outputList := []*ebiten.Image{}
	for i := range amountX {
		x := i * imageSize
		y := gridY * imageSize
		tempImage := getSubImage(rawImage, x, y, imageSize, imageSize)
		outputList = append(outputList, tempImage)
	}
	if err != nil {
		log.Fatal(err)
	}
	return outputList
}

func (em *EntityManager) selectImage(entityKind, state, frame int) *ebiten.Image {
	var imageList []*ebiten.Image
	switch entityKind {

	case 0, 4: // jackie
		switch state {
		case 0:
			imageList = em.jackieSC.walk
		case 1:
			imageList = em.jackieSC.attack
		case 2:
			imageList = em.jackieSC.stand

		}
	case 1: // dog
		switch state {
		case 0:
			imageList = em.dogSC.walk
		case 1:
			imageList = em.dogSC.attack
		case 2:
			imageList = em.dogSC.stand

		}

	case 2: //blob
		switch state {
		case 0:
			imageList = em.blobSC.walk
		case 1:
			imageList = em.blobSC.attack
		case 2:
			imageList = em.blobSC.stand

		}
	case 3: // golem
		switch state {
		case 0:
			imageList = em.golemSC.walk
		case 1:
			imageList = em.golemSC.attack
		case 2:
			imageList = em.golemSC.stand

		}
	case 5: // ant
		switch state {
		case 0:
			imageList = em.antSC.walk
		case 1:
			imageList = em.antSC.attack
		case 2:
			imageList = em.antSC.stand

		}
	default:
		switch state {
		case 0:
			imageList = em.antSC.walk
		case 1:
			imageList = em.antSC.attack
		case 2:
			imageList = em.antSC.stand

		}

	}

	//fmt.Println(frame)
	return imageList[frame]

}
