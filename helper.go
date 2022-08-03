package miyabi

import (
	"errors"
	"fmt"
)

func pathValidate(path string) string {
	if path == "" {
		return "/"
	}
	if path[0] != separator {
		return "/" + path
	}
	return path
}

func handlerValidate(handler HandlerFunc, path string) error {
	if handler == nil {
		return errors.New(fmt.Sprintf("%s: handler is nil", path))
	}
	return nil
}

func maxi(a, b int) int {
	if a < b {
		return a
	}
	return b
}
