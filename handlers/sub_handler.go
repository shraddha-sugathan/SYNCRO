package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"syncro/models"

	"github.com/google/uuid"
)

func HandleSubscribe(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.SubscribeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Extract isolated tenant/user ID directly from the JWT context
		userID := r.Context().Value("user_id").(uuid.UUID)

		var billingInterval string
		err := db.QueryRow("SELECT billing_interval FROM products WHERE id = $1", req.ProductID).Scan(&billingInterval)
		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		now := time.Now()
		var endDate time.Time
		if billingInterval == "monthly" {
			endDate = now.AddDate(0, 1, 0)
		} else {
			endDate = now.AddDate(1, 0, 0)
		}

		var subID uuid.UUID
		err = db.QueryRow(
			`INSERT INTO subscriptions (user_id, product_id, status, current_period_start, current_period_end) 
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			userID, req.ProductID, "active", now, endDate,
		).Scan(&subID)

		if err != nil {
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":         "Subscribed successfully",
			"subscription_id": subID.String(),
		})
	}
}
