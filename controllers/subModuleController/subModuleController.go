package submodulecontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"

	"github.com/gin-gonic/gin"
)

func GetSubModules(c *gin.Context) {
	moduleID := c.Param("module_id")

	var submodules []models.SubModule
	if err := config.DB.
		Where("module_id = ?", moduleID).
		Order("sub_modules.order ASC").
		Find(&submodules).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}

	c.JSON(200, models.Response{
		Message: "SubModules loaded successfully",
		Data: submodules,
	})
}

func GetSubModuleDetail(c *gin.Context) {
	id := c.Param("id")

	var submodule models.SubModule
	if err := config.DB.First(&submodule, id).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "SubModule not found"})
		return
	}

	c.JSON(200, models.Response{
		Message: "SubModule loaded successfully",
		Data: submodule,
	})
}

func CreateSubModule(c *gin.Context) {
	var subModule models.SubModule

	if err := c.ShouldBindJSON(&subModule); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error()})
		return
	}

	if subModule.Title == "" || subModule.Content == "" || subModule.Description == "" {
		c.JSON(400, models.Response{
			Message: "title, content, and description cannot be empty"})
		return
	}

	var maxOrder int
	config.DB.Model(&models.SubModule{}).
		Where("module_id = ?", subModule.ModuleID).
		Select("COALESCE(MAX(`order`), 0)").
		Scan(&maxOrder)

	subModule.Order = maxOrder + 1

	if err := config.DB.Create(&subModule).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}

	c.JSON(200, models.Response{
		Message: "SubModule created successfully",
		Data:    subModule})
}

func UpdateSubModule(c *gin.Context) {
	id := c.Param("id")
	var subModule models.SubModule

	if err := config.DB.First(&subModule, id).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "SubModule not found"})
		return
	}

	if err := c.ShouldBindJSON(&subModule); err != nil {
		c.JSON(400, models.Response{
			Message: err.Error()})
		return
	}

	if subModule.Title == "" || subModule.Content == "" || subModule.Description == "" {
		c.JSON(400, models.Response{
			Message: "title, content, and description cannot be empty"})
		return
	}

	if err := config.DB.Save(&subModule).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}

	c.JSON(200, models.Response{
		Message: "SubModule updated successfully",
		Data:    subModule})
}

func DeleteSubModule(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.SubModule{}, id).Error; err != nil {
		c.JSON(500, models.Response{
			Message: err.Error()})
		return
	}

	c.JSON(200, models.Response{
		Message: "SubModule deleted successfully"})
}
