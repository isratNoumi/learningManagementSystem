package middlewares

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"learningManagementSystem/internal-LMS/controllers"
	"learningManagementSystem/internal-LMS/models"
)

//var rolePermissions = map[string][]string{
//	"instructor": {"create", "read", "update", "delete"},
//	"user":       {"read"},
//}

// AddMiddleware sets up the CORS middleware for the Iris application.

func AddCorsMiddleware() iris.Handler {

	// Configure CORS to allow access from your frontend IP
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	return corsConfig

}

func AddCookieMiddleware() iris.Handler {

	// Wrapper middleware to extract JWT from cookie
	cookieJWTMiddleware := func(ctx iris.Context) {
		// Get JWT from cookie
		token := ctx.GetCookie("jwt_token")
		if token == "" {
			ctx.StatusCode(iris.StatusUnauthorized)
			_ = ctx.JSON(iris.Map{"error": "Missing JWT cookie"})
			return
		}

		// Set Authorization header for verifyMiddleware
		ctx.Request().Header.Set("Authorization", "Bearer "+token)

		ctx.Next()

	}
	return cookieJWTMiddleware
}

func VerifyMiddleware() iris.Handler {
	verifier := jwt.NewVerifier(jwt.HS256, controllers.JwtSecret)
	verifier.WithDefaultBlocklist()

	verifyMiddleware := verifier.Verify(func() interface{} {
		return new(models.Claims)
	})
	return verifyMiddleware

}

// RBACMiddleware checks if the user has the required permission
func RBACMiddleware(requiredPermission string) iris.Handler {
	return func(ctx iris.Context) {
		// Get the user role from the JWT claims
		claims := jwt.Get(ctx).(*models.Claims)
		userRole := claims.Role

		if userRole == 0 {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(iris.Map{"error": "User not authenticated"})
			return
		}

		exists := controllers.HasAccess(userRole, requiredPermission)
		if !exists {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{"error": "Role not found"})
			return
		}

		// If the user has the required permission, proceed to the next handler
		ctx.Next()
	}
}
