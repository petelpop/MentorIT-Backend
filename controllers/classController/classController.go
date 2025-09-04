package classcontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context)  {
	var classes []models.Class

	config.DB.Find(&classes)
	c.JSON(200, models.Response{
		Message: "Data successfully loaded",
		Data: classes,
	})
}

func Show(c *gin.Context)  {
	
}

func IndexByCategory(c *gin.Context)  {
	
}


// Admin, Teacher
func Create(c *gin.Context) {
	var input models.Class

	title := c.PostForm("title")
	description := c.PostForm("description")
	categoryName := c.PostForm("category_name")

	if title == "" || description == "" || categoryName == "" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Title, Description, and Category Name are required",
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

	openedFile, _ := file.Open()
	defer openedFile.Close()
	buffer := make([]byte, 512)
	_, _ = openedFile.Read(buffer)
	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Only JPEG and PNG images are allowed",
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
		Message: "Class created successfully",
		Data:    input,
	})
}

func Update(c *gin.Context) {
	id := c.Param("id")

	var class models.Class
	if err := config.DB.First(&class, id).Error; err != nil {
		c.AbortWithStatusJSON(404, models.Response{
			Message: "class not found",
		})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	categoryName := c.PostForm("category_name")

	if title != "" {
		class.Title = title
	}

	if description != "" {
		class.Description = description
	}

	if categoryName != "" {
		var category models.ClassCategory
		if err := config.DB.Where("name = ?", categoryName).First(&category).Error; err != nil {
			c.AbortWithStatusJSON(404, models.Response{
				Message: "Category not found",
			})
			return
		}
		class.CategoryName = categoryName
		class.ClassCategoryID = category.Id
	}

	file, err := c.FormFile("thumbnail")
	if err == nil {
		if file.Size > 5*1024*1024 {
			c.AbortWithStatusJSON(400, models.Response{
				Message: "File size exceeds 5MB",
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			c.AbortWithStatusJSON(400, models.Response{
				Message: "Only PNG, JPG, or JPEG files are allowed",
			})
			return
		}

		openedFile, _ := file.Open()
		defer openedFile.Close()
		buffer := make([]byte, 512)
		_, _ = openedFile.Read(buffer)
		contentType := http.DetectContentType(buffer)
		if contentType != "image/jpeg" && contentType != "image/png" {
			c.AbortWithStatusJSON(400, models.Response{
				Message: "Only JPEG and PNG images are allowed",
			})
			return
		}

		filename := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "Failed to save thumbnail",
			})
			return
		}
		class.Thumbnail = filename
	}

	if err := config.DB.Save(&class).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: "Failed to update class",
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Class update successfully",
		Data:    class,
	})

}
