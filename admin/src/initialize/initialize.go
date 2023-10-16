package initialize

import (
	"blitzSeckill/admin/src/cache"
	"blitzSeckill/admin/src/controller/activity"
	"blitzSeckill/admin/src/db"
	"fmt"
	"github.com/gin-gonic/gin"
)

// InitResource 初始化服务资源
func InitResource() {
	if err := db.InitDB(); err != nil {
		panic(fmt.Sprintf("InitDB err %s", err.Error()))
	}

	if err := cache.InitCache(); err != nil {
		panic(fmt.Sprintf("InitCache err %s", err.Error()))
	}
}

// RegisterRouter 注册路由
func RegisterRouter(router *gin.Engine) {
	{
		//创建活动
		router.POST("/activity/create", activity.CreateActivity)

		//查询活动列表
		router.GET("/activity/list", activity.GetActivityList)

		//限购
		router.GET("/sec/kill", activity.SecKill)

		router.GET("/sayhello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "sayhello",
			})
		})
		router.GET("/query", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "query",
			})
		})

		router.GET("/settlement/prePage", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "prePage",
			})
		})

		router.GET("/settlement/page", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "page",
			})
		})

		router.GET("settlement/initData", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "page",
			})
		})

		// 提交订单
		router.GET("settlement/submitData", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "page",
			})
		})
	}

}
