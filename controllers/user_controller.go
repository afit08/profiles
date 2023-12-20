package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"profile-api/configs"
	"profile-api/helpers"
	"profile-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// var validate = validator.New()

// type loginResponse struct {
// 	Username string `json:"username"`
// 	Token    string `json:"token"`
// }

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	return err == nil
}

func CreateUsers(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		// Use StatusJSON to set both HTTP status and JSON response
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid form data: %s", err.Error())})
		return
	}

	if user.Image == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	password := HashPassword(user.Password)

	user.ID = uuid.New().String()
	user.Password = password
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Upload gambar ke MinIO
	err := helpers.UploadImageToMinio(user.Image, user.Image.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading image to MinIO"})
		return
	}

	// Insert the categori into the database
	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		// Use StatusJSON for consistent response format
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Create Data"})
		return
	}

	result := gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"image":      user.Image.Filename,
		"desc":       user.Desc,
		"job_name":   user.JobName,
		"skill_name": user.Skills,
		"username":   user.Username,
		"password":   user.Password,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created",
		"data":    result,
	})
}
