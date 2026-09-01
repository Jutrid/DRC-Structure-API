package main

import (
	"log"
	"os"

	"Jutrid/DRC_structure_API/database"
	"Jutrid/DRC_structure_API/routes"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database/db.sqlite"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialise DB (migration + seed)
	database.InitDB(dbPath)
	log.Printf("Base de données SQLite initialisée: %s", dbPath)

	// Setup router
	r := routes.SetupRouter()

	log.Printf("Serveur démarré sur http://localhost:%s", port)
	log.Printf("Endpoints disponibles:")
	log.Printf("  GET /api/v1/provinces")
	log.Printf("  GET /api/v1/provinces/:id")
	log.Printf("  GET /api/v1/provinces/:id/cities")
	log.Printf("  GET /api/v1/provinces/:id/territories")
	log.Printf("  GET /api/v1/cities")
	log.Printf("  GET /api/v1/cities/:id")
	log.Printf("  GET /api/v1/territories")
	log.Printf("  GET /api/v1/territories/:id")
	log.Printf("  GET /api/v1/stats")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erreur serveur: %v", err)
	}
}
