package dnsudp

import (
	"github.com/0x2E/rawdns"
	"github.com/pkg/errors"
	"net"
	"time"
)

// Send 向conn中发送DNS请求
func Send(conn net.Conn, domain string, id uint16, qtype rawdns.QType) error {
	payload, err := rawdns.Marshal(id, 1, domain, qtype)
	if err != nil {
		return errors.Wrap(err, "DNS marshal error")
	}
	if _, err = conn.Write(payload); err != nil {
		return errors.Wrap(err, "UDP send error")
	}
	return nil
}

// Receive 从conn中接收DNS报文，返回解析后的结构体
func Receive(conn net.Conn, timeout int) (*rawdns.Message, error) {
	err := conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
	if err != nil {
		return nil, errors.Wrap(err, "failed to set read deadline")
	}
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	// 超时错误内容：read udp 192.168.0.102:54012->8.8.8.8:53: i/o timeout
	if err != nil {
		return nil, errors.Wrap(err, "UDP read error")
	}

	// 解析
	resp, err := rawdns.Unmarshal(buf[:n])
	if err != nil {
		err = errors.Wrap(err, "DNS unmarshal error")
	}
	return resp, err
}
