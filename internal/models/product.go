package models

import "encoding/json"

// Spec représente une ligne dans ton JSONB (ex: {"label": "Poids", "value": "45g"})
type Spec struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Product est l'objet principal
type Product struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Specs       []Spec   `json:"specs"`  // Mappé vers JSONB
	Images      []string `json:"images"` // Mappé vers TEXT[]
}

// Helper pour convertir les Specs en JSON string pour la DB
func (p *Product) SpecsToJSON() ([]byte, error) {
	return json.Marshal(p.Specs)
}