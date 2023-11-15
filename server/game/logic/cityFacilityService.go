package logic

import (
	"encoding/json"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/gameConfig"
	"sgserver/server/game/model/data"
	"xorm.io/xorm"
)

var CityFacilityService = &cityFacilityService{}

type cityFacilityService struct {
}

func (c cityFacilityService) TryCreate(cid, rid int, req *net.WsMsgReq) error {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(cf).Where("cityId=?", cid).Get(cf)
	if err != nil {
		log.Println("查询城池设施出错", err)
		return common.New(constant.DBError, "查询城池设施出错")
	}
	if ok {
		return nil
	}

	cf.RId = rid
	cf.CityId = cid
	list := gameConfig.FacilityConf.List
	facs := make([]data.Facility, len(list))
	for k, v := range list {
		fac := data.Facility{
			Name:         v.Name,
			Type:         v.Type,
			PrivateLevel: 0,
			UpTime:       0,
		}
		facs[k] = fac
	}
	dataJson, _ := json.Marshal(facs)
	cf.Facilities = string(dataJson)
	if session := req.Context.Get("dbSession"); session != nil {
		_, err = session.(*xorm.Session).Table(cf).Insert(cf)
	} else {
		_, err = db.Engine.Table(cf).Insert(cf)
	}

	if err != nil {
		log.Println("城池设施插入出错", err)
		return err
	}
	return nil
}
