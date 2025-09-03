package classcategorycontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"
	"strconv"

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

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	config.DB.Create(&input)
	c.JSON(200, models.Response{
		Message: "Success",
		Data:    input,
	})
}

func Update(c *gin.Context) {
	var category models.ClassCategory
	idParam := c.Param("id")

	if err := c.ShouldBindJSON(&category); err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: err.Error(),
		})
		return
	}

	if config.DB.Model(&category).Where("id = ?", idParam).Updates(category).RowsAffected == 0 {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Update failed",
		})
		return
	}

	idUint, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "Invalid ID format",
		})
		return
	}

	c.JSON(200, models.Response{
		Message: "Data successfully updated",
		Data: models.ClassCategory{
			Id:          uint(idUint),
			Name:        category.Name,
			Description: category.Description,
		},
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
