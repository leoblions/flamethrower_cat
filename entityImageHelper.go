package main

import (
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// type AllEntitySpriteCollections struct {
// 	jackieSC *EntitySpriteCollection
// 	dogSC    *EntitySpriteCollection
// 	blobSC   *EntitySpriteCollection
// 	golemSC  *EntitySpriteCollection
// 	antSC    *EntitySpriteCollection
// 	flySC    *EntitySpriteCollection
// }

type EntitySpriteCollection struct {
	walk   []*ebiten.Image
	attack []*ebiten.Image
	stand  []*ebiten.Image
}

func (em *EntityManager) initEntityImages() error {
	var err error
	// jackie 0
	ce := 0
	em.esCollections[ce] = &EntitySpriteCollection{}
	em.esCollections[0].walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_JACKIE, 1)
	em.esCollections[0].attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_JACKIE, 2)
	em.esCollections[0].stand = grabImagesRowToListFromFilename(IMAGES_JACKIE, 100, 3, 2)
	//fmt.Println("Jackie images", len(em.jackieImages))
	// robo dog 1
	ce = 1
	em.esCollections[ce] = &EntitySpriteCollection{}
	em.esCollections[1].walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)
	em.esCollections[1].attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)
	em.esCollections[1].stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)

	// worm blob 2
	ce = 2
	em.esCollections[ce] = &EntitySpriteCollection{}
	em.esCollections[2].walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 1)
	em.esCollections[2].attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 2)
	em.esCollections[2].stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 2)

	// golem 3
	ce = 3
	em.esCollections[3] = &EntitySpriteCollection{}
	em.esCollections[3].walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 3)
	em.esCollections[3].attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 4)
	em.esCollections[3].stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 4)

	// ant 5
	ce = 5
	em.esCollections[5] = &EntitySpriteCollection{}
	em.esCollections[5].walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_ANT, 1)
	em.esCollections[5].attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_ANT, 2)
	em.esCollections[5].stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_ANT, 2)

	// fly 6
	ce = 6
	em.esCollections[ce] = &EntitySpriteCollection{}
	em.esCollections[ce].walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_FLY, 1)
	em.esCollections[ce].attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_FLY, 2)
	em.esCollections[ce].stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_FLY, 2)
	// em.esCollections[6].walk = grabImagesRowToListFromFilenameWH(IMAGES_FLY, 200, 100, 0, 2)
	// em.esCollections[6].attack = grabImagesRowToListFromFilenameWH(IMAGES_FLY, 200, 100, 1, 2)
	// em.esCollections[6].stand = grabImagesRowToListFromFilenameWH(IMAGES_FLY, 200, 100, 1, 2)

	// shark 7
	ce = 7
	em.esCollections[ce] = &EntitySpriteCollection{}
	em.esCollections[ce].walk, err = getForwardAndReverseSpriteRowFromFile(IMAGES_SHARK, 1)
	em.esCollections[ce].attack, err = getForwardAndReverseSpriteRowFromFile(IMAGES_SHARK, 2)
	em.esCollections[ce].stand, err = getForwardAndReverseSpriteRowFromFile(IMAGES_SHARK, 2)
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

func grabImagesRowToListFromFilenameWH(filename string, imageWidth, imageHeight, gridY, amountX int) []*ebiten.Image {
	imageDir := path.Join(subdir, filename)
	rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	outputList := []*ebiten.Image{}
	for i := range amountX {
		x := i * imageWidth
		y := gridY * imageHeight
		tempImage := getSubImage(rawImage, x, y, imageWidth, imageHeight)
		outputList = append(outputList, tempImage)
	}
	if err != nil {
		log.Fatal(err)
	}
	return outputList
}

func (em *EntityManager) selectImage(entityKind, state, frame int) *ebiten.Image {
	var imageList []*ebiten.Image
	if entityKind == 4 {
		entityKind = 0
	}

	switch state {
	case 0:
		imageList = em.esCollections[entityKind].walk
	case 1:
		imageList = em.esCollections[entityKind].attack
	case 2:
		imageList = em.esCollections[entityKind].stand

	}

	//fmt.Println(frame)
	return imageList[frame]

}
