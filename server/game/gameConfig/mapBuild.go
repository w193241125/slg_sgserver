package gameConfig

import (
	"encoding/json"
	"os"
)

type cfg struct {
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	Grain    int    `json:"grain"`
	Wood     int    `json:"wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Durable  int    `json:"durable"`
	Defender int    `json:"defender"`
}

type mapBuildConf struct {
	Title  string `json:"title"`
	Cfg    []cfg  `json:"cfg"`
	cfgMap map[int8][]cfg
}

var MapBuildConf = &mapBuildConf{}

const mapBuildConfFile = "/conf/game/map_build.json"

func (m *mapBuildConf) Load() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + mapBuildConfFile

	lens := len(os.Args)
	if lens > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + mapBuildConfFile
		}
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		panic(err)
	}
}
