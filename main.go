package main

import (
	"fmt"
	"log"
	"time"
	"wxcloudrun-golang/internal/app/service"
	"wxcloudrun-golang/internal/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载指定时区的位置信息
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	// 设置全局默认时区
	time.Local = loc
	// 获取当前时间
	now := time.Now()
	fmt.Println(now)
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}
	service := service.NewService()
	router := gin.Default()
	router.POST("/auth/login", service.WeChatLogin)
	router.GET("/courts", service.GetCounts)
	router.GET("/courts/:id", service.GetCountInfo)
	router.GET("/courts/:id/judge", service.JudgeLocation)

	router.GET("/events", service.GetEvents)
	router.GET("/videos", service.GetEventInfo)
	router.POST("/collects", service.ToggleCollectVideo)
	router.GET("/user/collects", service.GetCollectVideos)

	router.GET("/recommend/videos", service.GetRecommendVideos)
	router.POST("user/phone", service.GetUserPhone)

	// 8080 port
	log.Fatal(router.Run())
}
