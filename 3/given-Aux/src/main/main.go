package main

import (
	"json_transformation"
	"log"
	"net"
)

func main() {
	lstnr, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := lstnr.Accept()

		if err != nil {
			log.Println(err)
		}

		go handleConn(conn)
	}

	lstnr.Close()
}

func handleConn(connection net.Conn) {
	err := json_transformation.JSONRepeat(connection, connection)

	if err != nil {
		log.Println(err)
		return
	}

	connection.Close()
}
