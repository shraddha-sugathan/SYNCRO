package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
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


func HandleGetAdminReport(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `
			SELECT u.email, p.name, s.status, s.current_period_end 
			FROM subscriptions s 
			JOIN users u ON s.user_id = u.id 
			JOIN products p ON s.product_id = p.id
			ORDER BY s.created_at DESC`
		
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var reports []map[string]interface{}
		for rows.Next() {
			var email, product, status string
			var endDate time.Time
			rows.Scan(&email, &product, &status, &endDate)
			reports = append(reports, map[string]interface{}{
				"email": email, "product": product, "status": status, "end_date": endDate,
			})
		}
		json.NewEncoder(w).Encode(reports)
	}
}