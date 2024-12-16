package log

import "encoding/json"

func ToJsonAndTruncateForLogging(req any) string {
	str, err := json.Marshal(req)
	if len(str) > 1000 {
		str = str[0:1000]
	}
	if err != nil {
		return "Error when serializing request: " + err.Error()
	}
	return string(str)
}
