package admincontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/helper"
	"MentorIT-Backend/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type CreateTeacherInput struct {
	Username string `json:"username" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func CreateTeacher(c *gin.Context) {
	var input CreateTeacherInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, models.Response{
			Message: "failed to hash password",
		})
		return
	}

	user := models.User{
		Username: input.Username,
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     string(helper.Teacher),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Teacher created successfully",
		"user":    user,
	})
}

func DeleteTeacher(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := config.DB.Where("id = ? AND role = ?", id, "teacher").First(&user).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "Teacher not found",
		})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "failed to delete teacher",
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Teacher deleted successfully",
	})
}
