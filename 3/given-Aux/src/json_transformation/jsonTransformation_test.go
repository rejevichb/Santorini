package json_transformation

import (
	"testing"
	"testing_utils"
	"unstructured_json"
)

func TestJSONTransform(t *testing.T) {
	testSlice := []interface{}{"a", "b", "c", "d"}
	expectedOut := []interface{}{[]interface{}{3, "a"},
		[]interface{}{2, "b"},
		[]interface{}{1, "c"},
		[]interface{}{0, "d"}}

	jsonTransform(testSlice)

	testing_utils.AssertTrue(unstructured_json.JsonEquality(testSlice, expectedOut), t)

	simpleJson := []interface{}{0, 2}

	jsonTransform(simpleJson)

	testing_utils.AssertTrue(unstructured_json.JsonEquality(simpleJson, []interface{}{[]interface{}{1, 0},
		[]interface{}{0, 2}}), t)

}
