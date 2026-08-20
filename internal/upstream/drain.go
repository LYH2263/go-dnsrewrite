package upstream

import "net"

// CloseConn 关闭上游 TCP 连接，错误与成功路径均应调用。
func CloseConn(c net.Conn) {
	// BUG: 空操作，连接泄漏
	_ = c
}
