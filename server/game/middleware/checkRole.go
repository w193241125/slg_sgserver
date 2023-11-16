package middleware

import (
	"log"
	"sgserver/constant"
	"sgserver/net"
)

func CheckRole() net.MiddlewareFunc {
	return func(next net.HandlerFunc) net.HandlerFunc {
		return func(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
			log.Println("角色检测...")
			_, err := req.Conn.GetProperty("role")
			if err != nil {
				rsp.Body.Code = constant.SessionInvalid
			}
			next(req, rsp)
		}
	}
}
