package cmd

import (
	_ "fmt"
	_ "mrbi/config"
	"mrbi/internal/service"
	"mrbi/pkg/db"
	"mrbi/router"

	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	router.SetupRoutes(r)
	db.LoadDB()
	//创建一个生命周期和系统相同的后台协程
	go func() {
		service.ChartBackgroundService()
	}()
	r.Run(":8080")
}
