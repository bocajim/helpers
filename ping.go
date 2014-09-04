package helpers

import (
	"bytes"
	"net"
	"os"
	"time"
)

const (
	ICMP_ECHO_REQUEST = 8
	ICMP_ECHO_REPLY   = 0
)

// returns a suitable 'ping request' packet, with id & seq and a
// payload length of pktlen
func makePingRequest(id, seq, pktlen int, filler []byte) []byte {
	p := make([]byte, pktlen)
	copy(p[8:], bytes.Repeat(filler, (pktlen-8)/len(filler)+1))

	p[0] = ICMP_ECHO_REQUEST // type
	p[1] = 0                 // code
	p[2] = 0                 // cksum
	p[3] = 0                 // cksum
	p[4] = uint8(id >> 8)    // id
	p[5] = uint8(id & 0xff)  // id
	p[6] = uint8(seq >> 8)   // sequence
	p[7] = uint8(seq & 0xff) // sequence

	// calculate icmp checksum
	cklen := len(p)
	s := uint32(0)
	for i := 0; i < (cklen - 1); i += 2 {
		s += uint32(p[i+1])<<8 | uint32(p[i])
	}
	if cklen&1 == 1 {
		s += uint32(p[cklen-1])
	}
	s = (s >> 16) + (s & 0xffff)
	s = s + (s >> 16)

	// place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero
	p[2] ^= uint8(^s & 0xff)
	p[3] ^= uint8(^s >> 8)

	return p
}

func parsePingReply(p []byte) (id, seq int) {
	id = int(p[4])<<8 | int(p[5])
	seq = int(p[6])<<8 | int(p[7])
	return
}

/*
func elapsedTime(start int64) float32 {
    t := float32(time.Nanoseconds() - start) / 1000000.0
    return t
}*/

func IcmpPing(addr string) (time.Duration, error) {

	raddr, err := net.ResolveIPAddr("ip4", addr) // *IPAddr
	if err != nil {
		return 0, err
	}

	ipconn, err := net.ListenIP("ip4:icmp", nil) // *IPConn (Conn 인터페이스 구현)
	if err != nil {
		return 0, err
	}

	sendid := os.Getpid() & 0xffff
	sendseq := 1
	pingpktlen := 64

	sendpkt := makePingRequest(sendid, sendseq, pingpktlen, []byte("Go Ping"))

	start := time.Now()

	n, err := ipconn.WriteToIP(sendpkt, raddr)
	if err != nil || n != pingpktlen {
		return 0, err
	}

	ipconn.SetDeadline(start.Add(time.Second * 5)) // 0.5 second

	resp := make([]byte, 1024)
	for {
		_, _, err := ipconn.ReadFrom(resp)
		//fmt.Printf("%d bytes from %s: icmp_req=%d time=%.2f ms\n", n, dst, sendseq, elapsedTime(start))

		// log.Printf("%x", resp)

		if err != nil {
			return 0, err
		}
		if resp[0] != ICMP_ECHO_REPLY {
			continue
		}
		rcvid, rcvseq := parsePingReply(resp)
		if rcvid != sendid || rcvseq != sendseq {
			return 0, err
		}
		break
	}

	return time.Now().Sub(start), nil
}
