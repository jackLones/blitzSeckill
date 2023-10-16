package main

import (
	"blitzSeckill/admin/src/config"
	"blitzSeckill/admin/src/initialize"
	"fmt"
	"github.com/niuniumart/gosdk/gin"
)

func main() {

	// 初始化配置
	config.Init()
	// 初始资源，主要是MySQL，Redis连接
	initialize.InitResource()
	// 创建一个web服务
	router := gin.CreateGin()
	// 这里跳进去就能看到有哪些接口
	initialize.RegisterRouter(router)
	fmt.Println("before router run")
	// 启动web server，这一步之后这个主协程启动会阻塞在这里，请求可以通过gin的子协程进来
	err := gin.RunByPort(router, config.Conf.Common.Port)
	fmt.Println(err)
}
