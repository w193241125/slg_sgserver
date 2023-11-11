package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// websocket 服务
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

func (w *WsServer) writeMsgLoop() {

	for {
		select {
		case msg := <-w.outChan:
			fmt.Println(msg)

		}
	}
}

func (w *WsServer) readMsgLoop() {
	//先读到客户端发送过来的数据, 处理后, 再回消息
	// 经过路由, 实际处理程序
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("接收消息出错:", err)
		}
		fmt.Println(data)
	}
	w.Close()
	//rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	//w.outChan <- rsp
}

func (w *WsServer) Close() {
	_ = w.wsConn.Close()
}
