package logic

import (
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/server/common"
	models "sgserver/server/models"
	"sgserver/server/web/model"
	"sgserver/utils"
	"time"
)

var DefaultAccountLogic = &AccountLogic{}

type AccountLogic struct {
}

func (l AccountLogic) Register(rq *model.RegisterReq) error {

	username := rq.Username
	user := &models.User{}
	get, err := db.Engine.Table(user).Where("username=?", username).Get(user)
	if err != nil {
		log.Println("注册查询用户失败")
		return common.New(constant.DBError, "数据库异常")
	}
	if get {
		//已存在用户
		return common.New(constant.UserExist, "用户已存在")
	} else {
		user.Mtime = time.Now()
		user.Ctime = time.Now()
		user.Username = rq.Username
		user.Passcode = utils.RandSeq(6)
		user.Passwd = utils.Password(rq.Password, user.Passcode)
		user.Hardware = rq.Hardware
		_, err := db.Engine.Table(user).Insert(user)
		if err != nil {
			log.Println("用户入库失败", err)
			return common.New(constant.DBError, "数据库异常")
		}
		return nil
	}
}
