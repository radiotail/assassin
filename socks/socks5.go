// RFC 1928 SOCKS Protocol Version 5
package socks

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
)

const ProtoVer = 5
const NegotiationMethod = 0 // NO AUTHENTICATION REQUIRED
const MaxBuffLen = 512

// SOCKS request commands
const (
	CmdConnect      = 1
	CmdBind         = 2
	CmdUdpassociate = 3
)

// SOCKS address types
const (
	AtypIpv4   = 1
	AtypDomain = 3
	AtypIpv6   = 4
)

//	recv:
// +----+----------+----------+
// |VER | NMETHODS | METHODS  |
// +----+----------+----------+
// | 1  |    1     | 1 to 255 |
// +----+----------+----------+

// send:
// +----+--------+
// |VER | METHOD |
// +----+--------+
// | 1  |   1    |
// +----+--------+
func Negotiation(client net.Conn) error {
	buffer := make([]byte, MaxBuffLen)
	// read VER, NMETHODS
	if _, err := client.Read(buffer[:2]); err != nil {
		return err
	}

	// read METHODS
	nMethods := buffer[1]
	if _, err := client.Read(buffer[:nMethods]); err != nil {
		return err
	}

	// replay: write VER METHOD
	if _, err := client.Write([]byte{ProtoVer, NegotiationMethod}); err != nil {
		return err
	}

	return nil
}

// recv:
// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+

// send:
// +----+-----+-------+------+----------+----------+
// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
func Request(client net.Conn) (string, error)  {
	buffer := make([]byte, MaxBuffLen)

	// read VER, CMD, RSV, ATYP
	if _, err := client.Read(buffer[:4]); err != nil {
		return "", err
	}

	cmd := buffer[1]
	var host, port string
	atyp := buffer[3]
	switch atyp {
	case AtypIpv4:
		if _, err := client.Read(buffer[:net.IPv4len + 2]); err != nil {
			return "", err
		}

		host = net.IP(buffer[:net.IPv4len]).String()
		log.Printf("ipv4 host: %v, %v\n", host, buffer[:net.IPv4len])
		port = strconv.Itoa(int(binary.LittleEndian.Uint16(buffer[net.IPv4len: net.IPv4len + 2])))
	case AtypDomain:
		if _, err := client.Read(buffer[:1]); err != nil {
			return "", err
		}

		domainLen := buffer[0]
		if _, err := client.Read(buffer[:domainLen + 2]); err != nil {
			return "", err
		}

		host = string(buffer[:domainLen])
		port = strconv.Itoa(int(binary.LittleEndian.Uint16(buffer[domainLen: domainLen + 2])))
	case AtypIpv6:
		if _, err := client.Read(buffer[:net.IPv6len + 2]); err != nil {
			return "", err
		}

		host = net.IP(buffer[:net.IPv6len]).String()
		log.Printf("ipv6 host: %v, %v\n", host, buffer[:net.IPv6len])
		port = strconv.Itoa(int(binary.LittleEndian.Uint16(buffer[net.IPv6len: net.IPv6len + 2])))
	default:
		err := fmt.Errorf("unknown atyp: %d", atyp)
		return "", err
	}
	log.Printf("host: %s, port: %d\n", host, port)

	// respone: write VER, REP, RSV, ATYP, BND.ADDR, BND.PORT
	switch cmd {
	case CmdConnect:
		if _, err := client.Write([]byte{ProtoVer, 0, 0, AtypIpv4, 0, 0, 0, 0, 0, 0}); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("request commands is support: %d", cmd)
	}

	return net.JoinHostPort(host, port), nil
}



