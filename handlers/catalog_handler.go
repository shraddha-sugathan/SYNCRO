package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"syncro/models"

	"github.com/google/uuid"
)

func HandleCreateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.ProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.BillingInterval != "monthly" && req.BillingInterval != "yearly" {
			http.Error(w, "Billing interval must be 'monthly' or 'yearly'", http.StatusBadRequest)
			return
		}

		var newID uuid.UUID
		err := db.QueryRow(
			"INSERT INTO products (name, price_cents, billing_interval) VALUES ($1, $2, $3) RETURNING id",
			req.Name, req.PriceCents, req.BillingInterval,
		).Scan(&newID)

		if err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product created", "id": newID.String()})
	}
}

func HandleGetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, price_cents, billing_interval, created_at FROM products")
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.PriceCents, &p.BillingInterval, &p.CreatedAt); err != nil {
				continue
			}
			products = append(products, p)
		}

		if products == nil {
			products = []models.Product{} // Return empty array instead of null
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}