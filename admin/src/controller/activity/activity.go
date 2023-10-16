package activity

import (
	"blitzSeckill/admin/src/cache"
	"blitzSeckill/admin/src/model"
	"blitzSeckill/admin/src/service"
	"blitzSeckill/until/com"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func CreateActivity(ctx *gin.Context) {
	activity := &model.Activity{}
	//活动名称
	activity.ActivityName = ctx.PostForm("activity_name")
	//商品Id
	activity.ProductId, _ = strconv.Atoi(ctx.PostForm("product_id"))
	//活动开始时间
	activity.StartTime, _ = strconv.Atoi(ctx.PostForm("start_time"))
	//活动结束时间
	activity.EndTime, _ = strconv.Atoi(ctx.PostForm("end_time"))
	//商品数量
	activity.Total, _ = com.StrTo(ctx.PostForm("total")).Int()
	//商品速度
	activity.SecSpeed, _ = strconv.Atoi(ctx.PostForm("speed"))
	//购买限制
	activity.BuyLimit, _ = com.StrTo(ctx.PostForm("buy_limit")).Int()
	activity.BuyRate = ctx.PostForm("buy_rate")
	activityServer := service.NewActivityService()
	if err := activityServer.CreateActivity(activity); err != nil {
		log.Printf("ActivityServer.CreateActivity, Error : %v", err)
		ctx.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "failed",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
	})
	return
}

func GetActivityList(ctx *gin.Context) {
	ActivityService := service.NewActivityService()
	activityList, err := ActivityService.GetActivityList()
	if err != nil {
		log.Printf("ActivityService.GetActivityList, err : %v", err)
		ctx.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "failed",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": activityList,
	})
	return
}

func SecKill(ctx *gin.Context) {
	count := cache.Evalsha("product_2023_1016", 2)
	fmt.Print("限购结果：", count)
	ctx.JSON(200, gin.H{
		"message": count,
	})
}
