package dnsudp

import (
	"errors"
	"github.com/0x2E/rawdns"
	"net"
	"time"
)

// Send 向conn中发送DNS请求
func Send(conn net.Conn, subdomain string, id uint16) error {
	payload, err := rawdns.Marshal(id, 1, subdomain, rawdns.QTypeA)
	if err != nil {
		return errors.New("DNS marshal error: " + err.Error())
	}
	if _, err = conn.Write(payload); err != nil {
		return errors.New("DNS send error: " + err.Error())
	}
	return nil
}

// Receive 从conn中接收DNS请求，返回经过解析的结构体
func Receive(conn net.Conn, timeout int) (*rawdns.Message, error) {
	if err := conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout))); err != nil {
		return nil, err
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	// 超时错误内容：read udp 192.168.0.102:54012->8.8.8.8:53: i/o timeout
	if err != nil {
		return nil, err
	}

	// 解析
	resp, err := rawdns.Unmarshal(buf[:n])
	return resp, err
}
