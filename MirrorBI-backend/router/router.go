package router

//全局路由注册
import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	v1 "mrbi/router/v1" // 导入 v1 路由
)

// SetupRoutes 全局路由设置
func SetupRoutes(r *gin.Engine) {
	//注册swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//启用 CORS 中间件，允许跨域资源共享
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://www.cloudhivegallery.cloud"}, // 允许的来源（前端地址）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                    // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},          // 允许的请求头
		ExposeHeaders:    []string{"Content-Length", "Authorization"},                            // 允许暴露的响应头
		AllowCredentials: true,                                                                   // 是否允许携带凭证（如 Cookies）
		AllowWildcard:    true,                                                                   // 是否允许任何来源
	}))
	// 注册 v1 路由
	v1.RegisterV1Routes(r)
}
