package utils

func MergeMap(first map[string]interface{}, second map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(first))
	for k, v := range first {
		out[k] = v
	}

	for k, v := range second {
		out[k] = v
	}
	return out
}
