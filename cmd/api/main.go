package main

import (
	"catalogue-backend/internal/database"
	"catalogue-backend/internal/handlers"
	"catalogue-backend/internal/repository"
	"log"
	"net/http"
	"os"
)

func main() {
	// 1. Initialiser la DB
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Impossible de se connecter à la DB: ", err)
	}
	defer db.Close()

	// 2. Initialiser les couches (Repo -> Handler)
	repo := repository.NewProductRepository(db)
	handler := &handlers.ProductHandler{Repo: repo}

	// 3. Définir les routes
	http.HandleFunc("/api/products", handler.HandleProducts)

	// 4. Lancer le serveur
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Serveur démarré sur le port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
