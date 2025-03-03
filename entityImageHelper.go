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

const (
	IMAGES_JACKIE   = "jackieRS.png"
	IMAGES_MONSTER  = "entitySheet1.png"
	IMAGES_ANT      = "entityAnt.png"
	IMAGES_FLY      = "entityFly.png"
	IMAGES_SHARK    = "entityShark.png"
	IMAGES_BIRD     = "entityBird.png"
	IMAGES_EARWIG   = "entityEarwig.png"
	IMAGES_B        = "B.png"
	IMAGES_BARNACLE = "barnacleFish.png"

	BOSS_SCALE_FACTOR_B = 3.0
)

type EntitySpriteCollection struct {
	walk   []*ebiten.Image
	attack []*ebiten.Image
	stand  []*ebiten.Image
	hit    []*ebiten.Image
}

func (em *EntityManager) initEntityImages() error {
	var err error
	// jackie 0
	ce := 0
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_JACKIE, 1, 2, 1, 1)
	em.esCollections[ce].stand = grabImagesRowToListFromFilename(IMAGES_JACKIE, 100, 3, 2)

	// robo dog 1
	ce = 1
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_MONSTER, 0, 0, 0, 0)

	// worm blob 2
	ce = 2
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_MONSTER, 1, 2, 2, 1)

	// golem 3
	ce = 3
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_MONSTER, 3, 4, 4, 3)

	// ant 5
	ce = 5
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_ANT, 1, 2, 2, 2)

	// fly 6
	ce = 6
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_FLY, 1, 2, 2, 2)

	// shark 7
	ce = 7
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_SHARK, 1, 2, 2, 2)

	// bird 8
	ce = 8
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_BIRD, 1, 2, 2, 2)

	// earwig 9
	ce = 9
	em.esCollections[ce], err = em.getEntityTypeSpriteCollectionX4(IMAGES_EARWIG, 1, 2, 2, 2)

	// moth boss 10
	ce = 10
	em.esCollections[ce] = &EntitySpriteCollection{}
	em.esCollections[ce].walk = grabImagesRowToListFromFilename(IMAGES_B, 100, 0, 1)
	em.esCollections[ce].attack = grabImagesRowToListFromFilename(IMAGES_B, 100, 0, 1)
	em.esCollections[ce].stand = grabImagesRowToListFromFilename(IMAGES_B, 100, 0, 1)
	em.esCollections[ce].hit = grabImagesRowToListFromFilename(IMAGES_B, 100, 0, 1)

	// barnacle boss 11
	ce = 11
	em.esCollections[ce] = &EntitySpriteCollection{}
	em.esCollections[ce].walk, err = getForwardAndReverseSpriteRowFromFileWHYA(IMAGES_BARNACLE, 100, 200, 0, 4)
	em.esCollections[ce].attack, err = getForwardAndReverseSpriteRowFromFileWHYA(IMAGES_BARNACLE, 100, 200, 0, 4)
	em.esCollections[ce].stand, err = getForwardAndReverseSpriteRowFromFileWHYA(IMAGES_BARNACLE, 100, 200, 1, 4)
	em.esCollections[ce].hit, err = getForwardAndReverseSpriteRowFromFileWHYA(IMAGES_BARNACLE, 100, 200, 1, 4)

	em.esCollections[ce].walk = rescaleSliceOfImages(em.esCollections[ce].walk, BOSS_SCALE_FACTOR_B)
	em.esCollections[ce].attack = rescaleSliceOfImages(em.esCollections[ce].attack, BOSS_SCALE_FACTOR_B)
	em.esCollections[ce].stand = rescaleSliceOfImages(em.esCollections[ce].stand, BOSS_SCALE_FACTOR_B)
	em.esCollections[ce].hit = rescaleSliceOfImages(em.esCollections[ce].hit, BOSS_SCALE_FACTOR_B)

	return err

}

func (em *EntityManager) getEntityTypeSpriteCollectionX4(imagePath string, rowWalk, rowAttack, rowStand, rowHit int) (*EntitySpriteCollection, error) {
	esCollection := &EntitySpriteCollection{}
	var err error
	esCollection.walk, err = getForwardAndReverseSpriteRowFromFile(imagePath, rowWalk)
	esCollection.attack, err = getForwardAndReverseSpriteRowFromFile(imagePath, rowAttack)
	esCollection.stand, err = getForwardAndReverseSpriteRowFromFile(imagePath, rowStand)
	esCollection.hit, err = getForwardAndReverseSpriteRowFromFile(imagePath, rowHit)
	return esCollection, err
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

	stopAnimating := (ent.velX == 0 && ent.kind != 11)

	leftRune = 'r'
	rightRune = 'l'
	changeDirection := false
	if ent.direction != ent.prevDirection {
		changeDirection = true
	}
	if ent.direction == leftRune {
		if ent.frame > 2 || changeDirection || stopAnimating {
			ent.frame = 0
		} else {
			ent.frame++
		}
	} else if ent.direction == rightRune {
		if ent.frame > 6 || changeDirection || stopAnimating {
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

func getForwardAndReverseSpriteRowFromFileWHYA(filename string, imageWidth, imageHeight, gridY, amountX int) ([]*ebiten.Image, error) {
	tempImagesList := []*ebiten.Image{}
	var err error
	path := path.Join(subdir, filename)
	//rawImage, _, err := ebitenutil.NewImageFromFile(imageDir)
	startYInPixels := gridY * imageHeight
	tempImagesList = append(tempImagesList, grabImagesRowToListFromPathWH(path, imageWidth, imageHeight, startYInPixels, amountX)...)
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
	path := path.Join(subdir, filename)

	return grabImagesRowToListFromPathWH(path, imageWidth, imageHeight, gridY, amountX)
}

func grabImagesRowToListFromPathWH(path string, imageWidth, imageHeight, gridY, amountX int) []*ebiten.Image {
	//imageDir := path.Join(subdir, filename)
	rawImage, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal("grabImagesRowToListFromFilenameWH\n", err)
	}
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
	case 3:
		imageList = em.esCollections[entityKind].hit

	}

	//fmt.Println(frame)
	return imageList[frame]

}
