package authcontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/helper"
	"MentorIT-Backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	Id           uint      `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Exp          int       `json:"exp"`
	AccessToken  string    `json:"access_token`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	hashedPassword, _ := helper.HashPassword(input.Password)

	user := models.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     string(helper.Student),
		Exp:      0,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(400, models.Response{
			Message: "Email or Username already exists",
		})
		return
	}

	accessToken, refreshToken, err := helper.GenerateTokens(user.Id, user.Role)

	if err != nil {
		c.JSON(500, models.Response{
			Message: "Token generation failed",
		})
		return
	}

	token := models.Token{
		UserID:       user.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}
	config.DB.Create(&token)

	c.JSON(200, models.Response{
		Message: "user registered successfully",
		Data: UserResponse{
			Id:           user.Id,
			Name:         user.Name,
			Username:     user.Username,
			Email:        user.Email,
			Role:         user.Role,
			Exp:          user.Exp,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    token.ExpiresAt,
		},
	})
}

func Login(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(401, models.Response{
			Message: "Invalid username or password",
		})
		return
	}

	if !helper.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(401, models.Response{
			Message: "Invalid username or password",
		})
		return
	}

	accessToken, refreshToken, err := helper.GenerateTokens(user.Id, user.Role)

	if err != nil {
		c.JSON(500, models.Response{
			Message: "Could not generate token",
		})
		return
	}

	var token models.Token
	if err := config.DB.Where("user_id = ?", user.Id).First(&token).Error; err != nil {
		token.AccessToken = accessToken
		token.RefreshToken = refreshToken
		token.ExpiresAt = time.Now().Add(15 * time.Minute)
		config.DB.Save(&token)
	} else {
		token = models.Token{
			UserID:       user.Id,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Add(15 * time.Minute),
		}
		config.DB.Create(&token)
	}

	c.JSON(200, models.Response{
		Message: "login successful",
		Data: UserResponse{
			Id:           user.Id,
			Username:     user.Username,
			Email:        user.Email,
			Role:         user.Role,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.ExpiresAt,
		},
	})
}

func RefreshToken(c *gin.Context) {
	var inputToken models.Token

	if err := c.ShouldBindJSON(&inputToken); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	var token models.Token
	if err := config.DB.Preload("User").Where("refresh_token = ?", inputToken.RefreshToken).First(&token).Error; err != nil {
		c.JSON(401, models.Response{
			Message: "Invalid refresh token",
		})
		return
	}

	accessToken, newRefreshToken, err := helper.GenerateTokens(token.UserID, token.User.Role)

	if err != nil {
		c.JSON(500, models.Response{
			Message: "Token generation failed",
		})
		return
	}

	token.AccessToken = accessToken
	token.RefreshToken = newRefreshToken
	token.ExpiresAt = time.Now().Add(15 * time.Minute)
	config.DB.Save(&token)

	c.JSON(200, models.Response{
		Message: "token refreshed successfully",
		Data: TokenResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    token.ExpiresAt,
		},
	})
}

func Logout(c *gin.Context) {
	var inputToken models.Token

	if err := c.ShouldBindJSON(&inputToken); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	var token models.Token
	if err := config.DB.Where("access_token = ?", inputToken.AccessToken).First(&token).Error; err != nil {
		c.JSON(401, models.Response{
			Message: "Invalid refresh token",
		})
		return
	}

	config.DB.Delete(&token)

	c.JSON(200, models.Response{
		Message: "logout successful",
	})
}
