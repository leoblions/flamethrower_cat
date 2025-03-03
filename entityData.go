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

func (em *EntityManager) saveDataToFile() {
	name := em.getDataFileURL()
	numericData := [][]int{}
	rows := len(em.entityList)
	for i := 0; i < rows; i++ {
		entObj := em.entityList[i]
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

func (em *EntityManager) getDataFileURL() string {
	filename := em.filename_base + strconv.Itoa(em.game.level) + GAME_DATA_MATRIX_END
	URL := path.Join(GAME_LEVEL_DATA_DIR, filename)
	return URL
}

func (em *EntityManager) loadDataFromFile() error {
	em.entityList = []*Entity{}
	//writeMapToFile(tm.tileData, TM_DEFAULT_MAP_FILENAME)
	name := em.getDataFileURL()
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
		entityTemp := em.createUniqueEntity(v[1], v[2], v[0])
		//entityTemp := NewEntity(v[0], v[1], v[2])
		entityTemp.uid = v[3]
		em.entityList = append(em.entityList, entityTemp)
		//fmt.Println("added entity ")
	}
	return nil
}

func (em *EntityManager) addEntity(kind, startGridX, startGridY int) *Entity {
	entityTemp := NewEntity(kind, startGridX, startGridY)
	em.entityList = append(em.entityList, entityTemp)
	return entityTemp
}

func (em *EntityManager) removeEntityByID(uid int) bool {
	var found = false
	for i, v := range em.entityList {
		if v.uid == uid {
			em.entityList[i] = nil
			found = true
			break
		}
	}
	return found
}

func (em *EntityManager) getUniqueUID() int {

	return 0
}

func (em *EntityManager) AddInstanceToGrid(gridX, gridY, kind int) {
	var entity *Entity
	if len(em.entityList) <= EN_MAX_ENTITIES_AT_ONCE {
		//uid := ent.getUniqueUID()
		entity = em.createUniqueEntity(gridX, gridY, kind)
		em.entityList = append(em.entityList, entity)
		log.Printf("Added Entity %d at %d, %d\n", kind, gridX, gridY)
	} else {
		log.Println("Failed to add Entity, no open slots")
	}

}

func (em *EntityManager) createUniqueEntity(gridX, gridY, kind int) *Entity {

	entity := NewEntity(kind, gridX, gridY)

	entity.uid = em.getUniqueUID()

	return entity
}

func (em *EntityManager) AddEntityToGrid(gridX, gridY, kind int) *Entity {
	var entity *Entity
	if len(em.entityList) <= EN_MAX_ENTITIES_AT_ONCE {
		//x := gridX
		//y := gridY
		//uid := ent.getUniqueUID()
		entity = em.createUniqueEntity(gridX, gridY, kind)
		em.entityList = append(em.entityList, entity)
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
