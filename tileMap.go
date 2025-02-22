package main

import (
	"fmt"
	"image/color"
	"log"
	"path"
	"slices"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	//mapRows        int    = 20
	//mapCols        int    = 20
	blankCellValue                   int     = 0
	brickImage0                      string  = "stoneWall1.png"
	brickImage1                      string  = "groundBrown.png"
	TM_SPRITESHEET_FILE              string  = "tileCT.png"
	TM_SPRITESHEET_FILE2             string  = "tilesB.png"
	TM_LAVA_IMAGES                   string  = "decorLava50.png"
	TM_WATER_IMAGES                  string  = "waterOverlay.png"
	TM_DEFAULT_MAP_FILENAME          string  = "map0.csv"
	TM_MAP_FILENAME_BASE             string  = "map"
	TM_MAP_FILENAME_END              string  = ".csv"
	TM_LOAD_MAP_ON_STARTUP           bool    = true
	TM_VIEWPORT_SHIFT_EDGE_RANGE     int     = 100
	TM_GRAVITY_ENABLED               bool    = true
	TM_GRAVITY_AMOUNT                float32 = 0.5
	TM_SPRITE_SIZE                   int     = 50
	TM_SPRITESHEET_ROWS              int     = 4
	TM_CREATE_BLANK_MAP_IF_NOT_EXIST bool    = true
	TM_CULLING_DISTANCE_TILES                = 10
	TM_LAVA_TILE_ID                          = 3
	TM_WATER_TILE_ID                         = 6

	TM_ANIMATED_TILE_FRAME_MAX = 4
	TM_ANIMATED_TILE_TICK_MAX  = 10
)

var (
	// become larger as player moves down and right
	worldOffsetX int = 0
	worldOffsetY int = 0
)

type TileMap struct {
	game          *Game
	tileSize      int
	rows          int
	cols          int
	tileData      [mapRows][mapCols]int
	images        []*ebiten.Image
	imagesLava    []*ebiten.Image
	imagesWater   []*ebiten.Image
	screenX       int
	screenY       int
	assetID       int
	tileKindMax   int
	runPan        bool
	filename_base string
	solidTiles    []int
	cullingRegion *CullingRegion

	currentAnimationTick  int
	currentAnimationFrame int
}

type CullingRegion struct {
	tileX1 int
	tileY1 int
	tileX2 int
	tileY2 int
}

type RectPointCollisionData struct {
	tlc bool
	top bool
	trc bool
	lef bool
	rig bool
	blc bool
	dow bool
	brc bool
}

type SideCollisionData struct {
	up    bool
	down  bool
	left  bool
	right bool
}

func NewTileMap(game *Game) *TileMap {
	tm := &TileMap{}
	tm.game = game
	tm.rows = mapRows
	tm.cols = mapCols
	tm.tileSize = GAME_TILE_SIZE
	tm.solidTiles = []int{0, 6, 14, 16, 21}
	//tm.tileData = initBlankGrid()
	tm.initTileMapImages()
	tm.initAnimatedTileImages()
	tm.tileKindMax = len(tm.images)
	tm.assetID = 0
	tm.filename_base = TM_MAP_FILENAME_BASE
	var err error
	if TM_LOAD_MAP_ON_STARTUP {
		err = tm.loadCurrentLevelMapFromFile()
		if err != nil {
			tm.tileData = initBlankGrid()
		}
	} else {
		tm.tileData = initBlankGrid()
	}
	return tm
}

func (tm *TileMap) updateCullingRegion() {
	if nil == tm.cullingRegion {
		tm.cullingRegion = &CullingRegion{}

	}

	HalfVP := tm.tileSize * TM_CULLING_DISTANCE_TILES / 2

	vpCenterX := (worldOffsetX + HalfVP) / tm.tileSize
	vpCenterY := (worldOffsetY + HalfVP) / tm.tileSize

	tm.cullingRegion.tileX1 = clamp(0, mapCols, vpCenterX-TM_CULLING_DISTANCE_TILES)
	tm.cullingRegion.tileX2 = clamp(0, mapCols, vpCenterX+TM_CULLING_DISTANCE_TILES)

	tm.cullingRegion.tileY1 = clamp(0, mapRows, vpCenterY-TM_CULLING_DISTANCE_TILES)
	tm.cullingRegion.tileY2 = clamp(0, mapRows, vpCenterY+TM_CULLING_DISTANCE_TILES)
}

func (pum *TileMap) getDataFileURL() string {
	filename := pum.filename_base + strconv.Itoa(pum.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	fmt.Println("Map file path is ", URL)
	return URL
}

func (tm *TileMap) getRectPointCollisionData(collider rect) *RectPointCollisionData {
	x2 := collider.x + collider.width
	y2 := collider.y + collider.height
	xm := collider.x + collider.width/2
	ym := collider.y + collider.height/2
	x1 := collider.x
	y1 := collider.y
	cd := &RectPointCollisionData{}
	cd.tlc = tm.pointCollidedWithSolidTile(x1, y1)
	cd.top = tm.pointCollidedWithSolidTile(xm, y1)
	cd.trc = tm.pointCollidedWithSolidTile(x2, y1)

	cd.lef = tm.pointCollidedWithSolidTile(x1, ym)
	//cd.tlc = tm.pointCollidedWithSolidTile(x1,y1)
	cd.rig = tm.pointCollidedWithSolidTile(x2, ym)

	cd.blc = tm.pointCollidedWithSolidTile(x1, y2)
	cd.dow = tm.pointCollidedWithSolidTile(xm, y2)
	cd.brc = tm.pointCollidedWithSolidTile(x2, y2)
	return cd

}

func (tm *TileMap) getSideCollisionData(collider rect) *SideCollisionData {
	pointColl := tm.getRectPointCollisionData(collider)
	// true if direction collides
	cdir := &SideCollisionData{}
	if (pointColl.tlc && pointColl.top) || (pointColl.top && pointColl.trc) {
		cdir.up = true
	} else {
		cdir.up = false
	}
	if (pointColl.blc && pointColl.dow) || (pointColl.dow && pointColl.brc) {
		cdir.down = true
	} else {
		cdir.down = false
	}
	if (pointColl.tlc && pointColl.lef) || (pointColl.lef && pointColl.tlc) {
		cdir.left = true
	} else {
		cdir.left = false
	}
	if (pointColl.trc && pointColl.rig) || (pointColl.rig && pointColl.brc) {
		cdir.right = true
	} else {
		cdir.right = false
	}

	return cdir

}

func (tm *TileMap) pointCollidedWithSolidTile(worldX, worldY int) bool {
	gridX := worldX / tm.tileSize
	gridY := worldY / tm.tileSize

	gridSizeX := len(tm.tileData[0])
	gridSizeY := len(tm.tileData)
	if gridX < 0 || gridX >= gridSizeX || gridY < 0 || gridY >= gridSizeY {
		return true
	}
	if kind := tm.tileData[gridY][gridX]; slices.Contains(tm.solidTiles, kind) {
		return false
	} else {
		return true
	}

}

func (tm *TileMap) pointCollidedWithGivenTileKind(worldX, worldY, kind int) bool {
	gridX := worldX / tm.tileSize
	gridY := worldY / tm.tileSize

	gridSizeX := len(tm.tileData[0])
	gridSizeY := len(tm.tileData)
	if gridX < 0 || gridX >= gridSizeX || gridY < 0 || gridY >= gridSizeY {
		return false
	}
	if kind == tm.tileData[gridY][gridX] {
		return true
	} else {
		return false
	}

}
func (tm *TileMap) solidUnderPlayer(distbelowFeet int) bool {
	//true if solid tile
	prect := tm.game.player.getWorldColliderRect()
	solid := tm.pointCollidedWithSolidTile(prect.x+(prect.width/2), prect.y+prect.height+distbelowFeet)
	//fmt.Println(solid)
	return solid

}

func (tm *TileMap) saveMapToFile() {
	name := tm.getDataFileURL()
	writeMapToFile(tm.tileData, name)
}
func (tm *TileMap) loadCurrentLevelMapFromFile() error {
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := tm.getDataFileURL()
	var err error
	td, err := loadMapFromFile(name)
	if err != nil {
		if TM_CREATE_BLANK_MAP_IF_NOT_EXIST {
			tm.tileData = initBlankGrid()
		}
		return err
	} else {
		tm.tileData = td
		return nil
	}
}

func (tm *TileMap) initTileMapImages_1() {
	imageDir := path.Join(subdir, brickImage0)
	//fmt.Println(imageDir)
	var rawImage *ebiten.Image
	var err error
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	//stretchedImage = copyAndStretchImage(rawImage, playerWidth, playerHeight)

	if err != nil {
		log.Fatal(err)
	}
	tm.images = []*ebiten.Image{}
	//tm.images := new(*ebiten.Image)
	blankImage := ebiten.NewImage(tm.tileSize, tm.tileSize)
	blankImage.Fill(color.RGBA{0x3f, 0x4f, 0x5f, 0x3f})
	tm.images = append(tm.images, blankImage)
	tm.images = append(tm.images, rawImage)

	imageDir = path.Join(subdir, brickImage1)
	rawImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	tm.images = append(tm.images, rawImage)
}

func (tm *TileMap) initTileMapImages() {
	imageDir := path.Join(subdir, TM_SPRITESHEET_FILE)
	//fmt.Println(imageDir)
	var spriteSheetImage *ebiten.Image
	var err error
	spriteSheetImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	//stretchedImage = copyAndStretchImage(rawImage, playerWidth, playerHeight)

	if err != nil {
		log.Fatal(err)
	}
	tm.images = cutSpriteSheet(spriteSheetImage, TM_SPRITE_SIZE,
		TM_SPRITE_SIZE, TM_SPRITESHEET_ROWS, TM_SPRITESHEET_ROWS)

	//2nd sheet
	imageDir = path.Join(subdir, TM_SPRITESHEET_FILE2)
	spriteSheetImage, _, err = ebitenutil.NewImageFromFile(imageDir)
	spriteSheet2 := cutSpriteSheet(spriteSheetImage, TM_SPRITE_SIZE,
		TM_SPRITE_SIZE, TM_SPRITESHEET_ROWS, TM_SPRITESHEET_ROWS)
	tm.images = append(tm.images, spriteSheet2...)

}

func (tm *TileMap) initAnimatedTileImages() {
	// lava images
	imageDir := path.Join(subdir, TM_LAVA_IMAGES)

	var spriteSheetImage *ebiten.Image
	var err error
	spriteSheetImage, _, err = ebitenutil.NewImageFromFile(imageDir)

	if err != nil {
		log.Fatal(err)
	}
	tm.imagesLava = cutSpriteSheet(spriteSheetImage, 50, 50, 5, 1)
	// water images
	imageDir = path.Join(subdir, TM_WATER_IMAGES)

	spriteSheetImage, _, err = ebitenutil.NewImageFromFile(imageDir)

	if err != nil {
		log.Fatal(err)
	}
	tm.imagesWater = cutSpriteSheet(spriteSheetImage, 50, 50, 5, 1)

}

func initBlankGrid() [mapRows][mapCols]int {
	outerArray := [mapRows][mapCols]int{}
	for y := 0; y < mapRows; y++ {
		for x := 0; x < mapCols; x++ {
			outerArray[y][x] = blankCellValue
		}
	}
	return outerArray

}

func (tm *TileMap) fillWithTile(kind int) {
	//outerArray := [mapRows][mapCols]int{}
	for y := 0; y < mapRows; y++ {
		for x := 0; x < mapCols; x++ {
			tm.tileData[y][x] = kind
		}
	}

}

func (tm *TileMap) cycleTileOverlay() {
	if tm.currentAnimationTick < TM_ANIMATED_TILE_TICK_MAX {
		tm.currentAnimationTick++
	} else {
		tm.currentAnimationTick = 0
		if tm.currentAnimationFrame < TM_ANIMATED_TILE_FRAME_MAX {
			tm.currentAnimationFrame++
		} else {
			tm.currentAnimationFrame = 0
		}
	}

}

func (tm *TileMap) Draw(screen *ebiten.Image) {
	startGX := tm.cullingRegion.tileX1
	startGY := tm.cullingRegion.tileY1
	endGX := tm.cullingRegion.tileX2
	endGY := tm.cullingRegion.tileY2
	for y := startGY; y < endGY; y++ {
		for x := startGX; x < endGX; x++ {
			cellValue := tm.tileData[y][x]
			_ = cellValue
			tileScreenX := (x * tm.tileSize) - worldOffsetX
			tileScreenY := (y * tm.tileSize) - worldOffsetY
			DrawImageAt(screen, tm.images[cellValue], tileScreenX, tileScreenY)
			if cellValue == TM_LAVA_TILE_ID {
				DrawImageAt(screen, tm.imagesLava[tm.currentAnimationFrame], tileScreenX, tileScreenY)
			} else if cellValue == TM_WATER_TILE_ID {
				DrawImageAt(screen, tm.imagesWater[tm.currentAnimationFrame], tileScreenX, tileScreenY)
			}
		}
	}

}

func (tm *TileMap) Update() {
	tm.cycleTileOverlay()
	tm.shiftViewportToFollowPlayer()
	tm.updateCullingRegion()

}

func (tm *TileMap) shiftViewportToFollowPlayer() {
	plr := tm.game.player.getColliderRect() //player screen loc
	x2 := plr.x + plr.width
	y2 := plr.y + plr.height
	//shiftAmount := defaultSpeed
	if plr.y < TM_VIEWPORT_SHIFT_EDGE_RANGE {
		worldOffsetY -= defaultSpeed

	}
	if y2 > tm.game.screenHeight-TM_VIEWPORT_SHIFT_EDGE_RANGE {
		worldOffsetY += defaultSpeed

	}

	panSpeedX := defaultSpeed
	if tm.runPan {
		panSpeedX += PL_RUN_BOOST
	}
	tm.runPan = false

	if plr.x < TM_VIEWPORT_SHIFT_EDGE_RANGE {
		worldOffsetX -= panSpeedX

	}
	if x2 > tm.game.screenWidth-TM_VIEWPORT_SHIFT_EDGE_RANGE {
		worldOffsetX += panSpeedX

	}

}

func (tm *TileMap) AddInstanceToGrid(tileX, tileY, assetID int) {
	if tileX >= mapCols || tileY >= mapRows ||
		tileX < 0 || tileY < 0 {
		fmt.Println("Can't put a tile there.")
		return
	}
	fmt.Printf("tile value %d  \n ", tm.tileData[tileY][tileX])

	fmt.Printf("set tile %d %d \n ", tileX, tileY)
	if tm.validatAssetID(assetID) {
		tm.assetID = assetID
	}
	tm.tileData[tileY][tileX] = tm.assetID
	fmt.Printf("tile value %d  \n ", tm.tileData[tileY][tileX])

}
func (tm *TileMap) validatAssetID(kind int) bool {
	if kind < len(tm.images) && kind > -1 {
		return true
	} else {
		return false
	}

}
func (tm *TileMap) CycleAssetKind(direction int) {
	propAssetID := tm.assetID + direction
	isValid := tm.validatAssetID(propAssetID)
	if isValid {
		tm.assetID = propAssetID
	}

	fmt.Println("Selected tile ", tm.assetID)

}
func (tm *TileMap) getAssetID() int {

	return tm.assetID

}

func (tm *TileMap) setAssetID(assetID int) {

	if assetID < len(tm.images) && assetID >= 0 {
		tm.assetID = assetID
	}

}
