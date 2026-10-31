package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"syncro/models"

	"github.com/google/uuid"
)

func HandleGetInvoices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, roleOK := r.Context().Value("role").(string)
		userID, userOK := r.Context().Value("user_id").(uuid.UUID)
		if !roleOK || !userOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if role != "customer" && role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		query := `SELECT id, subscription_id, user_id, amount_due_cents, status, due_date, created_at, updated_at FROM invoices`
		args := []interface{}{}
		if role == "customer" {
			query += " WHERE user_id = $1"
			args = append(args, userID)
		}
		query += " ORDER BY due_date DESC"

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var invoices []models.Invoice
		for rows.Next() {
			var inv models.Invoice
			if err := rows.Scan(&inv.ID, &inv.SubscriptionID, &inv.UserID, &inv.AmountDueCents, &inv.Status, &inv.DueDate, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			invoices = append(invoices, inv)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if invoices == nil {
			invoices = []models.Invoice{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}
