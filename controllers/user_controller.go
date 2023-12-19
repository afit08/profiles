package controllers

import (
	"profile-api/configs"
	"profile-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func createUsers(c *gin.Context) {
	var user models.User
}
