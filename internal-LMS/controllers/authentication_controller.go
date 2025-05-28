package controllers

import (
	_ "github.com/asaskevich/govalidator"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"learningManagementSystem/internal-LMS/database"
	models2 "learningManagementSystem/internal-LMS/models"
	"net/http"
	"time"
)

const JwtSecret = "your-secret-key"
const RefreshSecret = "refresh-secret-key"

func Validatepassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

// HashPassword Hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash Check password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUserRecord create user record
func CreateUserRecord(c iris.Context) {
	var user models2.User
	err := c.ReadJSON(&user)
	if err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid JSON. Expected an array of responses."})
		return
	}

	// Validate required fields
	if user.Name == "" || user.Password == "" || (user.Role != 1 && user.Role != 2) {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"Error": "Missing required fields: UserName, Password, or Role"})
		return
	}
	//Validate Password
	if !Validatepassword(user.Password) {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Password must be at least 8 characters long"})
		return
	}
	var userExists bool

	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE name = ? and role=?)",
		user.Name, user.Role).Scan(&userExists).Error
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "failed to verify user: " + err.Error()})
		return
	}
	if userExists {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "user account already exist "})
		return
	}
	// Insert  answers into User table
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		// Hash the password
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword

		// Insert the user into the database
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return nil

	})
	if txErr != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "failed to create user: " + txErr.Error()})
		return
	}
	c.StatusCode(iris.StatusCreated)
	c.JSON(iris.Map{"message": "User created successfully"})

}

// GenerateToken Generate JWT token
func GenerateToken(user models2.User) (string, error) {
	signer := jwt.NewSigner(jwt.HS256, JwtSecret, 1*time.Minute)
	claims := &models2.Claims{

		Username: user.Name,
		Userid:   user.ID,
		Role:     user.Role,
		//Exp:      int64(1 * time.Minute),
	}

	token, err := signer.Sign(claims)
	if err != nil {
		return " ", err
	}

	return string(token), nil
}

// GenerateRefreshToken Generate Refresh Token
func GenerateRefreshToken(user models2.User) (signedtoken string, err error) {
	signer := jwt.NewSigner(jwt.HS256, RefreshSecret, 24*time.Hour)
	claims := &models2.Claims{
		Username: user.Name,
		Userid:   user.ID,
	}

	token, err := signer.Sign(claims)
	if err != nil {

		return " ", err
	}

	return string(token), nil
}

func CheckAuthentication(c iris.Context) {

	var user models2.Userlogin
	if err := c.ReadJSON(&user); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid input"})
		return
	}
	var storedUser models2.User
	err := database.DB.Where("name = ?", user.Name).First(&storedUser).Error
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "No user found"})
		return
	}
	// Check if the password is correct
	if !CheckPasswordHash(user.Password, storedUser.Password) {
		c.StatusCode(iris.StatusUnauthorized)
		c.JSON(iris.Map{"error": "Invalid credentials"})
		return
	}
	// Check if the user is an instructor
	if storedUser.Role != 1 && storedUser.Role != 2 {
		c.StatusCode(iris.StatusUnauthorized)
		c.JSON(iris.Map{"error": "Unauthorized"})
		return

	}
	// Generate a JWT token
	token, err := GenerateToken(storedUser)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "Failed to generate token"})
		return
	}
	// Generate a refresh token
	t1, err := GenerateRefreshToken(storedUser)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "Failed to generate refresh token"})
		return
	}

	// Set JWT in a cookie
	cookies := []*http.Cookie{
		{
			Name:     "jwt_token",
			Value:    token,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   24 * 3600,
			Path:     "/",
		},
		{
			Name:     "refresh_token",
			Value:    t1,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   24 * 3600,
			Path:     "/",
		},
	}

	for _, cookie := range cookies {
		c.SetCookie(cookie, iris.CookieAllowSubdomains())
	}
	c.JSON(models2.LoginResponseJWT{
		Token:        token,
		RefreshToken: t1,
		Message:      "Login Successful",
		Username:     storedUser.Name,
	})

}

func Logout(ctx iris.Context) {
	token := ctx.GetCookie("jwt_token")
	if token == "" {
		// No token found, but still respond (user might already be logged out)
		ctx.StatusCode(iris.StatusOK)
		ctx.JSON(iris.Map{
			"message": "No active session found, logged out",
		})
		return
	}
	ctx.RemoveCookie("jwt_token", iris.CookiePath("/"))
	ctx.RemoveCookie("refresh_token", iris.CookiePath("/"))
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"message": "Logged out successfully",
	})

}

func ResetPassword(c iris.Context) {

	var req models2.ResetRequest
	if err := c.ReadJSON(&req); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid input"})
		return
	}
	if !Validatepassword(req.NewPassword) {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Password must be at least 8 characters long"})
		return
	}
	var user models2.User
	err := database.DB.Where("name = ?", req.Username).First(&user).Error
	if err != nil {
		c.StatusCode(iris.StatusNotFound)
		c.JSON(iris.Map{"error": "User not found"})
		return
	}
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "Failed to update password"})
		return
	}
	c.JSON(iris.Map{"message": "Password reset successful"})
}

func HasAccess(userID int, accessName string) bool {
	var count int
	err := database.DB.Table("user").Joins("inner join roles on user.role= roles.roles_id").
		Joins("inner join access_roles on roles.roles_id = access_roles.roles_id").
		Joins("inner join access on access_roles.access_id = access.access_id").
		Where("access_roles.access_id = ?", accessName).Where("users.id = ?", userID).Select("count(*)").Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0

}
