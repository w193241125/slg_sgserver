package net

import (
	"encoding/json"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"log"
	"sgserver/utils"
	"sync"
)

// WsServer websocket 服务
type WsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgRsp //通信管道...
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewWsServer(wsConn *websocket.Conn) *WsServer {
	return &WsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 1000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

func (w *WsServer) Router(router *Router) {
	w.router = router
}

func (w *WsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

func (w *WsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	return w.property[key], nil
}

func (w *WsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}

func (w *WsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()

}

func (w *WsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

func (w *WsServer) Start() {
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

// Write 写消息
func (w *WsServer) Write(msg *WsMsgRsp) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		//有加密
		key := secretKey.(string)
		//数据做加密
		data, _ = utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		//数据压缩
		if data, err := utils.Zip(data); err != nil {
			err := w.wsConn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

// 监控是否需要写消息
func (w *WsServer) writeMsgLoop() {

	for {
		select {
		case msg := <-w.outChan:
			fmt.Println(msg)
			w.Write(msg)
		}
	}
}

func (w *WsServer) readMsgLoop() {
	//先读到客户端发送过来的数据, 处理后, 再回消息
	// 经过路由, 实际处理程序
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("接收消息出错:", err)
			break
		}
		fmt.Println(data)
		//前端发送过来的消息就是JSON格式
		//1. data 解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压出错, 非法格式: ", err)
			continue
		}
		//2. 前端消息是加密的,需要解密
		secretKey, err := w.GetProperty("secretKey")
		if err == nil {
			key := secretKey.(string)
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println(err)
				//w.Handshake()
			} else {
				data = d
			}
		}
		//3. data 转为body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误, 非法格式:", err)
		} else {
			//获取到前端传递的数据, 去具体业务处理.
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			w.router.Run(req, rsp)
			w.outChan <- rsp
		}
	}
	w.Close()
	//rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	//w.outChan <- rsp
}

func (w *WsServer) Close() {
	_ = w.wsConn.Close()
}
