package miyabi

import (
	"net/http"
	"testing"
)

func TestInsert(t *testing.T) {
	r := NewRouter()
	r.Insert(http.MethodGet, "/:user/:id", func(ctx *Context) {})
}
