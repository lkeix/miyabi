package main

import (
	"fmt"
	"miyabi"
	"miyabi/sessions"
	"strconv"
)

func main() {
	myb := miyabi.New()
	router := miyabi.NewRouter()
	sessions.NewSessions()
	g1 := router.NewGroup("/group1")
	g1.Apply(middleware1, middleware2)
	g1.GET("/test1", test)
	router.AppendGroup(g1)
	router.GET("/", test)
	router.GET("/test1", test1)
	router.GET("/test2", test2)
	router.POST("/test3", test3)
	router.GET("/foo/:user", test4)
	router.GET("/repo/:user/:active", test5)
	router.GET("/fizz/:user/:active/:bool/:okok", test6)
	router.GET("/tmpl", test7)
	router.Apply(middleware3)
	myb.Router = router
	myb.Serve(":8000")
}

func test(ctx *miyabi.Context) {
	fmt.Println("hello")
	ctx.Response.WriteResponse("hello")
	sessions.Start(ctx)
}

func sessionTst(sess *sessions.Sessions) miyabi.HandlerFunc {
	return func(ctx *miyabi.Context) {
	}
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
	session := sessions.Start(ctx)
	var cnt int64
	if session.Get("cnt") != nil {
		cnt = session.Get("cnt").(int64)
	}
	cnt++
	session.Set("cnt", cnt)
	ctx.AddTemplates("test", "./templates/test.tmpl", "./templates/test1.tmpl")
	cntStr := strconv.Itoa(int(cnt))
	ctx.Execute("test", map[string]string{
		"Title": "Hello!",
		"test":  "test",
		"cnt":   cntStr,
	})
}

func middleware1(ctx *miyabi.Context) {
	fmt.Println("called Middleware1")
}

func middleware2(ctx *miyabi.Context) {
	fmt.Println("called Middleware2")
}

func middleware3(ctx *miyabi.Context) {
	fmt.Println("called Middleware3")
}
