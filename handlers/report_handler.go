package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func HandleGetRevenue(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var mrrCents int64
		err := db.QueryRow(`
			SELECT COALESCE(SUM(
				CASE p.billing_interval
					WHEN 'monthly' THEN p.price_cents
					WHEN 'yearly' THEN p.price_cents / 12
					ELSE 0
				END
			), 0)
			FROM subscriptions s
			JOIN products p ON s.product_id = p.id
			WHERE s.status = 'active'
		`).Scan(&mrrCents)

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int64{
			"mrr_cents": mrrCents,
		})
	}
}
