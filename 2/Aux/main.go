package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	decoder := json.NewDecoder(stdin)

	var objects = make([]interface{}, 0)
	var obj interface{}
	for {
		err := decoder.Decode(&obj)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		objects = append(objects, obj)
	}

	for idx, val := range objects {
		obj, _ := json.Marshal(val)
		fmt.Printf("[%v,%s]\n", len(objects)-idx-1, string(obj))
	}

	return
}
