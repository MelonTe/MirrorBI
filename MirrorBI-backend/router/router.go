package router

//全局路由注册
import (
	"mrbi/internal/consts"
	"mrbi/internal/controller"
	"mrbi/internal/middleware"
	"mrbi/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var userService = service.NewUserService()

// SetupRoutes 全局路由设置
func SetupRoutes(r *gin.Engine) {
	//注册swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//启用session中间件
	middleware.InitSession(r)
	//启用 CORS 中间件，允许跨域资源共享
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://www.cloudhivegallery.cloud"}, // 允许的来源（前端地址）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                    // 允许的 HTTP 方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},          // 允许的请求头
		ExposeHeaders:    []string{"Content-Length", "Authorization"},                            // 允许暴露的响应头
		AllowCredentials: true,                                                                   // 是否允许携带凭证（如 Cookies）
		AllowWildcard:    true,                                                                   // 是否允许任何来源
	}))

	//注册路由
	api := r.Group("/api")
	userAPI := api.Group("/user")
	{
		userAPI.POST("/register", controller.UserRegister)
		userAPI.POST("/login", controller.UserLogin)
		userAPI.GET("/get/login", controller.GetLoginUser)
		userAPI.POST("/logout", controller.UserLogout)
		userAPI.GET("/get/vo", controller.GetUserVOById)
		//以下需要权限
		userAPI.POST("/list/page/vo", middleware.AuthCheck(consts.ADMIN_ROLE), controller.ListUserVOByPage)
		userAPI.POST("/update", middleware.AuthCheck(consts.ADMIN_ROLE), controller.UpdateUser)
		userAPI.POST("/delete", middleware.AuthCheck(consts.ADMIN_ROLE), controller.DeleteUser)
		userAPI.POST("/add", middleware.AuthCheck(consts.ADMIN_ROLE), controller.AddUser)
		userAPI.GET("/get", middleware.AuthCheck(consts.ADMIN_ROLE), controller.GetUserById)
		//userAPI.POST("/avatar", middleware.LoginCheck(), controller.UploadAvatar)
		userAPI.POST("/edit", middleware.LoginCheck(), controller.EditUser)
	}
}
