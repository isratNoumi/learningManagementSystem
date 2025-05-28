package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"type:int;primaryKey"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
}

type Userlogin struct {
	Name     string `json:"userName"`
	Password string `json:"userPassword"`
}
type LoginResponseJWT struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
	Username     string `json:"username"`
	Message      string `json:"message"`
}

type Claims struct {
	Userid   int    `json:"user_id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}
type ResetRequest struct {
	Username    string `json:"username"`
	NewPassword string `json:"new_password"`
}
