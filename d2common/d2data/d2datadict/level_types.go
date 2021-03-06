package d2datadict

import (
	"log"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2common"
)

type LevelTypeRecord struct {
	Name      string
	ID        int
	Files     [32]string
	Beta      bool
	Act       int
	Expansion bool
}

var LevelTypes []LevelTypeRecord

func LoadLevelTypes(file []byte) {
	data := strings.Split(string(file), "\r\n")[1:]
	LevelTypes = make([]LevelTypeRecord, len(data))

	for i, j := 0, 0; i < len(data); i, j = i+1, j+1 {
		idx := -1
		inc := func() int {
			idx++
			return idx
		}

		if data[i] == "" {
			continue
		}

		parts := strings.Split(data[i], "\t")

		if parts[0] == "Expansion" {
			j--
			continue
		}

		LevelTypes[j].Name = parts[inc()]
		LevelTypes[j].ID = d2common.StringToInt(parts[inc()])

		for fileIdx := range LevelTypes[i].Files {
			LevelTypes[j].Files[fileIdx] = parts[inc()]
			if LevelTypes[j].Files[fileIdx] == "0" {
				LevelTypes[j].Files[fileIdx] = ""
			}
		}

		LevelTypes[j].Beta = parts[inc()] != "1"
		LevelTypes[j].Act = d2common.StringToInt(parts[inc()])
		LevelTypes[j].Expansion = parts[inc()] != "1"
	}
	log.Printf("Loaded %d LevelType records", len(LevelTypes))
}
