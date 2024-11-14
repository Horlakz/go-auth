package model

import (
	"github.com/google/uuid"

	"github.com/horlakz/go-auth/database"
)

var (
	ROLE_INDIVIDUAL = "individual"
	ROLE_JOINT      = "joint"
	ROLE_CORPORATE  = "corporate"
)

type User struct {
	database.BaseModel

	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`

	Roles []UserRole `json:"roles"`
}

type Role struct {
	database.BaseModel

	Name string `json:"name"`
}

type UserRole struct {
	database.BaseModel

	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}

type VerificationCode struct {
	database.BaseModel

	UserID  uuid.UUID `json:"user_id"`
	Code    string    `json:"code"`
	Purpose string    `json:"purpose"`
}
