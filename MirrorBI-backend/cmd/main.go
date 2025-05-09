package cmd

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "mrbi/config"
	"mrbi/pkg/db"
	"mrbi/router"
)

func Main() {
	r := gin.Default()
	router.SetupRoutes(r)
	db.LoadDB()
	r.Run(":8080")
}
