package handlers

import (
	"net/http"

	"Jutrid/DRC_structure_API/database"
	"Jutrid/DRC_structure_API/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/cities
func GetCities(c *gin.Context) {
	page, limit, offset := GetPagination(c)
	search := c.Query("search")
	provinceID := c.Query("province_id")

	query := database.DB.Model(&models.City{}).Preload("Province")
	countQuery := database.DB.Model(&models.City{})

	if provinceID != "" {
		query = query.Where("cities.province_id = ?", provinceID)
		countQuery = countQuery.Where("province_id = ?", provinceID)
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("cities.name LIKE ?", like)
		countQuery = countQuery.Where("name LIKE ?", like)
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors du comptage des villes")
		return
	}

	var cities []models.City
	if err := query.Order("cities.id ASC").Offset(offset).Limit(limit).Find(&cities).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors de la récupération des villes")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cities,
		"meta": BuildMeta(page, limit, total),
	})
}

// GET /api/v1/cities/:id
func GetCityByID(c *gin.Context) {
	id := c.Param("id")
	var city models.City
	if err := database.DB.Preload("Province").First(&city, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "ville non trouvée")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": city,
	})
}
