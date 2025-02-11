package main

import (
	"fmt"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type EntityManagerImageCollections struct {
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

func (em *EntityManager) initEntityImages() error {
	var err error
	// jackie
	em.jackieImages, err = getForwardAndReverseSpriteRowFromFile(IMAGES_JACKIE, 1)
	em.jackieImagesA, err = getForwardAndReverseSpriteRowFromFile(IMAGES_JACKIE, 2)
	em.jackieImagesS = grabImagesRowToListFromFilename(IMAGES_JACKIE, 100, 3, 2)
	fmt.Println("Jackie images", len(em.jackieImages))
	// robo dog
	em.dogImages, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)
	em.dogImagesA, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 0)

	// worm blob

	em.blobImages, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 1)
	em.blobImagesA, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 2)

	// worm blob

	em.golemImages, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 3)
	em.golemImagesA, err = getForwardAndReverseSpriteRowFromFile(IMAGES_MONSTER, 3)

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

	if ent.kind != 0 {
		leftRune = 'l'
		rightRune = 'r'
	} else {
		leftRune = 'r'
		rightRune = 'l'
	}

	if ent.direction == leftRune {
		if ent.frame > 2 || ent.velX == 0 {
			ent.frame = 0
		} else {
			ent.frame++
		}
	} else if ent.direction == rightRune {
		if ent.frame > 6 || ent.velX == 0 {
			ent.frame = 4
		} else {
			ent.frame++
		}
	} else {
		//fmt.Println("Ent unknown condition")
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

	case 0, 4:
		switch state {
		case 0:
			imageList = em.jackieImages
		case 1:
			imageList = em.jackieImagesA
		case 2:
			imageList = em.jackieImagesS

		}
	case 1:
		switch state {
		case 0:
			imageList = em.dogImages
		case 1:
			imageList = em.dogImagesA

		}

	case 2:
		switch state {
		case 0:
			imageList = em.blobImages
		case 1:
			imageList = em.blobImagesA

		}
	case 3:
		switch state {
		case 0:
			imageList = em.golemImages
		case 1:
			imageList = em.golemImagesA

		}

	}
	//fmt.Println(frame)
	return imageList[frame]

}
