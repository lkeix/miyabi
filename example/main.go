package main

import (
	"fmt"
	"miyabi"
)

func main() {
	myb := miyabi.New()
	routing := miyabi.NewRouter()
	g1 := routing.NewGroup("/group1")
	g1.Apply(middleware1, middleware2)
	g1.GET("/test1", test)
	routing.AppendGroup(g1)
	routing.GET("/", test)
	routing.GET("/test1", test1)
	routing.GET("/test2", test2)
	routing.POST("/test3", test3)
	routing.GET("/foo/:user", test4)
	routing.GET("/repo/:user/:active", test5)
	routing.GET("/fizz/:user/:active/:bool/:okok", test6)
	routing.GET("/tmpl", test7)
	routing.Apply(middleware3)
	myb.Routing = routing
	myb.Serve(":8000")
}

func test(ctx *miyabi.Context) {
	fmt.Println("hello")
	ctx.Response.WriteResponse("hello")
}

func test1(ctx *miyabi.Context) {
	ctx.Response.WriteResponse("see you")
}

func test2(ctx *miyabi.Context) {
	fmt.Println("test2")
	ctx.Response.WriteResponse("test2")
}

func test3(ctx *miyabi.Context) {
	fmt.Println(ctx.Request.Data)
}

func test4(ctx *miyabi.Context) {
	fmt.Println(ctx.Request.PathParams)
	ctx.Response.WriteResponse(ctx.Request.PathParams)
}

func test5(ctx *miyabi.Context) {
	ctx.Response.WriteResponse(ctx.Request.PathParams)
}

func test6(ctx *miyabi.Context) {
	fmt.Println(ctx.Request.QueryParams["page"])
	ctx.Response.WriteResponse(ctx.Request.PathParams)
}

func test7(ctx *miyabi.Context) {
	ctx.AddTemplates("test", "./templates/test.tmpl", "./templates/test1.tmpl")
	ctx.Execute("test", map[string]string{
		"Title": "Hello!",
		"test":  "test",
	})
}

func middleware1(ctx *miyabi.Context) {
	fmt.Println("called Middleware1")
}

func middleware2(ctx *miyabi.Context) {
	fmt.Println("called Middleware2")
}

func middleware3(ctx *miyabi.Context)  {
	fmt.Println("called Middleware3")
}
