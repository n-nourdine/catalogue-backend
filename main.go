package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lib/pq" // Utilisation explicite
)

// --- 1. MODÈLE ---
type Spec struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// On ajoute des tags DB pour scanner le SQL
type Product struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Specs       []Spec   `json:"specs"`
	Images      []string `json:"images"`
}

// Structure intermédiaire pour décoder le JSONB de Postgres
type ProductDB struct {
	ID          string
	Title       string
	Category    string
	Description string
	Price       int
	Specs       []uint8 // Postgres renvoie le JSONB en []byte
	Images      []uint8 // Postgres renvoie le Array en string bizarre parfois, on va gérer ça
}

var db *sql.DB

// --- 2. CONNEXION DB ---
func initDB() {
	var err error
	// On récupère l'URL de la base depuis les variables d'environnement (Sécurité !)
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Erreur: La variable DATABASE_URL est vide")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Impossible de se connecter à Supabase:", err)
	}
	fmt.Println("✅ Connecté à Supabase PostgreSQL")
}

// --- 3. HANDLERS ---

func handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// Headers CORS pour autoriser ton Frontend à parler au Backend
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, category, description, price, specs, images FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		var specsJSON []byte
		// var imagesArray []string

		// Postgres driver 'pq' gère bien les tableaux de strings avec pq.Array()
		// Mais ici on va faire simple.
		// Note: Pour specs (JSONB), on scanne dans un []byte puis on Unmarshal
		// Pour images (TEXT[]), le driver pq a une fonction spéciale mais restons standard.

		// Simplification pour l'exemple : On utilise une astuce pour lire le tableau Postgres
		// Idéalement on utilise `pq.Array(&p.Images)`

		// if err := rows.Scan(&p.ID, &p.Title, &p.Category, &p.Description, &p.Price, &specsJSON, (*pqStringArray)(&p.Images)); err != nil {
		// 	log.Println("Erreur scan:", err)
		// 	continue
		// }
		if err := rows.Scan(&p.ID, &p.Title, &p.Category, &p.Description, &p.Price, &specsJSON, pq.Array(&p.Images)); err != nil {
			log.Println("Erreur scan:", err)
			continue
		}

		// Convertir le JSONB (Specs) en Struct Go
		json.Unmarshal(specsJSON, &p.Specs)

		products = append(products, p)
	}

	json.NewEncoder(w).Encode(products)
}

// Petit helper pour lire les tableaux Postgres (TEXT[]) sans prise de tête
type pqStringArray []string

func (a *pqStringArray) Scan(src interface{}) error {
	// Cette fonction sert à traduire le format {img1,img2} de postgres en []string Go
	// Pour l'instant, utilisons la librairie 'pq' qui le fait mieux.
	// Astuce : utilise "github.com/lib/pq" et rows.Scan(..., pq.Array(&p.Images))
	return nil
}

func main() {
	// Initialisation
	initDB()

	http.HandleFunc("/api/products", handleGetProducts)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("🚀 Serveur démarré sur le port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
