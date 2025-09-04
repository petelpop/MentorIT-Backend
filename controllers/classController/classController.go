package classcontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Admin, Teacher
func Create(c *gin.Context) {
	var input models.Class

	title := c.PostForm("title")
	description := c.PostForm("description")
	categoryName := c.PostForm("category_name")
	
	if title == "" || description == "" || categoryName == "" {
		c.AbortWithStatusJSON(400, models.Response{
			Message : "Title, Description, and Category Name are required",
		})
		return
	}

	var category models.ClassCategory
	if err := config.DB.Where("name = ?", categoryName).First(&category).Error; err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Category not found",
		})
		return
	}

	file, err := c.FormFile("thumbnail")
	if err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Thumbnail is required",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Only .jpg, .jpeg, and .png files are allowed",
		})
		return
	}

	filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: "Failed to upload thumbnail",
		})
		return
	}

	input.Title = title
	input.Description = description
	input.Thumbnail = filename
	input.CategoryName = categoryName
	input.ClassCategoryID = category.Id

	if err := config.DB.Create(&input).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.Response{
		Message : "Class created successfully",
		Data: input,
	})
}