package classcategorycontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var categories []models.ClassCategory

	config.DB.Find(&categories)
	c.JSON(200, models.Response{
		Message: "Success",
		Data:    categories,
	})
}

func Show(c *gin.Context) {
	var category models.ClassCategory

	id := c.Param("id")

	if err := config.DB.First(&category, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(404, models.Response{
				Message: "Data not found",
			})
			return
		default:
			c.AbortWithStatusJSON(500, models.Response{
				Message: err.Error(),
			})
			return
		}
	}

	c.JSON(200, models.Response{
		Message: "Success",
		Data:    category,
	})
}

// Admin
func Create(c *gin.Context) {
	var input models.ClassCategory

	Name := c.PostForm("name")
	Description := c.PostForm("description")

	if Name == "" || Description == "" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Name and Description are required",
		})
		return
	}
	input.Name = Name
	input.Description = Description

	file, err := c.FormFile("icon")
	
	if err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Icon is required",
		})
		return
	}

	if file.Size > 5*1024*1024 {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "File size exceeds 5MB",
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Only .jpg, .jpeg, .png files are allowed",
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
			Message: "Failed to upload file",
		})
		return
	}
	input.Icon = filename

	if err := config.DB.Create(&input).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Success Created Data",
		Data:    input,
	})
}

func Update(c *gin.Context) {
	var category models.ClassCategory
	idParam := c.Param("id")

	if err := config.DB.First(&category, idParam).Error; err != nil {
		c.AbortWithStatusJSON(404, models.Response{
			Message: "Data not found",
		})
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	if name != "" {
		category.Name = name
	}
	if description != "" {
		category.Description = description
	}

	file, err := c.FormFile("icon")
	if err == nil {
		if file.Size > 5*1024*1024 {
			c.AbortWithStatusJSON(400, models.Response{
				Message: "File size exceeds 5MB",
			})
			return
		}
	

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Only .jpg, .jpeg, .png files are allowed",
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
			Message: "Failed to upload file",
		})
		return
	}
	category.Icon = filename
	}
	
	if err := config.DB.Save(&category).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Data successfully updated",
		Data:    category,
	})

}

func Delete(c *gin.Context) {
	var category models.ClassCategory

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Invalid ID",
		})
		return
	}

	if config.DB.Delete(&category, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Delete Failed",
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Data successfully deleted",
	})
}
