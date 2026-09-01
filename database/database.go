package database

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"Jutrid/DRC_structure_API/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := DB.AutoMigrate(&models.Province{}, &models.City{}, &models.Territory{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := SeedIfEmpty(); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	return DB
}

func SeedIfEmpty() error {
	var count int64
	DB.Model(&models.Province{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	log.Println("Seeding database from Data_in_json...")

	// Resolve Data_in_json path relative to project root.
	// Try multiple locations: relative to cwd and relative to this file.
	candidates := []string{
		"Data_in_json",
		filepath.Join("..", "Data_in_json"),
		filepath.Join(".", "Data_in_json"),
	}
	var dataDir string
	for _, c := range candidates {
		if _, err := os.Stat(filepath.Join(c, "provinces.json")); err == nil {
			dataDir = c
			break
		}
	}
	if dataDir == "" {
		// Fallback: search from executable directory or current directory
		cwd, _ := os.Getwd()
		dataDir = filepath.Join(cwd, "Data_in_json")
	}

	if err := seedProvinces(filepath.Join(dataDir, "provinces.json")); err != nil {
		return err
	}
	if err := seedCities(filepath.Join(dataDir, "cities.json")); err != nil {
		return err
	}
	if err := seedTerritories(filepath.Join(dataDir, "territories.json")); err != nil {
		return err
	}

	log.Println("Database seeded successfully")
	return nil
}

func seedProvinces(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var provinces []models.Province
	if err := json.Unmarshal(data, &provinces); err != nil {
		return err
	}
	if len(provinces) == 0 {
		return nil
	}
	return DB.Create(&provinces).Error
}

func seedCities(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// JSON has extra field province_name that we ignore
	type rawCity struct {
		ID         uint    `json:"id"`
		ProvinceID uint    `json:"province_id"`
		Name       string  `json:"name"`
		Code       *string `json:"code"`
	}
	var raw []rawCity
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	cities := make([]models.City, 0, len(raw))
	for _, r := range raw {
		cities = append(cities, models.City{
			ID:         r.ID,
			ProvinceID: r.ProvinceID,
			Name:       r.Name,
			Code:       r.Code,
		})
	}
	return DB.Create(&cities).Error
}

func seedTerritories(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	type rawTerr struct {
		ID         uint    `json:"id"`
		ProvinceID uint    `json:"province_id"`
		Name       string  `json:"name"`
		Code       *string `json:"code"`
	}
	var raw []rawTerr
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	terrs := make([]models.Territory, 0, len(raw))
	for _, r := range raw {
		terrs = append(terrs, models.Territory{
			ID:         r.ID,
			ProvinceID: r.ProvinceID,
			Name:       r.Name,
			Code:       r.Code,
		})
	}
	return DB.Create(&terrs).Error
}
