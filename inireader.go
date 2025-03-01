package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	READINI_SEPERATOR = '='
	READINI_COMMENT   = '#'
	READINI_NEWLINE   = "\n"
)

func readIni(filePath string) (map[string]string, error) {
	outmap := map[string]string{}
	var err error = nil

	stringList := []string{}
	type innerList []int

	file, err := os.Open(filePath)
	defer file.Close()
	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()
		if line != nil && len(line) > 0 {
			//fmt.Printf("ReadLine: %q\n", line)
			lineAsString := string(line)
			stringList = append(stringList, lineAsString)
		} else {
			break
		}
		if err != nil {
			log.Println(err)
			return outmap, err
		}
	}
	if err != nil {
		log.Println(err)
		return outmap, err
	}
	rows := len(stringList)
	if rows == 0 {
		return outmap, nil
	}

	for _, v := range stringList {
		LV := len(v)
		if rune(v[0]) == READINI_COMMENT {
			continue
		} else if strings.Contains(v, string(READINI_SEPERATOR)) &&
			v[0] != READINI_SEPERATOR &&
			v[LV-1] != READINI_SEPERATOR {
			sepPosition := strings.Index(v, string(READINI_SEPERATOR))
			var substr1, substr2 string
			substr1 = string(v[0 : sepPosition-1])  //pos 0 to eq-1
			substr2 = string(v[sepPosition+1 : LV]) // pos eq+1 to length-1
			outmap[substr1] = substr2

		}
	}
	return outmap, err

}

func defaultDataMap() map[string]string {
	outmap := map[string]string{}
	outmap["width"] = "800"
	outmap["height"] = "600"
	outmap["console"] = "true"
	return outmap
}

func openOrCreateDefaultIni(filePath string) (map[string]string, error) {
	dataMap, err := readIni(filePath)
	if err != nil {
		log.Println("inireader failed to open file ", filePath)
		dataMap = defaultDataMap()
		errW := writeHashmapToIni(dataMap, filePath)
		if errW != nil {
			log.Println("inireader failed to write default file ", filePath)
			err = errW
		}
	}
	return dataMap, err
}

func writeHashmapToIni(dataMap map[string]string, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
		log.Fatal(fmt.Sprintf("writeHashmapToIni failed to write file: %s/n", filePath))
		return err
	}
	defer file.Close() // close file at end of this function

	var sb strings.Builder
	writer := bufio.NewWriter(file)
	for k, v := range dataMap {
		sb.Reset()
		sb.WriteString(k)
		sb.WriteRune(READINI_SEPERATOR)
		sb.WriteString(v)
		sb.WriteString(READINI_NEWLINE)

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

func pprintMap(dataMap map[string]string) {

	for k, v := range dataMap {

		fmt.Printf("%s = %s\n", k, v)
	}

}
