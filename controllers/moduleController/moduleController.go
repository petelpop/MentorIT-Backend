package modulecontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"

	"github.com/gin-gonic/gin"
)

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