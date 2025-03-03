package main

import (
	"fmt"
	"log"
	"path"
	"strconv"
)

const (
	EM_BARNACLEFISH_W = 300
	EM_BARNACLEFISH_H = 600
)

func (ent *EntityManager) saveDataToFile() {
	name := ent.getDataFileURL()
	numericData := [][]int{}
	rows := len(ent.entityList)
	for i := 0; i < rows; i++ {
		entObj := ent.entityList[i]
		if entObj != nil {
			record := []int{entObj.kind, entObj.startGridX, entObj.startGridY, entObj.uid}
			numericData = append(numericData, record)
		}

	}
	if rows != 0 {

		write2DIntListToFile(numericData, name)
	} else {
		log.Println("Entitys: no data to write, ", name)
	}
}

func (ent *EntityManager) getDataFileURL() string {
	filename := ent.filename_base + strconv.Itoa(ent.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}

func (ent *EntityManager) loadDataFromFile() error {
	ent.entityList = []*Entity{}
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := ent.getDataFileURL()
	numericData, err := loadDataListFromFile(name)
	rows := len(numericData)
	if rows == 0 {
		log.Println("Entity loadDataFromFile no data to load")
		return nil
	}
	if err != nil {
		return err
	}
	//var entity *Entity
	for i := 0; i < EN_MAX_ENTITIES_AT_ONCE && i < rows; i++ {
		v := numericData[i]
		//entity = ent.AddInstanceToGrid(v[0], v[1], v[2])
		entityTemp := NewEntity(v[0], v[1], v[2])
		entityTemp.uid = v[3]
		ent.entityList = append(ent.entityList, entityTemp)
		//fmt.Println("added entity ")
	}
	return nil
}

func (ent *EntityManager) addEntity(kind, startGridX, startGridY int) *Entity {
	entityTemp := NewEntity(kind, startGridX, startGridY)
	ent.entityList = append(ent.entityList, entityTemp)
	return entityTemp
}

func (ent *EntityManager) removeEntityByID(uid int) bool {
	var found = false
	for i, v := range ent.entityList {
		if v.uid == uid {
			ent.entityList[i] = nil
			found = true
			break
		}
	}
	return found
}

func (ent *EntityManager) getUniqueUID() int {

	return 0
}

func (ent *EntityManager) AddInstanceToGrid(gridX, gridY, kind int) {
	var entity *Entity
	if len(ent.entityList) <= EN_MAX_ENTITIES_AT_ONCE {
		//uid := ent.getUniqueUID()
		entity = ent.createEntityInstance(gridX, gridY, kind)
		ent.entityList = append(ent.entityList, entity)
		log.Printf("Added Entity %d at %d, %d\n", kind, gridX, gridY)
	} else {
		log.Println("Failed to add Entity, no open slots")
	}

}

func (ent *EntityManager) createEntityInstance(gridX, gridY, kind int) *Entity {
	x := gridX
	y := gridY
	//uid := ent.getUniqueUID()
	entity := NewEntity(kind, x, y)
	entity.alive = true
	entity.uid = ent.getUniqueUID()
	entity.width = EN_SPRITE_W
	entity.height = EN_SPRITE_H
	if entity.kind == EM_BARNACLEFISH_TYPE {
		entity.width = EM_BARNACLEFISH_W
		entity.height = EM_BARNACLEFISH_H
	}
	return entity
}

func (ent *EntityManager) AddEntityToGrid(gridX, gridY, kind int) *Entity {
	var entity *Entity
	if len(ent.entityList) <= EN_MAX_ENTITIES_AT_ONCE {
		//x := gridX
		//y := gridY
		//uid := ent.getUniqueUID()
		entity = ent.createEntityInstance(gridX, gridY, kind)
		ent.entityList = append(ent.entityList, entity)
		log.Printf("Added Entity %d at %d, %d\n", kind, gridX, gridY)
	} else {
		log.Println("Failed to add Entity, no open slots")
	}
	return entity
}

func (tm *EntityManager) validatAssetID(kind int) bool {
	if kind < len(tm.images) && kind > -1 {
		return true
	} else {
		return false
	}

}

func (tm *EntityManager) CycleAssetKind(direction int) {
	propAssetID := tm.assetID + direction
	isValid := tm.validatAssetID(propAssetID)
	if isValid {
		tm.assetID = propAssetID
	}

	fmt.Println("Selected Entity ", tm.assetID)

}

func (tm *EntityManager) getAssetID() int {
	fmt.Println("EntityManager getAssetID", tm.assetID)
	return tm.assetID

}

func (tm *EntityManager) setAssetID(assetID int) {

	if assetID < EM_KIND_MAX && assetID >= 0 {
		tm.assetID = assetID
	}
	tm.assetID = assetID
	fmt.Println("EntityManager Selected entity type ", tm.assetID)

}
