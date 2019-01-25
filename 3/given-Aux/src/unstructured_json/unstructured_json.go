// Package unstructured_json provides a boolean function that can
// check for the equality of structure and content of two untyped
// inputs.
package unstructured_json

func JsonEquality(first interface{}, second interface{}) bool {
	switch fT := first.(type) {
	case int:
		other, ok := second.(int)
		return ok && fT == other
	case float64:
		other, ok := second.(float64)
		return ok && fT == other
	case bool:
		other, ok := second.(bool)
		return ok && fT == other
	case string:
		other, ok := second.(string)
		return ok && fT == other
	case []interface{}:
		other, ok := second.([]interface{})
		return ok && jsonArrEq(fT, other)
	case map[string]interface{}:
		other, ok := second.(map[string]interface{})
		return ok && jsonObjEq(fT, other)
	default:
		return false
	}
}

func jsonArrEq(first []interface{}, second []interface{}) bool {
	if len(first) != len(second) {
		return false
	}

	for i, v := range first {
		if !JsonEquality(v, second[i]) {
			return false
		}
	}

	return true
}

func jsonObjEq(first map[string]interface{}, second map[string]interface{}) bool {
	for k, v := range first {
		if oV, ok := second[k]; !ok || !JsonEquality(oV, v) {
			return false
		}
	}
	return true
}
