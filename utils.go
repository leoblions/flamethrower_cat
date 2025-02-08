package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	UTIL_LINE_SEPARATOR   = '\n'
	UTIL_NUMBER_SEPARATOR = ','
)

func loadMapFromFile(filePath string) ([mapRows][mapCols]int, error) {
	stringList := []*string{}
	intMatrix := [mapRows][mapCols]int{}

	file, err := os.Open(filePath)
	defer file.Close()
	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if line != nil && len(line) > 0 {
			//fmt.Printf("ReadLine: %q\n", line)
			lineAsString := string(line)
			stringList = append(stringList, &lineAsString)
		} else {
			break
		}
		if err != nil {
			log.Println(err)
			return intMatrix, err
		}
	}
	if err != nil {
		log.Println(err)
		return intMatrix, err
	}
	cols := len(strings.Split(*stringList[0], string(UTIL_NUMBER_SEPARATOR)))
	//intMatrix := [mapRows][mapCols]int{}
	for y := 0; y < cols; y++ {
		tempRow := [mapCols]int{}
		lineString := strings.Split(*stringList[y], string(UTIL_NUMBER_SEPARATOR))
		for x, value := range lineString {
			if value == "\x00" {
				continue
			}
			intValue, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}
			tempRow[x] = intValue

		}
		intMatrix[y] = tempRow

	}
	return intMatrix, err

}

func loadDataListFromFile(filePath string) ([][]int, error) {
	stringList := []*string{}
	type innerList []int
	intMatrix := [][]int{}

	file, err := os.Open(filePath)
	defer file.Close()
	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if line != nil && len(line) > 0 {
			//fmt.Printf("ReadLine: %q\n", line)
			lineAsString := string(line)
			stringList = append(stringList, &lineAsString)
		} else {
			break
		}
		if err != nil {
			log.Println(err)
			return intMatrix, err
		}
	}
	if err != nil {
		log.Println(err)
		return intMatrix, err
	}
	rows := len(stringList)
	if rows == 0 {
		return intMatrix, nil
	}
	cols := len(strings.Split(*stringList[0], string(UTIL_NUMBER_SEPARATOR)))
	//intMatrix := [mapRows][mapCols]int{}
	_ = cols
	for y := 0; y < rows; y++ {
		tempRow := []int{}
		lineString := strings.Split(*stringList[y], string(UTIL_NUMBER_SEPARATOR))
		for _, value := range lineString {
			if value == "\x00" {
				continue
			}
			strippedValue := strings.Trim(value, " ")
			intValue, err := strconv.Atoi(strippedValue)
			if err != nil {
				log.Fatal(err)
			}
			tempRow = append(tempRow, intValue)

		}
		intMatrix = append(intMatrix, tempRow)

	}
	return intMatrix, err

}

func writeMapToFile(mapData [mapRows][mapCols]int, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
		log.Fatal(fmt.Sprintf("writeMapToFile failed to write file: %s/n", filePath))
		return err
	}
	defer file.Close() // close file at end of this function
	rows := len(mapData)
	cols := len(mapData[0])
	var sb strings.Builder
	writer := bufio.NewWriter(file)
	for y := 0; y < rows; y++ {
		sb.Reset()
		for x := 0; x < cols; x++ {
			intVal := mapData[y][x]
			//sb.WriteRune(rune(intVal - 48))
			sb.WriteString(strconv.Itoa(intVal))

			if x != cols-1 {
				sb.WriteRune(UTIL_NUMBER_SEPARATOR)
			} else {
				sb.WriteRune(UTIL_LINE_SEPARATOR)
			}
		}
		_, err := writer.WriteString(sb.String())
		if err != nil {
			return err
		}

	}
	if err := writer.Flush(); err != nil {
		log.Fatal(err)
		return err
	} else {
		log.Println(fmt.Sprintf("Wrote file %s successfully.", filePath))
		return nil
	}

}

func write2DIntListToFile(list2D [][]int, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
		log.Fatal(fmt.Sprintf("writeMapToFile failed to write file: %s/n", filePath))
		return err
	}
	defer file.Close() // close file at end of this function
	rows := len(list2D)
	var cols int
	if rows == 0 {
		cols = 0
	} else {
		cols = len(list2D[0])
	}

	var sb strings.Builder
	writer := bufio.NewWriter(file)
	for y := 0; y < rows; y++ {
		sb.Reset()
		for x := 0; x < cols; x++ {
			intVal := list2D[y][x]
			//sb.WriteRune(rune(intVal - 48))
			sb.WriteString(strconv.Itoa(intVal))

			if x != cols-1 {
				sb.WriteRune(UTIL_NUMBER_SEPARATOR)
			} else {
				sb.WriteRune(UTIL_LINE_SEPARATOR)
			}
		}
		_, err := writer.WriteString(sb.String())
		if err != nil {
			return err
		}

	}
	if err := writer.Flush(); err != nil {
		log.Fatal(err)
		return err
	} else {
		log.Println(fmt.Sprintf("Wrote file %s successfully.", filePath))
		return nil
	}

}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func getListOfLinesFromFile(filename string) []*string {
	var linesList = []*string{}
	// open file handle
	file, err := os.Open(filename)
	defer file.Close()

	r := bufio.NewReader(file)

	// Section 2
	for {
		line, _, err := r.ReadLine()
		if line != nil && len(line) > 0 {
			//fmt.Printf("ReadLine: %q\n", line)
			lineAsString := string(line)
			linesList = append(linesList, &lineAsString)
		} else {
			break
		}
		if err != nil {
			break
		}
	}

	_ = err

	return linesList
}

func getSubImage(orig *ebiten.Image, x, y, width, height int) *ebiten.Image {
	outputImage := ebiten.NewImage(width, height)

	rectangle := image.Rect(x, y, x+width, y+height)
	cutImage := orig.SubImage(rectangle).(*ebiten.Image)
	op := &ebiten.DrawImageOptions{}
	outputImage.DrawImage(cutImage, op)
	return outputImage
}

func cutSpriteSheet(orig *ebiten.Image, spriteWidth, spriteHeight, cols, rows int) []*ebiten.Image {
	imagesList := []*ebiten.Image{}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			spriteX := x * spriteWidth
			spriteY := y * spriteHeight
			imageTemp := getSubImage(orig, spriteX, spriteY, spriteWidth, spriteHeight)
			imagesList = append(imagesList, imageTemp)
		}
	}
	return imagesList
}

func clampInt(min, max, test int) int {
	if test > max {
		return max
	} else if test < min {
		return min
	} else {
		return test
	}
}

// generic functions
func clamp[T int64 | float64 | int | int32 | float32](min, max, test T) T {
	if test > max {
		return max
	} else if test < min {
		return min
	} else {
		return test
	}
}

func abs[T int | float64 | int64](input T) T {
	if input < T(0) {
		return T(-1) * input
	} else {
		return input
	}
}

func DrawImageAt(screen *ebiten.Image, image *ebiten.Image, x int, y int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(image, op)
}

func copyAndStretchImage(orig *ebiten.Image, newX, newY int) *ebiten.Image {
	outputImage := ebiten.NewImage(newX, newY)
	oldX, oldY := orig.Bounds().Size().X, orig.Bounds().Size().Y
	//fmt.Println(newX)
	factorX := float64(newX) / float64(oldX)
	factorY := float64(newY) / float64(oldY)
	//factorX = 5
	//factorY = 5
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(factorX, factorY)
	outputImage.DrawImage(orig, op)
	return outputImage

}

func FlipHorizontal(source *ebiten.Image) *ebiten.Image {
	// https://ebitengine.org/en/tour/geom.html
	width := source.Bounds().Dx()
	height := source.Bounds().Dy()
	result := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)                 //multiply BRC x value of image by -1
	op.GeoM.Translate(float64(width), 0) // shift image left by its width
	result.DrawImage(source, op)         //apply geometry matrix
	return result
}

func FlipImageXorY(source *ebiten.Image, horizontal bool) *ebiten.Image {
	// https://ebitengine.org/en/tour/geom.html
	width := source.Bounds().Dx()
	height := source.Bounds().Dy()
	result := ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}
	if horizontal {
		// flip horizontal
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(width), 0)
	} else {
		// flip vertically
		op.GeoM.Scale(1, -1)
		op.GeoM.Translate(0, float64(height))
	}

	result.DrawImage(source, op) //apply geometry matrix
	return result
}

type rect struct {
	x      int
	y      int
	width  int
	height int
}

func (r *rect) Clone() *rect {
	r2 := &rect{}
	r2.x = r.x
	r2.y = r.y
	r2.width = r.width
	r2.height = r.height
	return r2
}

func collideRect(rect1, rect2 rect) bool {
	r1x2 := rect1.x + rect1.width
	r1y2 := rect1.y + rect1.height
	r2x2 := rect2.x + rect2.width
	r2y2 := rect2.y + rect2.height

	if rect1.x > r2x2 || rect2.x > r1x2 ||
		rect1.y > r2y2 || rect2.y > r1y2 {
		return false
	} else {
		return true
	}

}

func withinRange(known, test, radius int) bool {
	if abs(known-test) < radius || abs(test-known) < radius {
		return true
	} else {
		return false
	}
}

func attenuate(principal, amountToDecrease float32) float32 {
	if principal > 0 {
		output := principal - amountToDecrease
		if output < 0 {
			return 0.0
		} else {
			return output
		}
	} else if principal < 0 {
		output := principal + amountToDecrease
		if output > 0 {
			return 0.0
		} else {
			return output
		}
	} else {
		return 0.0
	}
}
