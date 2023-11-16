package logic

import (
	"log"
	"math/rand"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/global"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"sync"
	"time"
	"xorm.io/xorm"
)

var RoleCityService = &roleCityService{
	posRC:  make(map[int]*data.MapRoleCity),
	roleRC: make(map[int][]*data.MapRoleCity),
}

type roleCityService struct {
	//key 位置 posID
	posRC map[int]*data.MapRoleCity
	//key 角色id
	roleRC map[int][]*data.MapRoleCity
	mutex  sync.RWMutex
}

func (r *roleCityService) Load() {
	//查询所有角色城池
	dbRC := make(map[int]*data.MapRoleCity)
	db.Engine.Find(dbRC)

	for _, v := range dbRC {
		posId := global.ToPosition(v.X, v.Y)
		r.posRC[posId] = v
		_, ok := r.roleRC[v.RId]
		if !ok {
			r.roleRC[v.RId] = make([]*data.MapRoleCity, 0)
		} else {
			r.roleRC[v.RId] = append(r.roleRC[v.RId], v)
		}
	}

}

func (r *roleCityService) InitCity(rid int, nickname string, req *net.WsMsgReq) error {
	//根据用户ID查询对应的游戏角色..
	roleCity := &data.MapRoleCity{}
	get, err := db.Engine.Table(roleCity).Where("rid=?", rid).Get(roleCity)
	if err != nil {
		log.Println("查询角色城池出错", err)
		return common.New(constant.DBError, "查询数据库rid出错")
	}
	if get {
		return nil
	} else {
		//城池是否能在这个坐标创建, 需要判断. 系统城池/玩家城池 五格之内不能有玩家.
		//系统城池
		session := req.Context.Get("dbSession")
		for {
			//初始化
			roleCity.X = rand.Intn(global.MapWith)
			roleCity.Y = rand.Intn(global.MapHeight)
			if IsCanBuild(roleCity.X, roleCity.Y) {
				roleCity.RId = rid
				roleCity.Name = nickname
				roleCity.CurDurable = gameConfig.Base.City.Durable
				roleCity.CreatedAt = time.Now()
				roleCity.IsMain = 1
				if session != nil {
					_, err = session.(*xorm.Session).Table(roleCity).Insert(roleCity)
				} else {
					_, err = db.Engine.Table(roleCity).Insert(roleCity)
				}
				if err != nil {
					log.Println("城池初始化失败", err)
					return common.New(constant.DBError, "插入城池初始化信息失败")
				}
				//初始化城池设施
				if err := CityFacilityService.TryCreate(roleCity.CityId, rid, req); err != nil {
					log.Println("城池设施初始化失败", err)
					return common.New(err.(common.MyError).Code(), err.Error())
				}
				break
			}

		}
	}
	return nil
}

func IsCanBuild(x int, y int) bool {
	confs := gameConfig.MapRes.Confs
	pIndex := global.ToPosition(x, y)
	_, ok := confs[pIndex]
	if !ok {
		return false
	}
	sysBuild := gameConfig.MapRes.SysBuild
	for _, v := range sysBuild {
		if v.Type == gameConfig.MapBuildSysFortress {
			if x >= v.X-5 &&
				x <= v.X+5 &&
				y >= v.Y-5 &&
				y <= v.Y+5 {
				return false
			}
		}

	}

	return true
}

func (r *roleCityService) GetRoleCity(rid int) ([]model.MapRoleCity, error) {
	citys := make([]data.MapRoleCity, 0)
	city := &data.MapRoleCity{}
	err := db.Engine.Table(city).Where("rid=?", rid).Find(&citys)

	modelCitys := make([]model.MapRoleCity, 0)
	if err != nil {
		log.Println("查询角色城池出错", err)
		return modelCitys, err
	}
	for _, v := range citys {
		modelCitys = append(modelCitys, v.ToModel().(model.MapRoleCity))
	}
	return modelCitys, nil
}

func (r *roleCityService) ScanBlock(req *model.ScanBlockReq) ([]model.MapRoleCity, error) {
	x := req.X
	y := req.Y
	length := req.Length
	var mrcs []model.MapRoleCity
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return mrcs, nil
	}
	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)
	for i := x - length; i <= maxX; i++ {
		for j := y - length; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			mrb, ok := r.posRC[posId]
			if ok {
				mrcs = append(mrcs, mrb.ToModel().(model.MapRoleCity))
			}

		}
	}
	return mrcs, nil
}
