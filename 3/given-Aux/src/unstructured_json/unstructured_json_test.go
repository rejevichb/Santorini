package unstructured_json

import (
	"testing"
	"testing_utils"
)

func TestJsonEquality(t *testing.T) {
	testing_utils.AssertTrue(JsonEquality(1, 1), t)
	testing_utils.AssertTrue(JsonEquality("hi", "hi"), t)
	testing_utils.AssertTrue(JsonEquality(3.0, 3.0), t)
	testing_utils.AssertTrue(!JsonEquality(1.0, 2.0), t)
	testing_utils.AssertTrue(JsonEquality([]interface{}{1, 2, "a"}, []interface{}{1, 2, "a"}), t)
	testing_utils.AssertTrue(!JsonEquality([]interface{}{1, 2, "a"}, []interface{}{1, 2, []interface{}{1, 2}}), t)
	testing_utils.AssertTrue(JsonEquality(map[string]interface{}{"Hey": 1, "Hi": []interface{}{1, 2, 3}},
		map[string]interface{}{"Hey": 1, "Hi": []interface{}{1, 2, 3}}), t)
	testing_utils.AssertTrue(!JsonEquality(map[string]interface{}{"Hey": 1, "Hi": []interface{}{1, 2, 3}},
		map[string]interface{}{"Hey": 1, "Hi": []interface{}{1, 2, 4}}), t)
}
