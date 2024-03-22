package check_alive

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"time"
	"zfofa/backend/core/fetch"
)

// checkAliveWithIcmpEcho 使用 ICMP ECHO 方式探测主机是否存活
func checkAliveWithIcmpEcho(ip string) bool {
	rIpv4 := &net.IPAddr{IP: net.ParseIP(ip)}
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return false
	}
	defer conn.Close()

	// 构造ICMP Echo消息
	echoMsg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xfff,
			Seq:  1,
			Data: []byte("ICMP test by zfofa."),
		},
	}

	msgB, _ := echoMsg.Marshal(nil)
	_, err = conn.WriteTo(msgB, rIpv4)
	if err != nil {
		return false
	}

	// 设置超时5秒
	reply := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, _, err = conn.ReadFrom(reply)
	if err != nil {
		return false
	}

	return true
}

func checkAliveWithHttp(url string) bool {
	r := fetch.CreateRequest(true, 10*time.Second, nil)
	_, err := r.Get(url, nil, nil)
	if err != nil {
		return false
	}

	return true
}

// checkAliveWithTcp 使用TCP探测主机存活
func checkAliveWithTcp(ip string, port int) bool {
	rHost := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", rHost, 5*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

// checkAliveWithUdp 使用UDP连接的方式探测存活
func checkAliveWithUdp(ip string, port int) bool {
	rAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	conn, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		return false
	}
	defer conn.Close()

	// 设置超时5秒
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, err = conn.Write([]byte("UDP test by zfofa."))
	if err != nil {
		return false
	}

	// 只要能够收到回复，说明存活
	reply := make([]byte, 1024)
	_, _, err = conn.ReadFrom(reply)
	if err != nil {
		return false
	}

	return true
}
