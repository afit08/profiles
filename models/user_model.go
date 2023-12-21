package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	ID        string                `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string                `form:"name,omitempty" bson:"name,omitempty"`
	Image     *multipart.FileHeader `form:"image,omitempty" bson:"image,omitempty"`
	Desc      string                `form:"desc,omitempty" bson:"desc,omitempty"`
	JobName   string                `form:"job_name,omitempty" bson:"job_name,omitempty"`
	Skills    []string              `form:"skills,omitempty" bson:"skills,omitempty"`
	Roles     string                `form:"roles,omitempty" bson:"roles,omitempty"`
	Username  string                `form:"username,omitempty" bson:"username,omitempty"`
	Password  string                `form:"password,omitempty" bson:"password,omitempty"`
	CreatedAt time.Time             `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
