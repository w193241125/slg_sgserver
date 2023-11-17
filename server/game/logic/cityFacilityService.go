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

func (c *cityFacilityService) TryCreate(cid, rid int, req *net.WsMsgReq) error {
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
		return common.New(constant.DBError, "城池设施插入出错")
	}
	return nil
}
func (c *cityFacilityService) GetById(rid int) ([]*data.CityFacility, error) {
	cf := make([]*data.CityFacility, 0)
	ct := &data.CityFacility{}
	err := db.Engine.Table(ct).Where("rid=?", rid).Find(&cf)
	if err != nil {
		return cf, common.New(constant.DBError, "查询城池设施出错")
	}
	return cf, nil
}

func (c *cityFacilityService) GetYield(rid int) data.Yield {
	//查询表中设施, 获取到
	//设施不同, 去配置中查询匹配, 增加产量的设施, 木头 金钱
	//设施等级不同,产量也不一样.
	cfs, err := c.GetById(rid)
	var y data.Yield
	if err == nil {
		for _, v := range cfs {
			facilities := v.Facility()
			for _, fa := range facilities {
				//计算等级 资源的产出是不同的
				if fa.GetLevel() > 0 {
					//计算等级(不同等级,产量不同)
					values := gameConfig.FacilityConf.GetValues(fa.Type, fa.GetLevel())
					adds := gameConfig.FacilityConf.GetAdditions(fa.Type)
					for i, aType := range adds {
						if aType == gameConfig.TypeWood {
							y.Wood += values[i]
						} else if aType == gameConfig.TypeGrain {
							y.Grain += values[i]
						} else if aType == gameConfig.TypeIron {
							y.Iron += values[i]
						} else if aType == gameConfig.TypeStone {
							y.Stone += values[i]
						} else if aType == gameConfig.TypeTax {
							y.Gold += values[i]
						}
					}
				}
			}
		}
	}
	return y
}
