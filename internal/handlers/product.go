package handlers

import (
	"catalogue-backend/internal/models"
	"catalogue-backend/internal/repository"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	Repo *repository.ProductRepository
}

// EnableCORS ajoute les headers pour que le frontend puisse parler au backend
func (h *ProductHandler) EnableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	h.EnableCORS(&w)

	// Gérer la requête "preflight" des navigateurs
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		h.getProducts(w, r)
	case "POST":
		h.createProduct(w, r)
	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Erreur DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON Invalide", http.StatusBadRequest)
		return
	}

	newProduct, err := h.Repo.Create(p)
	if err != nil {
		http.Error(w, "Erreur création: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}
