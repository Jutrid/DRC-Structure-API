package handlers

import (
	"net/http"

	"Jutrid/DRC_structure_API/database"
	"Jutrid/DRC_structure_API/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/stats
func GetStats(c *gin.Context) {
	var provinceCount, cityCount, territoryCount int64
	database.DB.Model(&models.Province{}).Count(&provinceCount)
	database.DB.Model(&models.City{}).Count(&cityCount)
	database.DB.Model(&models.Territory{}).Count(&territoryCount)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"provinces":   provinceCount,
			"cities":      cityCount,
			"territories": territoryCount,
			"total":       provinceCount + cityCount + territoryCount,
		},
	})
}
