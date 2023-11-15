package logic

import (
	"encoding/json"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/common"
	"sgserver/server/game/model"
	"sgserver/server/game/model/data"
	"sync"
	"xorm.io/xorm"
)

var RoleAttrService = &roleAttrService{
	attrs: make(map[int]*data.RoleAttribute),
}

type roleAttrService struct {
	attrs map[int]*data.RoleAttribute
	mutex sync.RWMutex
}

func (r *roleAttrService) TryCreate(rid int, req *net.WsMsgReq) error {
	//根据用户ID查询对应的游戏角色..
	role := &data.RoleAttribute{}
	get, err := db.Engine.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return common.New(constant.DBError, "查询数据库rid出错")
	}
	if get {
		r.mutex.Lock()
		r.attrs[rid] = role
		defer r.mutex.Unlock()
		return nil
	} else {
		role.RId = rid
		role.UnionId = 0
		role.ParentId = 0
		role.PosTags = ""
		if session := req.Context.Get("dbSession"); session != nil {
			_, err = session.(*xorm.Session).Table(role).Insert(role)
		} else {
			_, err = db.Engine.Table(role).Insert(role)
		}

		if err != nil {
			log.Println("插入初始角色属性失败", err)
			return common.New(constant.DBError, "插入初始角色属性失败")
		}
		r.mutex.Lock()
		defer r.mutex.Unlock()
		r.attrs[rid] = role

	}
	return nil
}

func (r *roleAttrService) GetTagList(rid int) ([]model.PosTag, error) {

	ra, ok := r.attrs[rid]
	if !ok {
		ra := &data.RoleAttribute{}
		var err error
		ok, err = db.Engine.Table(ra).Where("rid=?", rid).Get(ra)
		if err != nil {
			log.Println("getTagList err", err)
			return nil, common.New(constant.DBError, "数据库出错")
		}
	}

	posTags := make([]model.PosTag, 0)
	if ok {
		tags := ra.PosTags
		if tags != "" {
			err := json.Unmarshal([]byte(tags), &posTags)
			if err != nil {
				return nil, common.New(constant.DBError, "数据库错误")
			}

		}
	}
	return posTags, nil
}
