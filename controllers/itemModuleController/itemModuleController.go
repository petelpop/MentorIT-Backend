package submodulecontroller

import (
	"MentorIT-Backend/config"
	"MentorIT-Backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ============ ADMIN ====================
func CreateModuleItem(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, models.Response{Message: "cannot read body"})
		return
	}

	var meta struct {
		ItemType string `json:"item_type" form:"item_type"`
	}
	_ = json.Unmarshal(bodyBytes, &meta)

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if meta.ItemType == "" && c.ContentType() == "multipart/form-data" {
		meta.ItemType = c.PostForm("item_type")
	}

	if meta.ItemType == "" {
		c.JSON(400, models.Response{Message: "item_type is required"})
		return
	}

	fmt.Println("Content-Type:", c.ContentType())
	fmt.Printf("Parsed ItemType: %+v\n", meta.ItemType)

	switch meta.ItemType {
	case "submodule":
		UploadSubModule(c)
	case "quiz":
		UploadQuiz(c)
	case "project":
		CreateProjectPage(c)
	default:
		c.JSON(400, models.Response{
			Message: "item_type must be one of: submodule, quiz, project",
		})
	}
}

// ----------------- SUBMODULE -----------------
func UploadSubModule(c *gin.Context) {
	var body struct {
		ModuleID    uint   `json:"module_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Content     string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "invalid body: " + err.Error()})
		return
	}

	if body.Title == "" || body.Description == "" || body.Content == "" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "title, description, content cannot be empty"})
		return
	}

	submodule := models.SubModule{
		Title:       body.Title,
		Description: body.Description,
		Content:     body.Content,
	}

	if err := config.DB.Create(&submodule).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{Message: err.Error()})
		return
	}

	var maxOrder int
	config.DB.Model(&models.ModuleItem{}).
		Where("module_id = ?", body.ModuleID).
		Select("COALESCE(MAX(`order`),0)").Scan(&maxOrder)

	item := models.ModuleItem{
		ModuleID: body.ModuleID,
		ItemType: "submodule",
		ItemID:   submodule.Id,
		Order:    maxOrder + 1,
	}
	config.DB.Create(&item)

	c.JSON(200, models.Response{
		Message: "SubModule uploaded",
		Data:    submodule})
}

// ----------------- QUIZ -----------------
func UploadQuiz(c *gin.Context) {
	var body struct {
		ModuleID  uint   `json:"module_id"`
		Title     string `json:"title"`
		Questions []struct {
			Question string   `json:"question"`
			Options  []string `json:"options"`
			Answer   string   `json:"answer"`
			Order    int      `json:"order"`
		} `json:"questions"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "invalid body: " + err.Error(),
		})
		return
	}

	if body.Title == "" || len(body.Questions) == 0 {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "title and questions required",
		})
		return
	}

	quiz := models.Quiz{
		Title:    body.Title,
		ModuleID: body.ModuleID,
	}
	if err := config.DB.Create(&quiz).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: err.Error(),
		})
		return
	}

	for i, q := range body.Questions {
		opts, _ := json.Marshal(q.Options) 

		question := models.QuizQuestion{
			QuizID:   quiz.Id,
			Question: q.Question,
			Options:  string(opts),
			Answer:   q.Answer,
			Order:    q.Order,
		}
		if question.Order == 0 {
			question.Order = i + 1
		}

		if err := config.DB.Create(&question).Error; err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "failed to save question: " + err.Error(),
			})
			return
		}
	}

	var maxOrder int
	config.DB.Model(&models.ModuleItem{}).
		Where("module_id = ?", body.ModuleID).
		Select("COALESCE(MAX(`order`),0)").
		Scan(&maxOrder)

	item := models.ModuleItem{
		ModuleID: body.ModuleID,
		ItemType: "quiz",
		ItemID:   quiz.Id,
		Order:    maxOrder + 1,
	}
	config.DB.Create(&item)

	var savedQuiz models.Quiz
	if err := config.DB.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("`order` ASC")
	}).First(&savedQuiz, quiz.Id).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: "failed to reload quiz: " + err.Error(),
		})
		return
	}

	type QuestionResponse struct {
		ID       uint     `json:"id"`
		QuizID   uint     `json:"quiz_id"`
		Question string   `json:"question"`
		Options  []string `json:"options"`
		Answer   string   `json:"answer"`
		Order    int      `json:"order"`
	}

	var questions []QuestionResponse
	for _, q := range savedQuiz.Questions {
		var opts []string
		_ = json.Unmarshal([]byte(q.Options), &opts)

		questions = append(questions, QuestionResponse{
			ID:       q.Id,
			QuizID:   q.QuizID,
			Question: q.Question,
			Options:  opts,
			Answer:   q.Answer,
			Order:    q.Order,
		})
	}

	c.JSON(200, models.Response{
		Message: "Quiz uploaded",
		Data: map[string]interface{}{
			"id":        savedQuiz.Id,
			"title":     savedQuiz.Title,
			"module_id": savedQuiz.ModuleID,
			"questions": questions,
		},
	})
}

// ----------------- FINAL PROJECT -----------------
func CreateProjectPage(c *gin.Context) {
	var body struct {
		ModuleID    uint   `json:"module_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Guide       string `json:"guide"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "invalid body: " + err.Error(),
		})
		return
	}

	if body.ModuleID == 0 || body.Title == "" || body.Description == "" {
		c.AbortWithStatusJSON(400, models.Response{
			Message: "module_id, title, and description are required",
		})
		return
	}

	project := models.ProjectPage{
		ModuleID:    body.ModuleID,
		Title:       body.Title,
		Description: body.Description,
		Guide:       body.Guide,
	}

	if err := config.DB.Create(&project).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: err.Error(),
		})
		return
	}

	var maxOrder int
	config.DB.Model(&models.ModuleItem{}).
		Where("module_id = ?", body.ModuleID).
		Select("COALESCE(MAX(`order`),0)").Scan(&maxOrder)

	item := models.ModuleItem{
		ModuleID: body.ModuleID,
		ItemType: "project",
		ItemID:   project.Id,
		Order:    maxOrder + 1,
	}
	config.DB.Create(&item)

	c.JSON(200, models.Response{
		Message: "Final project template created successfully",
		Data:    project,
	})
}

func DeleteModuleItem(c *gin.Context) {
	id := c.Param("id")

	var item models.ModuleItem
	if err := config.DB.First(&item, id).Error; err != nil {
		c.AbortWithStatusJSON(404, models.Response{
			Message: "Module item not found"})
		return
	}

	switch item.ItemType {
	case "submodule":
		if err := config.DB.Delete(&models.SubModule{}, item.ItemID).Error; err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "failed to delete submodule"})
			return
		}
	case "quiz":
		if err := config.DB.Where("quiz_id = ?", item.ItemID).Delete(&models.QuizQuestion{}).Error; err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "failed to delete quiz questions"})
			return
		}
		if err := config.DB.Delete(&models.Quiz{}, item.ItemID).Error; err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "failed to delete quiz"})
			return
		}
	case "project":
		if err := config.DB.Delete(&models.ProjectPage{}, item.ItemID).Error; err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "failed to delete project"})
			return
		}
	default:
		c.AbortWithStatusJSON(400, models.Response{
			Message: "unknown item type"})
		return
	}

	if err := config.DB.Delete(&item).Error; err != nil {
		c.AbortWithStatusJSON(500, models.Response{
			Message: "failed to delete module item"})
		return
	}

	c.JSON(200, models.Response{
		Message: "Module item deleted successfully",
	})
}

func UpdateModuleItem(c *gin.Context) {
	id := c.Param("id")

	var item models.ModuleItem
	if err := config.DB.First(&item, id).Error; err != nil {
		c.AbortWithStatusJSON(404, models.Response{
			Message: "Module item not found"})
		return
	}

	switch item.ItemType {
	case "submodule":
		var input models.SubModule
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(400, models.Response{
				Message: err.Error()})
			return
		}

		if input.Title == "" || input.Content == "" || input.Description == "" {
			c.AbortWithStatusJSON(400, models.Response{
				Message: "title, content, and description cannot be empty"})
			return
		}

		if err := config.DB.Model(&models.SubModule{}).
			Where("id = ?", item.ItemID).
			Updates(map[string]interface{}{
				"title":       input.Title,
				"content":     input.Content,
				"description": input.Description,
			}).Error; err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "failed to update submodule"})
			return
		}

	case "quiz":
		var input models.Quiz
		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(400, models.Response{
				Message: err.Error()})
			return
		}

		if err := config.DB.Model(&models.Quiz{}).
			Where("id = ?", item.ItemID).
			Updates(map[string]interface{}{
				"title": input.Title,
			}).Error; err != nil {
			c.AbortWithStatusJSON(500, models.Response{
				Message: "failed to update quiz"})
			return
		}

		if len(input.Questions) > 0 {
			config.DB.Where("quiz_id = ?", item.ItemID).Delete(&models.QuizQuestion{})
			for _, q := range input.Questions {
				q.QuizID = item.ItemID
				if err := config.DB.Create(&q).Error; err != nil {
					c.JSON(500, models.Response{
						Message: "failed to update quiz questions"})
					return
				}
			}
		}

	case "project":
		var input models.ProjectPage
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, models.Response{
				Message: err.Error()})
			return
		}

		if err := config.DB.Model(&models.ProjectPage{}).
			Where("id = ?", item.ItemID).
			Updates(map[string]interface{}{
				"title":       input.Title,
				"description": input.Description,
				"guide":       input.Guide,
			}).Error; err != nil {
			c.JSON(500, models.Response{
				Message: "failed to update project"})
			return
		}

	default:
		c.JSON(400, models.Response{
			Message: "unknown item type"})
		return
	}

	c.JSON(200, models.Response{
		Message: "Module item updated successfully",
	})
}

// ============ STUDENT ====================
func GetModuleItems(c *gin.Context) {
	moduleID := c.Param("id")
	fmt.Printf("Getting module items for module_id: %s\n", moduleID)

	var items []models.ModuleItem
	if err := config.DB.
		Where("module_id = ?", moduleID).
		Order("`order` ASC").
		Find(&items).Error; err != nil {
		fmt.Printf("Error getting module items: %s\n", err.Error())
		c.AbortWithStatusJSON(500, models.Response{
			Message: err.Error(),
		})
		return
	}

	fmt.Printf("Found %d module items\n", len(items))
	result := []map[string]interface{}{}

	for _, item := range items {
		fmt.Printf("Processing item: ID=%d, Type=%s, ItemID=%d\n", item.Id, item.ItemType, item.ItemID)
		entry := map[string]interface{}{
			"id":        item.Id,
			"item_type": item.ItemType,
			"order":     item.Order,
		}

		switch item.ItemType {
		case "submodule":
			var sm models.SubModule
			if err := config.DB.First(&sm, item.ItemID).Error; err == nil {
				entry["data"] = sm
				fmt.Printf("Successfully loaded submodule: %s\n", sm.Title)
			} else {
				fmt.Printf("Error loading submodule with ID %d: %s\n", item.ItemID, err.Error())
				entry["data"] = nil
			}
		case "quiz":
			var quiz models.Quiz
			if err := config.DB.Preload("Questions", func(db *gorm.DB) *gorm.DB {
				return db.Order("`order` ASC")
			}).First(&quiz, item.ItemID).Error; err == nil {
				entry["data"] = quiz
				fmt.Printf("Successfully loaded quiz: %s with %d questions\n", quiz.Title, len(quiz.Questions))
			} else {
				fmt.Printf("Error loading quiz with ID %d: %s\n", item.ItemID, err.Error())
				entry["data"] = nil
			}
		case "project":
			var project models.ProjectPage
			if err := config.DB.First(&project, item.ItemID).Error; err == nil {
				entry["data"] = project
				fmt.Printf("Successfully loaded project: %s\n", project.Title)
			} else {
				fmt.Printf("Error loading project with ID %d: %s\n", item.ItemID, err.Error())
				entry["data"] = nil
			}
		default:
			fmt.Printf("Unknown item type: %s\n", item.ItemType)
			entry["data"] = nil
		}

		result = append(result, entry)
	}

	c.JSON(200, models.Response{
		Message: "Module items loaded successfully",
		Data:    result,
	})
}

func GetModuleItemDetail(c *gin.Context) {
	id := c.Param("id")

	var item models.ModuleItem
	if err := config.DB.First(&item, id).Error; err != nil {
		c.AbortWithStatusJSON(404, models.Response{
			Message: "Module item not found"})
		return
	}

	entry := map[string]interface{}{
		"id":        item.Id,
		"item_type": item.ItemType,
		"order":     item.Order,
	}

	switch item.ItemType {
	case "submodule":
		var sm models.SubModule
		if err := config.DB.First(&sm, item.ItemID).Error; err == nil {
			entry["data"] = sm
		} else {
			entry["data"] = nil
		}
	case "quiz":
		var quiz models.Quiz
		if err := config.DB.Preload("Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC")
		}).First(&quiz, item.ItemID).Error; err == nil {
			entry["data"] = quiz
		} else {
			entry["data"] = nil
		}
	case "project":
		var project models.ProjectPage
		if err := config.DB.First(&project, item.ItemID).Error; err == nil {
			entry["data"] = project
		} else {
			entry["data"] = nil
		}
	default:
		entry["data"] = nil
	}

	c.JSON(200, models.Response{
		Message: "Module item loaded successfully",
		Data:    entry,
	})
}
