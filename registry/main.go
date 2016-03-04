package main

import (
	"registry/middleware"
	"registry/server"
)

type serverOpts struct {
	host string
}

const (
	host = ":9090"
)

func main() {
	opts := serverOpts{host: host}
	//创建Url路由
	router := server.NewRouter()
	//添加web中间件
	n := middleware.New()
	//将router添加到中间件栈最底端
	n.UseHandler(router)
	n.Run(opts.host)
}
