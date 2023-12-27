package routes

import (
	"profile-api/controllers"
	"profile-api/helpers"
	"profile-api/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	users := router.Group("/profile/api/users")
	{
		users.POST("/create", controllers.CreateUsers)
		users.PUT("/update/:id", controllers.UpdateUsers)
		users.GET("/show", middleware.EnsureUser(), controllers.ShowUser)
		users.POST("/signin", controllers.Login)
		users.POST("/signout", controllers.Logout)
		users.GET("/image/:filename", helpers.ShowImageFromMinio)
	}
}
