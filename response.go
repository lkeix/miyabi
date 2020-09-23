package miyabi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response is http response class.
type Response struct {
	writer    *http.ResponseWriter
	Status    int
	Size      int64
	Committed bool
}

// NewResponse create response instance.
func NewResponse(w *http.ResponseWriter) *Response {
	return &Response{writer: w}
}

func (resp *Response) reset() {
	resp.writer = nil
	resp.Size = -1
	resp.Status = 200
}

// WriteResponse write interface
// response arge type is byte[], string, interface.
func (resp *Response) WriteResponse(response interface{}) error {
	w := *resp.writer
	responseStr := ""
	var err error
	switch response.(type) {
	case string:
		responseStr = response.(string)
	case ([]byte):
		responseStr = string(response.([]byte))
	default:
		var tmp []byte
		tmp, err = json.Marshal(response)
		responseStr = string(tmp)
	}
	if err == nil {
		fmt.Fprintf(w, responseStr)
	}
	return err
}
