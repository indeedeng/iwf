package interpreter

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func caller(skip int) string {
	_, path, lineno, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v:%v", filepath.Base(path), lineno)
}

func LastCaller() string {
	return caller(3)
}
