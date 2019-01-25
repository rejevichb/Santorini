// Package json_transformation provides the ability to read
// a stream of data, package it as JSON data, and output
// in a specified reverse-chronological format.
package json_transformation

import (
	"encoding/json"
	"errors"
	"io"
)

func JSONRepeat(rdr io.Reader, wrtr io.Writer) error {
	dec := json.NewDecoder(rdr)
	enc := json.NewEncoder(wrtr)

	var jsons []interface{}

	for {
		var currJSON interface{}

		err := dec.Decode(&currJSON)

		if err != nil {
			if err != io.EOF {
				return errors.New("Error while attempting to read stream.")
			}

			break
		}

		jsons = append(jsons, currJSON)
	}

	jsonTransform(jsons)

	for _, obj := range jsons {
		enc.Encode(obj)
	}

	return nil
}

// jsonTransform transforms a slice of JSON expressions. Each JSON expression becomes the second
// member of a JSON array where the first member is the number of expressions in the slice
// - the index of the expression in the slice - 1.
// Therefore the first expression becomes associated with the last possible index and the last with the first.
func jsonTransform(objs []interface{}) {
	for i := 0; i < len(objs); i++ {
		inverseInd := len(objs) - i - 1
		objs[i] = []interface{}{inverseInd, objs[i]}
	}
}
