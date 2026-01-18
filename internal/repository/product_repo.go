package repository

import (
	"catalogue-backend/internal/models"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT id, title, category, description, price, specs, images FROM products`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		var specsJSON []byte // Variable temporaire pour le JSONB

		// Scan : On mappe les colonnes SQL vers nos variables
		// pq.Array(&p.Images) gère la conversion TEXT[] -> []string
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Category,
			&p.Description,
			&p.Price,
			&specsJSON,
			pq.Array(&p.Images),
		)
		if err != nil {
			return nil, err
		}

		// On décode le JSON brut des specs vers la structure Go
		if len(specsJSON) > 0 {
			_ = json.Unmarshal(specsJSON, &p.Specs)
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) Create(p models.Product) (models.Product, error) {
	query := `
		INSERT INTO products (title, category, description, price, specs, images)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	specsJSON, _ := p.SpecsToJSON()

	// QueryRow pour récupérer l'ID généré
	err := r.DB.QueryRow(
		query,
		p.Title,
		p.Category,
		p.Description,
		p.Price,
		specsJSON,
		pq.Array(p.Images),
	).Scan(&p.ID)

	if err != nil {
		return models.Product{}, err
	}
	return p, nil
}
