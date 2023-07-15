package main

import (
	"fmt"
	"log"
	"wxcloudrun-golang/internal/app/service"
	"wxcloudrun-golang/internal/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}
	service := service.NewService()
	router := gin.Default()
	router.POST("/auth/login", service.WeChatLogin)
	router.POST("/user/court", service.StoreCourt)
	router.GET("/user/download", service.GetUserDownload)
	router.GET("/courts", service.GetCounts)
	router.GET("/courts/:id", service.GetCountInfo)
	router.GET("/courts/:id/judge", service.JudgeLocation)

	router.GET("/events", service.GetEvents)
	router.GET("/videos", service.GetVideos)
	router.POST("/videos", service.StoreVideo)
	router.GET("/records", service.GetRecords)
	router.POST("/collects", service.ToggleCollectVideo)
	router.POST("/user/event", service.CollectUserEvent)
	router.GET("/user/collects", service.GetCollectVideos)

	router.GET("/recommend/videos", service.GetRecommendVideos)
	router.POST("/user/phone", service.GetUserPhone)
	router.POST("/survey", service.CollectSurvey)

	// 8080 port
	log.Fatal(router.Run())
}
