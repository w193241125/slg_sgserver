package general

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
)

type general struct {
	Title            string                `json:"title"`
	GArr             []generalDetail       `json:"list"` //武将属性
	GMap             map[int]generalDetail // CfgId 与武将的映射
	totalProbability int
}
type generalDetail struct {
	Name         string `json:"name"`
	CfgId        int    `json:"cfgId"`
	Force        int    `json:"force"`    //武力
	Strategy     int    `json:"strategy"` //策略
	Defense      int    `json:"defense"`  //防御
	Speed        int    `json:"speed"`    //速度
	Destroy      int    `json:"destroy"`  //破坏力
	ForceGrow    int    `json:"force_grow"`
	StrategyGrow int    `json:"strategy_grow"`
	DefenseGrow  int    `json:"defense_grow"`
	SpeedGrow    int    `json:"speed_grow"`
	DestroyGrow  int    `json:"destroy_grow"`
	Cost         int8   `json:"cost"`
	Probability  int    `json:"probability"`
	Star         int8   `json:"star"`
	Arms         []int  `json:"arms"`
	Camp         int8   `json:"camp"`
}

var General = &general{}
var generalFile = "/conf/game/general/general.json"

func (g *general) Load() {
	g.GArr = make([]generalDetail, 0)
	g.GMap = make(map[int]generalDetail)
	//获取当前文件路径
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//配置文件位置
	cf := currentDir + generalFile
	//打包后 程序参数加入配置文件路径
	length := len(os.Args)
	if length > 1 {
		path := os.Args[1]
		if path != "" {
			cf = path + generalFile
		}
	}
	data, err := os.ReadFile(cf)
	if err != nil {
		log.Println("武将配置读取失败")
		panic(err)
	}
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Println("武将配置格式定义失败")
		panic(err)
	}

	for _, v := range g.GArr {
		g.GMap[v.CfgId] = v
		g.totalProbability += v.Probability
	}
}

func (g *general) Rand() int {
	rate := rand.Intn(g.totalProbability)
	var cur = 0
	for _, v := range g.GArr {
		if rate >= cur && rate < cur+v.Probability {
			return v.CfgId
		}
		cur += v.Probability
	}
	return 0
}
