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

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

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
	user.Roles = "user"
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
		"role_name":  user.Roles,
		"username":   user.Username,
		"password":   user.Password,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created",
		"data":    result,
		"status":  http.StatusCreated,
	})
}

func ShowUser(c *gin.Context) {
	filter := bson.M{"_id": c.Value("id")}

	var users models.User
	err := userCollection.FindOne(context.Background(), filter).Decode(&users)

	if err != nil {
		// Check if the user is not found
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}

	result := []gin.H{
		{
			"id":       users.ID,
			"name":     users.Name,
			"image":    users.Image.Filename,
			"desc":     users.Desc,
			"job_name": users.JobName,
			"skills":   users.Skills,
			"roles":    users.Roles,
		},
	}

	if len(result) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No Data Roles",
			"data":    []gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get Users",
		"data":    result,
		"status":  http.StatusOK,
	})
}

func Login(c *gin.Context) {
	var user models.User

	request := new(models.LoginRequest)
	if err := c.ShouldBind(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	err := userCollection.FindOne(context.Background(), bson.M{"username": request.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "login or password is incorrect",
		})
		return
	}

	passwordIsValid := VerifyPassword(request.Password, user.Password)
	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid password",
		})
		return
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["roleType"] = user.Roles
	claims["exp"] = time.Now().Add(time.Hour * 3).Unix() // expired in 3 hours

	token, errGenerateToken := helpers.GenerateAllTokens(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		c.Status(http.StatusUnauthorized)
		return
	}
	c.SetCookie("jwt", token, 10800, "/", "0.0.0.0", false, true)
	result := gin.H{
		"name":  user.Name,
		"token": token,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successfully!!!",
		"status":  http.StatusOK,
		"data":    result,
	})
}

func Logout(c *gin.Context) {
	_, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Cookie not found",
		})
		return
	}
	c.SetCookie("jwt", "", -1, "/", "0.0.0.0", true, false)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func UpdateUsers(c *gin.Context) {
	var updateUser models.User
	userID := c.Param("id")

	uuid, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := c.ShouldBind(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateUser.ID = uuid.String()
	filter := bson.M{"_id": updateUser.ID}
	update := bson.M{"$set": updateUser}

	if updateUser.Image != nil {
		err := helpers.UploadImageToMinio(updateUser.Image, updateUser.Image.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading image to MinIO"})
			return
		}
	}

	_, err = userCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating"})
		return
	}
	fmt.Print(updateUser)
	result := gin.H{
		"id":         updateUser.ID,
		"name":       updateUser.Name,
		"image":      updateUser.Image.Filename,
		"desc":       updateUser.Desc,
		"job_name":   updateUser.JobName,
		"skill_name": updateUser.Skills,
		"role_name":  updateUser.Roles,
		"username":   updateUser.Username,
		"password":   updateUser.Password,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categori updated",
		"data":    result,
	})
}
