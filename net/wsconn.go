package net

import "sync"

// ReqBody 请求体格式
type ReqBody struct {
	Seq   int64       `json:"seq"`
	Name  string      `json:"name"`
	Msg   interface{} `json:"msg"`
	Proxy string      `json:"proxy"`
}

// RspBody 响应体格式
type RspBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

// WsMsgReq websocket 请求体格式
type WsMsgReq struct {
	Body    *ReqBody
	Conn    WSConn
	Context *WsContext
}
type WsContext struct {
	mutex    sync.RWMutex
	property map[string]interface{}
}

func (w *WsContext) Set(key string, value interface{}) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.property[key] = value
}

func (w *WsContext) Get(key string) interface{} {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	value, ok := w.property[key]
	if ok {
		return value
	}
	return nil
}

type WsMsgRsp struct {
	Body *RspBody
}

// WSConn 理解为 request 请求, 会有参数,  这里就是放参数 取参数的接口.
type WSConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

type Handshake struct {
	Key string `json:"key"`
}

type Heartbeat struct {
	CTime int64 `json:"ctime"`
	STime int64 `json:"stime"`
}
