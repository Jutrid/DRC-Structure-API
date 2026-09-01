package handlers

import (
	"net/http"

	"Jutrid/DRC_structure_API/database"
	"Jutrid/DRC_structure_API/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/territories
func GetTerritories(c *gin.Context) {
	page, limit, offset := GetPagination(c)
	search := c.Query("search")
	provinceID := c.Query("province_id")

	query := database.DB.Model(&models.Territory{}).Preload("Province")
	countQuery := database.DB.Model(&models.Territory{})

	if provinceID != "" {
		query = query.Where("territories.province_id = ?", provinceID)
		countQuery = countQuery.Where("province_id = ?", provinceID)
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("territories.name LIKE ?", like)
		countQuery = countQuery.Where("name LIKE ?", like)
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors du comptage des territoires")
		return
	}

	var territories []models.Territory
	if err := query.Order("territories.id ASC").Offset(offset).Limit(limit).Find(&territories).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors de la récupération des territoires")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": territories,
		"meta": BuildMeta(page, limit, total),
	})
}

// GET /api/v1/territories/:id
func GetTerritoryByID(c *gin.Context) {
	id := c.Param("id")
	var territory models.Territory
	if err := database.DB.Preload("Province").First(&territory, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "territoire non trouvé")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": territory,
	})
}
