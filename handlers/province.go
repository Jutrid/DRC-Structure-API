package handlers

import (
	"net/http"

	"Jutrid/DRC_structure_API/database"
	"Jutrid/DRC_structure_API/models"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/provinces
func GetProvinces(c *gin.Context) {
	page, limit, offset := GetPagination(c)
	search := c.Query("search")

	query := database.DB.Model(&models.Province{})
	countQuery := database.DB.Model(&models.Province{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR principal_town LIKE ?", like, like)
		countQuery = countQuery.Where("name LIKE ? OR principal_town LIKE ?", like, like)
	}

	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors du comptage")
		return
	}

	var provinces []models.Province
	if err := query.Order("id ASC").Offset(offset).Limit(limit).Find(&provinces).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors de la récupération des provinces")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": provinces,
		"meta": BuildMeta(page, limit, total),
	})
}

// GET /api/v1/provinces/:id
func GetProvinceByID(c *gin.Context) {
	id := c.Param("id")
	var province models.Province
	if err := database.DB.First(&province, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "province non trouvée")
		return
	}

	// Compte les villes et territoires liés
	var cityCount, territoryCount int64
	database.DB.Model(&models.City{}).Where("province_id = ?", province.ID).Count(&cityCount)
	database.DB.Model(&models.Territory{}).Where("province_id = ?", province.ID).Count(&territoryCount)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"province":         province,
			"cities_count":      cityCount,
			"territories_count": territoryCount,
		},
	})
}

// GET /api/v1/provinces/:id/cities
func GetCitiesByProvince(c *gin.Context) {
	id := c.Param("id")

	// Vérifie existence province
	var province models.Province
	if err := database.DB.First(&province, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "province non trouvée")
		return
	}

	page, limit, offset := GetPagination(c)
	search := c.Query("search")

	query := database.DB.Where("province_id = ?", province.ID)
	countQuery := database.DB.Model(&models.City{}).Where("province_id = ?", province.ID)

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ?", like)
		countQuery = countQuery.Where("name LIKE ?", like)
	}

	var total int64
	countQuery.Count(&total)

	var cities []models.City
	if err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&cities).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors de la récupération des villes")
		return
	}
	// Nettoie la relation Province pour éviter un objet vide (non-nécessaire ici, province déjà connue)
	for i := range cities {
		cities[i].Province = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"data": cities,
		"meta": BuildMeta(page, limit, total),
	})
}

// GET /api/v1/provinces/:id/territories
func GetTerritoriesByProvince(c *gin.Context) {
	id := c.Param("id")

	var province models.Province
	if err := database.DB.First(&province, id).Error; err != nil {
		ErrorResponse(c, http.StatusNotFound, "province non trouvée")
		return
	}

	page, limit, offset := GetPagination(c)
	search := c.Query("search")

	query := database.DB.Where("province_id = ?", province.ID)
	countQuery := database.DB.Model(&models.Territory{}).Where("province_id = ?", province.ID)

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ?", like)
		countQuery = countQuery.Where("name LIKE ?", like)
	}

	var total int64
	countQuery.Count(&total)

	var territories []models.Territory
	if err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&territories).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "erreur lors de la récupération des territoires")
		return
	}
	for i := range territories {
		territories[i].Province = nil
	}

	c.JSON(http.StatusOK, gin.H{
		"data": territories,
		"meta": BuildMeta(page, limit, total),
	})
}
