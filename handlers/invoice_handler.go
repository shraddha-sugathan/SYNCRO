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
		userID := r.Context().Value("user_id").(uuid.UUID)

		rows, err := db.Query(
			`SELECT id, subscription_id, user_id, amount_due_cents, status, due_date, created_at 
			 FROM invoices WHERE user_id = $1 ORDER BY due_date DESC`, userID,
		)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var invoices []models.Invoice
		for rows.Next() {
			var inv models.Invoice
			if err := rows.Scan(&inv.ID, &inv.SubscriptionID, &inv.UserID, &inv.AmountDueCents, &inv.Status, &inv.DueDate, &inv.CreatedAt); err != nil {
				continue
			}
			invoices = append(invoices, inv)
		}

		if invoices == nil {
			invoices = []models.Invoice{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoices)
	}
}