package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"LecCenterGo/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTopFiveLeaderboard(c *gin.Context) {

	var users []models.User
	// Order by 'count' in descending order and limit to 5
	result := dao.DB.Order("count DESC").Limit(5).Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to retrieve leaderboard data",
			"error":   result.Error.Error(),
		})
		return
	}

	var leaderboardData []response.LeaderboardResponse
	for _, user := range users {
		leaderboardData = append(leaderboardData, response.LeaderboardResponse{
			Username: user.Username,
			Avatar:   user.Avatar,
			Count:    user.Count,
			Major:    user.Major,
			Grade:    user.Grade,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Successfully retrieved top 5 leaderboard.",
		"data":    leaderboardData,
	})
}
