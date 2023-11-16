package controller

import (
	"github.com/mitchellh/mapstructure"
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/logic"
	"sgserver/server/game/middleware"
	"sgserver/server/game/model"
)

var DefaultNationMapController = &nationMapController{}

type nationMapController struct {
}

func (n *nationMapController) Router(router *net.Router) {
	g := router.Group("nationMap")
	g.Use(middleware.Log())
	g.AddRouter("config", n.config)
	g.AddRouter("scanBlock", n.scanBlock)
}
func (n *nationMapController) config(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//reqObj := &model.ConfigReq{}
	rspObj := &model.ConfigRsp{}
	cfgs := gameConfig.MapBuildConf.Cfg

	rspObj.Confs = make([]model.Conf, len(cfgs))
	for index, v := range cfgs {
		rspObj.Confs[index].Type = v.Type
		rspObj.Confs[index].Name = v.Name
		rspObj.Confs[index].Level = v.Level
		rspObj.Confs[index].Defender = v.Defender
		rspObj.Confs[index].Durable = v.Durable
		rspObj.Confs[index].Grain = v.Grain
		rspObj.Confs[index].Iron = v.Iron
		rspObj.Confs[index].Stone = v.Stone
		rspObj.Confs[index].Wood = v.Wood
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}

func (n *nationMapController) scanBlock(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.ScanBlockReq{}
	rspObj := &model.ScanBlockRsp{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK

	mrb, err := logic.RoleBuildService.ScanBlock(reqObj)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.MRBuilds = mrb
	//扫描角色城池
	mrc, err := logic.RoleCityService.ScanBlock(reqObj)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.MCBuilds = mrc
	//扫描玩家军队
	armys, err := logic.ArmyService.ScanBlock(reqObj)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Armys = armys
	rsp.Body.Msg = rspObj
}
