// RFC 1928 SOCKS Protocol Version 5
package socks
import "net"

const ProtoVer = 5
const NegotiationMethod = 0 // NO AUTHENTICATION REQUIRED
const MaxBuffLen = 512

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

+----+-----+-------+------+----------+----------+
|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
+----+-----+-------+------+----------+----------+
| 1  |  1  | X'00' |  1   | Variable |    2     |
+----+-----+-------+------+----------+----------+

+----+-----+-------+------+----------+----------+
|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
+----+-----+-------+------+----------+----------+
| 1  |  1  | X'00' |  1   | Variable |    2     |
+----+-----+-------+------+----------+----------+
func request(client net.Conn) error  {
	buffer := make([]byte, MaxBuffLen)

	// read VER, NMETHODS
	if _, err := client.Read(buffer[:2]); err != nil {
		return err
	}
}



