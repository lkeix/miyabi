package main

import (
	"fmt"
	"miyabi"
)

func main() {
	myb := miyabi.New()
	routing := miyabi.NewRouter()
	routing.GET("/", test)
	routing.GET("/test1", test1)
	routing.GET("/test2", test2)
	myb.Routing = routing
	myb.Serve(":8000")
}

func test(ctx *miyabi.Context) {
	fmt.Println("hello")
	ctx.Response.WriteResponse("hello")
}

func test1(ctx *miyabi.Context) {
	fmt.Println("see you")
	ctx.Response.WriteResponse("see you")
}
func test2(ctx *miyabi.Context) {
	fmt.Println("test2")
	ctx.Response.WriteResponse("test2")
}
