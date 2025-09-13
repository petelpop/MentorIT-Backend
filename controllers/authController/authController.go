package authcontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/helper"
	"MentorIT-Backend/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type ResetPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordWithTokenRequest struct {
	Email       string `json:"email" binding:"required,email"`
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
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
			Username:     user.Username,
			Name:         user.Name,
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
		token.ExpiresAt = time.Now().Add(1000000 * time.Hour)
		config.DB.Save(&token)
	} else {
		token = models.Token{
			UserID:       user.Id,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    time.Now().Add(1000000 * time.Hour),
		}
		config.DB.Create(&token)
	}

	c.JSON(200, models.Response{
		Message: "login successful",
		Data: UserResponse{
			Id:           user.Id,
			Username:     user.Username,
			Name:         user.Name,
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
	token.ExpiresAt = time.Now().Add(1000000 * time.Hour)
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

func ResetPassword(c *gin.Context) {
	var request ResetPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, models.Response{
			Message: "Unauthorized",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "User not found",
		})
		return
	}

	if !helper.CheckPasswordHash(request.OldPassword, user.Password) {
		c.JSON(400, models.Response{
			Message: "Invalid old password",
		})
		return
	}

	hashedPassword, err := helper.HashPassword(request.NewPassword)
	if err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to hash password",
		})
		return
	}

	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to update password",
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Password updated successfully",
	})
}

func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ForgotPassword(c *gin.Context) {
	var request ForgotPasswordRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(200, models.Response{
			Message: "If the email exists, a reset link has been sent",
		})
		return
	}

	resetToken, err := generateResetToken()
	if err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to generate reset token",
		})
		return
	}

	config.DB.Where("user_id = ? AND used = ?", user.Id, false).Delete(&models.ResetToken{})

	resetTokenRecord := models.ResetToken{
		UserID:    user.Id,
		Token:     resetToken,
		Email:     request.Email,
		ExpiresAt: time.Now().Add(1 * time.Hour), 
		Used:      false,
	}

	if err := config.DB.Create(&resetTokenRecord).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to create reset token",
		})
		return
	}

	emailService := helper.NewEmailService()
	if err := emailService.SendResetPasswordEmail(request.Email, user.Name, resetToken); err != nil {
		fmt.Printf("Failed to send reset email to %s: %v\n", request.Email, err)
	}

	c.JSON(200, models.Response{
		Message: "If the email exists, a reset link has been sent",
	})
}

func ResetPasswordWithToken(c *gin.Context) {
	var request ResetPasswordWithTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	var resetToken models.ResetToken
	if err := config.DB.Where("token = ? AND email = ? AND used = ? AND expires_at > ?",
		request.ResetToken, request.Email, false, time.Now()).First(&resetToken).Error; err != nil {
		c.JSON(400, models.Response{
			Message: "Invalid or expired reset token",
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, resetToken.UserID).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "User not found",
		})
		return
	}

	hashedPassword, err := helper.HashPassword(request.NewPassword)
	if err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to hash password",
		})
		return
	}

	if err := config.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to update password",
		})
		return
	}

	config.DB.Model(&resetToken).Update("used", true)

	config.DB.Where("user_id = ?", user.Id).Delete(&models.Token{})

	c.JSON(200, models.Response{
		Message: "Password reset successfully",
	})
}

func TestEmailConfig(c *gin.Context) {
	emailService := helper.NewEmailService()

	if err := emailService.TestEmailConfiguration(); err != nil {
		c.JSON(500, models.Response{
			Message: "Email configuration test failed",
			Data:    map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Email configuration is working properly",
	})
}
