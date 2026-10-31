package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func HandleGetRevenue(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var totalRevenueCents int
		// Sum up the price of all currently active subscriptions
		err := db.QueryRow(`
			SELECT COALESCE(SUM(p.price_cents), 0) 
			FROM subscriptions s 
			JOIN products p ON s.product_id = p.id 
			WHERE s.status = 'active'
		`).Scan(&totalRevenueCents)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{
			"active_mrr_cents": totalRevenueCents,
		})
	}
}