// RFC 1928 SOCKS Protocol Version 5
package socks

import (
	"encoding/binary"
	"log"
	"net"
	"strconv"
)

const ProtoVer = 5
const NegotiationMethod = 0 // NO AUTHENTICATION REQUIRED
const MaxBuffLen = 512

const (
	AType_IPV4 = 1
	AType_DOMAIN = 3
	AType_IPV6 = 4
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
func negotiation(client net.Conn) error {
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

	// write VER METHOD
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
func request(client net.Conn) error  {
	buffer := make([]byte, MaxBuffLen)

	// read VER, REP, RSV, ATYP
	if _, err := client.Read(buffer[:4]); err != nil {
		return err
	}

	var host, port string
	atyp := buffer[3]
	switch atyp {
	case AType_IPV4:
		if _, err := client.Read(buffer[:net.IPv4len + 2]); err != nil {
			return err
		}

		host = net.IP(buffer[:net.IPv4len]).String()
		log.Printf("ipv4 host: %v, %v\n", host, buffer[:net.IPv4len])
		port = strconv.Itoa(int(binary.LittleEndian.Uint16(buffer[net.IPv4len: net.IPv4len + 1])))
	case AType_DOMAIN:
	case AType_IPV6:
		if _, err := client.Read(buffer[:net.IPv6len + 2]); err != nil {
			return err
		}

		host = net.IP(buffer[:net.IPv6len]).String()
		log.Printf("ipv6 host: %v, %v\n", host, buffer[:net.IPv6len])
		port = strconv.Itoa(int(binary.LittleEndian.Uint16(buffer[net.IPv6len: net.IPv6len + 1])))
	}

	return nil
}



