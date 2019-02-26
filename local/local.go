package local

import (
	"log"
	"net"
)

func newClient(client net.Conn) {

}

func Start(addr string)  {
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
		return
	}

	log.Printf("client listening on %s\n", addr)
	for {
		client, err := ls.Accept()
		if nil != err {
			log.Printf("failed to accept: %v\n", err)
			continue
		}

		go newClient(client)
	}
}
