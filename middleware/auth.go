package middleware

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(requiredRoles ...string) gin.HandlerFunc{
	 return func(c *gin.Context) {
	 authHeader := c.GetHeader("Authorization")
	 if authHeader == "" {
		c.JSON(401, models.Response{
			Message: "Authorization header required"})
			c.Abort()
			return
			}
			
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				c.JSON(401, models.Response{
					Message: "Invalid token format",
				})
				c.Abort()
				return
			}

			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return config.JWTKey, nil
			})

			if err != nil || !token.Valid {
				c.JSON(401, models.Response{
					Message: "Invalid token",
				})
				c.Abort()
				return
			}

			allowed := false
			for _, role := range requiredRoles {
				if claims.Role == role {
					allowed = true
					break
				}
			}

			if !allowed {
				c.JSON(403, models.Response{
					Message: "Forbidden: insufficient permissions",
				})
				c.Abort()
				return
			}
			
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			
			c.Next()
	}
}