package gameConfig

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

var Skill skill

type skill struct {
	skills       []Conf
	skillConfMap map[int]Conf
	outline      outline
}

type trigger struct {
	Type int    `json:"type"`
	Des  string `json:"des"`
}

type triggerType struct {
	Des  string    `json:"des"`
	List []trigger `json:"list"`
}

type effect struct {
	Type   int    `json:"type"`
	Des    string `json:"des"`
	IsRate bool   `json:"isRate"`
}

type effectType struct {
	Des  string   `json:"des"`
	List []effect `json:"list"`
}

type target struct {
	Type int    `json:"type"`
	Des  string `json:"des"`
}

type targetType struct {
	Des  string   `json:"des"`
	List []target `json:"list"`
}

type outline struct {
	TriggerType triggerType `json:"trigger_type"` //触发类型
	EffectType  effectType  `json:"effect_type"`  //效果类型
	TargetType  targetType  `json:"target_type"`  //目标类型
}

type level struct {
	Probability int   `json:"probability"`  //发动概率
	EffectValue []int `json:"effect_value"` //效果值
	EffectRound []int `json:"effect_round"` //效果持续回合数
}

type Conf struct {
	CfgId         int     `json:"cfgId"`
	Name          string  `json:"name"`
	Trigger       int     `json:"trigger"` //发起类型
	Target        int     `json:"target"`  //目标类型
	Des           string  `json:"des"`
	Limit         int     `json:"limit"`          //可以被武将装备上限
	Arms          []int   `json:"arms"`           //可以装备的兵种
	IncludeEffect []int   `json:"include_effect"` //技能包括的效果
	Levels        []level `json:"levels"`
}

const skillFile = "/conf/game/skill/skill_outline.json"
const skillPath = "/conf/game/skill/"

func (s *skill) Load() {
	s.skills = make([]Conf, 0)
	s.skillConfMap = make(map[int]Conf)
	currentDir, _ := os.Getwd()
	cf := currentDir + skillFile
	cp := currentDir + skillPath
	if len(os.Args) > 1 {
		if path := os.Args[1]; path != "" {
			cf = path + skillFile
			cp = path + skillPath
		}
	}
	data, err := os.ReadFile(cf)
	if err != nil {
		log.Println("读取配置出错")
		panic(err)
	}
	err = json.Unmarshal(data, &s.outline)
	if err != nil {
		log.Println("技能配置格式定义失败")
		panic(err)
	}
	files, err := os.ReadDir(cp)
	if err != nil {
		log.Println("技能文件读取失败")
		panic(err)
	}
	for _, v := range files {
		if v.IsDir() {
			name := v.Name()
			dirPath := cp + name
			skillFiles, err := os.ReadDir(dirPath)
			if err != nil {
				log.Println("读取文件技能失败")
				panic(err)
			}
			for _, sv := range skillFiles {
				if sv.IsDir() {
					continue
				}
				fileJson := path.Join(dirPath, sv.Name())
				conf := Conf{}
				data, err = os.ReadFile(fileJson)
				if err != nil {
					log.Println(name + "技能文件格式错误")
					panic(err)
				}
				err := json.Unmarshal(data, &conf)
				if err != nil {
					log.Println(name + "技能文件格式错误")
					panic(err)
				}
				s.skills = append(s.skills, conf)
				s.skillConfMap[conf.CfgId] = conf
			}
		}
	}
}
