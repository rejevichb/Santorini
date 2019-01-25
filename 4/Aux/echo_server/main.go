package main

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/google/uuid"
)

func main() {
	serv, _ := net.Listen("tcp", ":8000")

	conn, _ := serv.Accept()

	buf := make([]byte, 1024)

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	n, _ := r.Read(buf)
	fmt.Println(string(buf[:n]))

	id, _ := uuid.NewRandom()
	w.Write([]byte(id.String()))
	w.Flush()

	for {
		read, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				return
			}
			fmt.Println(err)
		} else {
			fmt.Println(string(buf[:read]))
			w.Write([]byte("hello world"))
			w.Flush()
		}
	}
}
