package routes

import (
	"net/http"

	"Jutrid/DRC_structure_API/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware CORS simple
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Health / Welcome
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Bienvenue sur l'API DRC Structure 🇨🇩",
			"version": "v1",
			"endpoints": gin.H{
				"provinces":            "/api/v1/provinces",
				"province_detail":      "/api/v1/provinces/:id",
				"province_cities":      "/api/v1/provinces/:id/cities",
				"province_territories": "/api/v1/provinces/:id/territories",
				"cities":               "/api/v1/cities?province_id=&search=&page=&limit=",
				"city_detail":          "/api/v1/cities/:id",
				"territories":          "/api/v1/territories?province_id=&search=&page=&limit=",
				"territory_detail":     "/api/v1/territories/:id",
				"stats":                "/api/v1/stats",
			},
			"docs": "Voir README.md pour la documentation complète",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		// Provinces
		v1.GET("/provinces", handlers.GetProvinces)
		v1.GET("/provinces/:id", handlers.GetProvinceByID)
		v1.GET("/provinces/:id/cities", handlers.GetCitiesByProvince)
		v1.GET("/provinces/:id/territories", handlers.GetTerritoriesByProvince)

		// Cities
		v1.GET("/cities", handlers.GetCities)
		v1.GET("/cities/:id", handlers.GetCityByID)

		// Territories
		v1.GET("/territories", handlers.GetTerritories)
		v1.GET("/territories/:id", handlers.GetTerritoryByID)

		// Stats
		v1.GET("/stats", handlers.GetStats)
	}

	return r
}
