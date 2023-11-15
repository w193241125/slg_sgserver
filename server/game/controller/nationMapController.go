package controller

import (
	"sgserver/constant"
	"sgserver/net"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/model"
)

var DefaultNationMapController = &nationMapController{}

type nationMapController struct {
}

func (n *nationMapController) Router(router *net.Router) {
	g := router.Group("nationMap")
	g.AddRouter("config", n.config)
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