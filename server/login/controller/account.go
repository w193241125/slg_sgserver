package controller

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"sgserver/constant"
	"sgserver/db"
	"sgserver/net"
	"sgserver/server/login/model"
	"sgserver/server/login/proto"
	"sgserver/utils"
	"time"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(router *net.Router) {
	g := router.Group("account")
	g.AddRouter("login", a.login)

}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	/*
		1.用户名 密码 硬件ID
		2. 查user表, 得到数据,
		3. 密码对比,正确->成功.
		4. 保存用户登录记录
		5. 保存用户最后登录信息.
		6. 生成 session 返回给客户端 , jwt 生成加密字符串的加密算法
		7. 客户端发起需要登录的行为时, 判断用户是否合法.
	*/
	loginReq := &proto.LoginReq{}
	loginRes := &proto.LoginRsp{}
	mapstructure.Decode(req.Body.Msg, loginReq)
	user := &model.User{}
	get, err := db.Engine.Table(user).Where("username=?", loginReq.Username).Get(user)
	if err != nil {
		log.Println("用户表查询出错", err)
		return
	}
	fmt.Println(loginReq.Username)
	fmt.Println(get)
	if !get {
		rsp.Body.Code = constant.UserNotExist
		return
	}
	pwd := utils.Password(loginReq.Password, user.Passcode)
	if pwd != user.Passwd {
		rsp.Body.Code = constant.PwdIncorrect
		return
	}
	//jwt A.B.C 三部分组成, A 定义加密算法, B 定义放入的数据, C 根据密钥 + A和B 生成加密字符串.
	token, _ := utils.Award(user.UId)
	rsp.Body.Code = constant.OK

	loginRes.UId = user.UId
	loginRes.Username = user.Username
	loginRes.Session = token
	loginRes.Password = ""
	rsp.Body.Msg = loginRes
	//保存用户登录记录
	ul := &model.LoginHistory{
		UId: user.UId, CTime: time.Now(), Ip: loginReq.Ip,
		Hardware: loginReq.Hardware, State: model.Login,
	}
	_, err = db.Engine.Table(ul).Insert(ul)
	if err != nil {
		log.Println("记录登录历史出错", err)
		return
	}

	//最后一次登录状态记录
	ll := &model.LoginLast{}
	println(user.UId)
	get, _ = db.Engine.Table(ll).Where("uid=?", user.UId).Get(ll)
	fmt.Println(get)
	if get {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		db.Engine.Table(ll).Update(ll)
	} else {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		ll.UId = user.UId
		_, err := db.Engine.Table(ll).Insert(ll)
		if err != nil {
			log.Println(err)
			return
		}
	}

	// 缓存一下, 当前用户与当前ws的连接
	net.Mgr.UserLogin(req.Conn, user.UId, token)
}
