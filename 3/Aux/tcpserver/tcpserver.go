package tcpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

const (
	CONN_PORT = "8000"
	CONN_TYPE = "tcp"
)

func main() {
	// Start TCP Server
	if listener, err := net.Listen(CONN_TYPE, ":"+CONN_PORT); err != nil {
		panic(err)
	} else {
		defer listener.Close()

		if conn, err := listener.Accept(); err != nil {
			panic(err)
		} else {
			defer conn.Close()
			handleIO(conn)
		}
	}
}

func handleIO(stream io.ReadWriter) {
	var validObject interface{}
	objects := make([]interface{}, 0)
	reader := json.NewDecoder(stream)

	for {
		if err := reader.Decode(&validObject); err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		} else {
			objects = append(objects, validObject)
		}
	}

	postJSON(stream, objects)
}

func postJSON(writer io.Writer, objects []interface{}) {
	for idx, val := range objects {
		obj, _ := json.Marshal(val)
		outString := fmt.Sprintf("[%v,%s]\n", len(objects)-idx-1, string(obj))

		writer.Write([]byte(outString))
	}
}
