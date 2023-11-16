package logic

import (
	"log"
	"sgserver/db"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/global"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sgserver/utils"
	"sync"
)

var RoleBuildService = &roleBuildService{
	posRB:  make(map[int]*data.MapRoleBuild),
	roleRB: make(map[int][]*data.MapRoleBuild),
}

type roleBuildService struct {
	//key 位置 posID
	posRB map[int]*data.MapRoleBuild
	//key 角色id
	roleRB map[int][]*data.MapRoleBuild
	mutex  sync.RWMutex
}

func (r *roleBuildService) Load() {
	//加载系统以及玩家建筑
	//首先需要判断数据库是否保存了系统建筑,没有就进行一个保存
	count, err := db.Engine.Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).Count(new(data.MapRoleBuild))
	if err != nil {
		return
	}
	if int64(len(gameConfig.MapRes.SysBuild)) != count {
		db.Engine.Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).Delete(new(data.MapRoleBuild))
		//证明系统数据库存储的系统建筑有问题
		for _, v := range gameConfig.MapRes.SysBuild {
			build := &data.MapRoleBuild{
				RId:   0,
				Type:  v.Type,
				X:     v.X,
				Y:     v.Y,
				Level: v.Level,
			}
			build.Init()
			_, err := db.Engine.InsertOne(build)
			if err != nil {
				log.Println("保存系统建筑失败", err)
			}
		}
	}
	//查询所有角色建筑
	dbRB := make(map[int]*data.MapRoleBuild)
	db.Engine.Find(dbRB)
	for _, v := range dbRB {
		posId := global.ToPosition(v.X, v.Y)
		r.posRB[posId] = v
		_, ok := r.roleRB[v.RId]
		if !ok {
			r.roleRB[v.RId] = make([]*data.MapRoleBuild, 0)
		} else {
			r.roleRB[v.RId] = append(r.roleRB[v.RId], v)
		}
	}

}

func (r *roleBuildService) GetBuilds(rid int) ([]model.MapRoleBuild, error) {
	builds := make([]data.MapRoleBuild, 0)
	build := &data.MapRoleBuild{}
	err := db.Engine.Table(build).Where("rid=?", rid).Find(&builds)

	modelBuilds := make([]model.MapRoleBuild, 0)
	if err != nil {
		log.Println("查询建筑出错", err)
		return modelBuilds, err
	}
	for _, v := range builds {
		modelBuilds = append(modelBuilds, v.ToModel().(model.MapRoleBuild))
	}
	return modelBuilds, nil
}

func (r *roleBuildService) ScanBlock(req *model.ScanBlockReq) ([]model.MapRoleBuild, error) {
	x := req.X
	y := req.Y
	length := req.Length
	var mrbs = make([]model.MapRoleBuild, 0)
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return mrbs, nil
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)
	for i := x - length; i <= maxX; i++ {
		for j := y - length; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			mrb, ok := r.posRB[posId]
			if ok {
				mrbs = append(mrbs, mrb.ToModel().(model.MapRoleBuild))
			}

		}
	}
	return mrbs, nil
}
