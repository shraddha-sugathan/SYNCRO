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
		w.Header().Set("Content-Type", "application/json")

		// Decode the incoming request
		var req models.SubscribeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate the product ID
		if req.ProductID == uuid.Nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Safely extract the user ID from JWT context
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok || userID == uuid.Nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Retrieve the product's billing interval
		var billingInterval string
		err := db.QueryRow(
			"SELECT billing_interval FROM products WHERE id = $1",
			req.ProductID,
		).Scan(&billingInterval)

		if err == sql.ErrNoRows {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Calculate the subscription period
		now := time.Now()
		var endDate time.Time

		switch billingInterval {
		case "monthly":
			endDate = now.AddDate(0, 1, 0)
		case "yearly":
			endDate = now.AddDate(1, 0, 0)
		default:
			http.Error(w, "Unsupported billing interval", http.StatusBadRequest)
			return
		}

		// Create the subscription
		var subID uuid.UUID
		err = db.QueryRow(
			`INSERT INTO subscriptions
			 (user_id, product_id, status, current_period_start, current_period_end)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id`,
			userID, req.ProductID, "active", now, endDate,
		).Scan(&subID)

		if err != nil {
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}

		// Return the created subscription
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":         "Subscribed successfully",
			"subscription_id": subID.String(),
		})
	}
}

func HandleGetSubscriptions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Safely extract the user ID
		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok || userID == uuid.Nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(
			`SELECT id, user_id, product_id, status,
			        current_period_start, current_period_end, created_at
			 FROM subscriptions
			 WHERE user_id = $1`,
			userID,
		)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var subs []models.Subscription

		for rows.Next() {
			var s models.Subscription

			if err := rows.Scan(
				&s.ID,
				&s.UserID,
				&s.ProductID,
				&s.Status,
				&s.CurrentPeriodStart,
				&s.CurrentPeriodEnd,
				&s.CreatedAt,
			); err != nil {
				http.Error(w, "Failed to read subscriptions", http.StatusInternalServerError)
				return
			}

			subs = append(subs, s)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if subs == nil {
			subs = []models.Subscription{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(subs)
	}
}
