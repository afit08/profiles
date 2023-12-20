package routes

import (
	"profile-api/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	users := router.Group("/profile/api/users")
	{
		users.POST("/create", controllers.CreateUsers)
	}
}
