package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"sgserver/constant"
	"sgserver/utils"
	"sync"
	"time"
)

type syncCtx struct {
	//goroutine 的上下文, 包含goroutine的状态,环境,现场等信息
	ctx     context.Context
	cancel  context.CancelFunc
	outChan chan *RspBody
}

func NewSyncCtx() *syncCtx {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	return &syncCtx{
		ctx:     ctx,
		cancel:  cancel,
		outChan: make(chan *RspBody), // 无缓存channel
	}
}

func (s *syncCtx) wait() *RspBody {
	select {
	case msg := <-s.outChan:
		return msg
	case <-s.ctx.Done():
		log.Println("代理服务响应超时")
		return nil
	}
}

type ClientConn struct {
	wsConn        *websocket.Conn
	handshake     bool
	handshakeChan chan bool
	isClosed      bool
	property      map[string]interface{}
	propertyLock  sync.RWMutex
	Seq           int64
	onPush        func(conn *ClientConn, body *RspBody)
	onClose       func(conn *ClientConn)
	syncCtxMap    map[int64]*syncCtx
	syncCtxLock   sync.RWMutex
}

func (c *ClientConn) waitHandShake() bool {
	//等待握手消息
	// 还需处理程序异常超时问题.(一直无法响应)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case _ = <-c.handshakeChan:
		log.Println("握手成功")
		return true
	case <-ctx.Done():
		log.Println("握手超时")
		return false
	}

}

func (c *ClientConn) Start() bool {
	//一直不停收消息
	//等待握手的消息返回
	c.handshake = false
	go c.wsReadLoop()
	return c.waitHandShake()
}

func (c *ClientConn) wsReadLoop() {
	//for {
	//	_, data, err := c.wsConn.ReadMessage()
	//	if err != nil {
	//		log.Println("读取消息出错了....", err)
	//	}
	//	fmt.Println(data)
	//	//收到握手消息
	//	c.handshake = true
	//	c.handshakeChan <- true
	//
	//}

	defer func() {
		if err := recover(); err != nil {
			log.Println("捕捉到异常", err)
			c.Close()
		}
	}()
	for {
		_, data, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Println("clientconn收消息出错:", err)
			break
		}
		//fmt.Println(data)
		//前端发送过来的消息就是JSON格式
		//1. data 解压 unzip
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错, 非法格式: ", err)
			continue
		}
		//2. 前端的消息加密的,需要解密
		secretKey, err := c.GetProperty("secretKey")
		if err == nil {
			//有加密
			key := secretKey.(string)
			//客户端传过来的数据是加密的 需要解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("Client数据格式有误，解密失败:", err)
			} else {
				data = d
			}
		}
		//3. data 转为body
		body := &RspBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误, 非法格式:", err)
		} else {
			//握手 别的一些请求.
			if body.Seq == 0 {
				if body.Name == HandshakeMsg {
					hs := &Handshake{}
					mapstructure.Decode(body.Msg, hs)
					if hs.Key != "" {
						c.SetProperty("secretKey", hs.Key)
					} else {
						c.RemoveProperty("secretKey")
					}

					c.handshake = true
					c.handshakeChan <- true
				} else {
					//不是握手就通知其他业务
					if c.onPush != nil {
						c.onPush(c, body)
					}
				}

			} else {
				c.syncCtxLock.RLock()
				ctx, ok := c.syncCtxMap[body.Seq]
				c.syncCtxLock.RUnlock()
				if ok {
					ctx.outChan <- body
				} else {
					log.Println("no seq syncCtx find")
				}

			}
		}
	}
	c.Close()
}
func (c *ClientConn) Close() {
	_ = c.wsConn.Close()
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn {
	return &ClientConn{
		wsConn:        wsConn,
		handshakeChan: make(chan bool),
		Seq:           0,
		isClosed:      false,
		property:      make(map[string]interface{}),
		syncCtxMap:    map[int64]*syncCtx{},
	}
}

func (c *ClientConn) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *ClientConn) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("property no found")
	}
}

func (c *ClientConn) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}

func (c *ClientConn) Addr() string {
	return c.wsConn.RemoteAddr().String()

}

func (c *ClientConn) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	//c.outChan <- rsp
	fmt.Println(rsp)
	c.Write(rsp.Body)
}
func (c *ClientConn) SetOnPush(hook func(conn *ClientConn, body *RspBody)) {
	c.onPush = hook
}

func (c *ClientConn) Write(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	secretKey, err := c.GetProperty("secretKey")
	if err == nil {
		//有加密
		key := secretKey.(string)
		//数据做加密
		data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
		if err != nil {
			log.Println("加密失败", err)
			return err
		}
	}
	//数据压缩
	if data, err := utils.Zip(data); err == nil {
		err := c.wsConn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Println("写数据失败", err)
			return err
		}
	} else {
		log.Println("压缩数据失败", err)
		return err
	}
	return nil
}

func (c *ClientConn) Send(name string, msg interface{}) *RspBody {
	//把请求发送给代理服务器 登录服务器 等待返回.
	c.Seq += 1
	seq := c.Seq
	sc := NewSyncCtx()
	c.syncCtxLock.Lock()
	c.syncCtxMap[seq] = sc
	c.syncCtxLock.Unlock()
	rsp := &RspBody{Name: name, Seq: seq, Code: constant.OK}

	//req请求
	req := &ReqBody{Seq: seq, Name: name, Msg: msg}
	err := c.Write(req)

	if err != nil {
		sc.cancel()
	} else {
		r := sc.wait()
		if r == nil {
			rsp.Code = constant.ProxyConnectError
		} else {
			rsp = r
		}
	}
	c.syncCtxLock.Lock()
	delete(c.syncCtxMap, seq)
	c.syncCtxLock.Unlock()

	return rsp
}
