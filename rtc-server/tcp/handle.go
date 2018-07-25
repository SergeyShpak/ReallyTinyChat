package tcp

import (
	"log"
	"net"
)

func Handle(conn net.Conn) error {
	log.Println("Handling connection: ", conn)
	return nil
}
