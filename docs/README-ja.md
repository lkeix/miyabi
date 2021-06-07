# 使い方

* GET のハンドラーを定義

```Go
package main

import "miyabi"

func main() {
  myb := miyabi.New()
  router := miyabi.NewRouter()
  router.GET("/", func (ctx *miyabi.Context){
    ctx.Response.WriteResponse("Hello World")
  })
  myb.Router = router
	myb.Serve(":8000")
}
```

* POST のハンドラーを定義

```Go
package main

import "miyabi"

func main() {
  myb := miyabi.New()
  router := miyabi.NewRouter()
  router.POST("/", func (ctx *miyabi.Context){
    ctx.Response.WriteResponse("Hello World")
  })
  myb.Router = router
	myb.Serve(":8000")
}
```

* パスパラメータを定義

```Go
package main

import "miyabi"

func main() {
  myb := miyabi.New()
  router := miyabi.NewRouter()
  router.GET("/repo/:user", func (ctx *miyabi.Context){  
    req := ctx.Request
    req.Parse()
    ctx.Response.WriteResponse(req.PathParams)
  })
  myb.Router = router
	myb.Serve(":8000")
}
```


## 各ファイル、ディレクトリの内容

* miyabi.go

  ホスティングのエントリーポイントなどに関するモジュール

* context.go

  リクエストやレスポンスに関する部分を集約、ハンドラの引数として使用する。

* request.go, response.go

  リクエストやレスポンスに関する部分を実装、PHPライクにアクセスできるようにプリプロセッサを追加している。
  
* sessions
  
  セッションに関するモジュール群

* router.go
  Webルーティングに関する部分のモジュール、ルーティングはトライ木で実装

* logger.go
  ログに関するモジュール

* example

  miyabiを使ってホスティングする際の例