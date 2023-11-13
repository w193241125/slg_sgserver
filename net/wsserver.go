package net

import (
	"encoding/json"
	"errors"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"sgserver/utils"
	"sync"
	"time"
)

// WsServer websocket 服务
type WsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgRsp //通信管道...
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex
	needSecret   bool
}

var cid int64

func NewWsServer(wsConn *websocket.Conn, needSecret bool) *WsServer {
	s := &WsServer{
		wsConn:     wsConn,
		outChan:    make(chan *WsMsgRsp, 1000),
		property:   make(map[string]interface{}),
		Seq:        0,
		needSecret: needSecret,
	}

	cid++
	s.SetProperty("cid", cid)
	return s
}

func (w *WsServer) NeedSecret(needSecret bool) {
	w.needSecret = needSecret
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
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("property no found")
	}
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

// 监控是否需要写消息
func (w *WsServer) writeMsgLoop() {

	for {
		select {
		case msg := <-w.outChan:
			w.Write(msg.Body)
		}
	}
}

// Write 写消息
func (w *WsServer) Write(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		//有加密
		key := secretKey.(string)
		//数据做加密
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)

	}
	//数据压缩
	if data, err := utils.Zip(data); err == nil {
		err := w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Println("服务端写数据出错", err)
		}
		d, _ := json.Marshal(msg)
		log.Println("服务端写数据", string(d))
	}
}

func (w *WsServer) readMsgLoop() {
	//先读到客户端发送过来的数据, 处理后, 再回消息
	// 经过路由, 实际处理程序
	defer func() {
		if err := recover(); err != nil {
			log.Println("wsserver readmsgLoop 异常", err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("wsserver收消息出错:", err)
			break
		}
		//fmt.Println(data)
		//前端发送过来的消息就是JSON格式
		//1. data 解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压出错, 非法格式: ", err)
			continue
		}
		//2. 前端消息是加密的,需要解密
		//if w.needSecret {
		secretKey, err := w.GetProperty("secretKey")
		if err == nil {
			key := secretKey.(string)
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误，解密失败:", err)
				w.Handshake()
			} else {
				data = d
			}
		}
		//}

		//3. data 转为body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误, 非法格式:", err)
		} else {
			//获取到前端传递的数据, 去具体业务处理.
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			if req.Body.Name == "heartbeat" {
				//回心跳消息
				h := &Heartbeat{}
				mapstructure.Decode(req.Body.Msg, h)
				h.STime = time.Now().UnixNano() / 1e6
				rsp.Body.Msg = h
			} else {
				if w.router != nil {
					w.router.Run(req, rsp)
				}

			}

			w.outChan <- rsp
		}
	}
	w.Close()
}

func (w *WsServer) Close() {
	_ = w.wsConn.Close()
}

const HandshakeMsg = "handshake"

// Handshake 当客户端发送请求, 会先进行握手协议
// 后端会发送对应的加密key给客户端.
// 客户端再在发送数据的时候,用这个key解密请求.
func (w *WsServer) Handshake() {
	secretKey := ""
	key, err := w.GetProperty("secretKey")
	if err == nil {
		secretKey = key.(string)
	} else {
		secretKey = utils.RandSeq(16)
	}
	handshake := &Handshake{Key: secretKey}

	body := &RspBody{Name: HandshakeMsg, Msg: handshake}

	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		if data, err := utils.Zip(data); err == nil {
			err := w.wsConn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Println("write msg err:", err)
			}
		}
	}
}
