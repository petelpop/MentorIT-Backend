package modulecontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetModules(c *gin.Context) {
	classID := c.Param("class_id")

	var modules []models.Module
	if err := config.DB.Preload("SubMods", func() *gorm.DB {
			return config.DB.Order("sub_modules.order ASC")
		}).
		Where("class_id = ?", classID).
		Order("modules.order ASC").
		Find(&modules).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}

	c.JSON(200, models.Response{
		Message: "Loaded modules successfully",
		Data: modules,
	})
}

func GetModuleDetail(c *gin.Context) {
	id := c.Param("id")

	var module models.Module
	if err := config.DB.
		Preload("SubMods", func(db *gorm.DB) *gorm.DB {
			return db.Order("sub_modules.order ASC")
		}).
		First(&module, id).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "Module not found"})
		return
	}

	c.JSON(200, models.Response{
		Message: "Module loaded successfully",
		Data: module,
	})
}

func CreateModule(c *gin.Context) {
	var module models.Module

	if err := c.ShouldBindJSON(&module); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error()})
		return
	}

	if module.Title == "" {
		c.JSON(400, models.Response{
			Message: "title, content, and description cannot be empty",
		})
		return
	}

	var maxOrder int
	config.DB.Model(&models.Module{}).
		Where("class_id = ?", module.ClassID).
		Select("COALESCE(MAX(`order`), 0)").
		Scan(&maxOrder)

	module.Order = maxOrder + 1 

	if err := config.DB.Create(&module).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}
	c.JSON(200, models.Response{
		Message: "Module created successfully",
	})
}

func UpdateModule(c *gin.Context) {
	id := c.Param("id")

	var module models.Module
	
	if err := config.DB.First(&module, id).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "Module not found"})
		return
	}

	if err := c.ShouldBindJSON(&module); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error()})
		return
	}

	if module.Title == "" {
		c.JSON(400, models.Response{
			Message: "title, content, and description cannot be empty"})
		return
	}

	if err := config.DB.Save(&module).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}
	c.JSON(200, models.Response{
		Message: "Module update successfully",
	})
}

func DeleteModule(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Module{}, id).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}
	c.JSON(200, models.Response{
		Message: "Module deleted successfully"})
}