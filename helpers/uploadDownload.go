package helpers

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var ErrInvalidFileExtension = errors.New("invalid file extension")

func UploadImageToMinio(imageData *multipart.FileHeader, objectName string) error {
	// Check if the file has a valid extension (PNG or JPG)
	ext := filepath.Ext(objectName)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return ErrInvalidFileExtension
	}

	// Implementasi koneksi ke MinIO
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretKey := os.Getenv("MINIO_SECRET_ACCESS_KEY_ID")
	bucketName := "gin-api"

	fileData, err := imageData.Open()
	if err != nil {
		return err
	}
	defer fileData.Close()

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Ganti menjadi true jika menggunakan koneksi aman (HTTPS)
	})
	if err != nil {
		return err
	}

	// Upload file content ke MinIO
	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, fileData, imageData.Size, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func DownloadImage(c *gin.Context) {
	// Implementasi koneksi ke MinIO
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretKey := os.Getenv("MINIO_SECRET_ACCESS_KEY_ID")

	filename := c.Param("filename")
	bucketName := "gin-api"
	objectName := filename

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Ganti menjadi true jika menggunakan koneksi aman (HTTPS)
	})
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Set expiration duration for the presigned URL
	expiration := 24 * time.Hour // You can adjust the expiration duration as needed

	// Create a presigned URL for the image
	url, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expiration, url.Values{})
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Redirect the user to the presigned URL
	c.Redirect(http.StatusTemporaryRedirect, url.String())
}

func ShowImageFromMinio(c *gin.Context) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY_ID")
	secretKey := os.Getenv("MINIO_SECRET_ACCESS_KEY_ID")

	filename := c.Param("filename")
	bucketName := "gin-api"
	objectName := filename

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // Ganti menjadi true jika menggunakan koneksi aman (HTTPS)
	})
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Fetch the object from Minio
	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer object.Close()

	// Set the appropriate headers for an image response
	c.Header("Content-Type", "image/jpeg") // Adjust the content type based on your image type

	// Stream the object to the client
	_, err = io.Copy(c.Writer, object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
}
