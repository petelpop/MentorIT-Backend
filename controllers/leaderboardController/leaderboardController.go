package leaderboardcontroller

import (
	"strconv"

	"MentorIT-Backend/config"
	"MentorIT-Backend/models"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var users []models.User
	if err := config.DB.
		Select("id, username, name, email, role, exp").
		Where("role = ?", "student").
		Order("exp DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to get leaderboard",
			Data:    nil,
		})
		return
	}

	leaderboard := make([]map[string]interface{}, len(users))
	for i, user := range users {
		leaderboard[i] = map[string]interface{}{
			"rank":     offset + i + 1,
			"id":       user.Id,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
			"exp":      user.Exp,
		}
	}

	c.JSON(200, models.Response{
		Message: "Leaderboard retrieved successfully",
		Data:    leaderboard,
	})
}

func GetUserRank(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(400, models.Response{
			Message: "Invalid user ID",
			Data:    nil,
		})
		return
	}

	var user models.User
	if err := config.DB.
		Select("id, username, name, email, role, exp").
		First(&user, userId).Error; err != nil {
		c.JSON(404, models.Response{
			Message: "User not found",
			Data:    nil,
		})
		return
	}

	if user.Role != "student" {
		c.JSON(400, models.Response{
			Message: "Only students are included in leaderboard rankings",
			Data:    nil,
		})
		return
	}

	var rank int64
	if err := config.DB.
		Model(&models.User{}).
		Where("exp > ? AND role = ?", user.Exp, "student").
		Count(&rank).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to calculate user rank",
			Data:    nil,
		})
		return
	}

	userRank := rank + 1

	var totalUsers int64
	config.DB.Model(&models.User{}).Where("role = ?", "student").Count(&totalUsers)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
			"exp":      user.Exp,
		},
		"rank":        userRank,
		"total_users": totalUsers,
		"percentile":  float64(totalUsers-userRank+1) / float64(totalUsers) * 100,
	}

	c.JSON(200, models.Response{
		Message: "User rank retrieved successfully",
		Data:    response,
	})
}

func GetTopUsers(c *gin.Context) {
	topStr := c.DefaultQuery("top", "5")
	top, err := strconv.Atoi(topStr)
	if err != nil || top <= 0 {
		top = 5
	}

	if top > 50 {
		top = 50
	}

	var users []models.User
	if err := config.DB.
		Select("id, username, name, email, role, exp").
		Where("role = ?", "student").
		Order("exp DESC").
		Limit(top).
		Find(&users).Error; err != nil {
		c.JSON(500, models.Response{
			Message: "Failed to get top users",
			Data:    nil,
		})
		return
	}

	topUsers := make([]map[string]interface{}, len(users))
	for i, user := range users {
		topUsers[i] = map[string]interface{}{
			"rank":     i + 1,
			"id":       user.Id,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
			"exp":      user.Exp,
		}
	}

	c.JSON(200, models.Response{
		Message: "Top users retrieved successfully",
		Data:    topUsers,
	})
}
