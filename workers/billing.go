package workers

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// StartBillingEngine boots up the background goroutine
func StartBillingEngine(db *sql.DB) {
	log.Println("Starting background billing engine...")
	ticker := time.NewTicker(1 * time.Minute) // Sweeps every 1 minute
	
	go func() {
		for {
			<-ticker.C
			processRenewals(db)
		}
	}()
}

func processRenewals(db *sql.DB) {
	now := time.Now()

	// Find active subscriptions that need to be renewed
	rows, err := db.Query(`
		SELECT s.id, s.user_id, p.price_cents, p.billing_interval, s.current_period_end
		FROM subscriptions s
		JOIN products p ON s.product_id = p.id
		WHERE s.status = 'active' AND s.current_period_end <= $1
	`, now)
	
	if err != nil {
		log.Printf("Billing engine error fetching renewals: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var subID, userID uuid.UUID
		var priceCents int
		var interval string
		var currentEnd time.Time

		if err := rows.Scan(&subID, &userID, &priceCents, &interval, &currentEnd); err != nil {
			log.Printf("Error scanning subscription: %v\n", err)
			continue
		}

		// Calculate the next billing cycle
		newStart := currentEnd
		var newEnd time.Time
		if interval == "monthly" {
			newEnd = newStart.AddDate(0, 1, 0)
		} else {
			newEnd = newStart.AddDate(1, 0, 0)
		}

		// Use an atomic transaction so we don't accidentally update a subscription 
		// without generating the invoice if the server crashes halfway through.
		tx, err := db.Begin()
		if err != nil {
			continue
		}

		// 1. Generate the pending invoice (Due in 7 days)
		dueDate := now.AddDate(0, 0, 7)
		_, err = tx.Exec(`
			INSERT INTO invoices (subscription_id, user_id, amount_due_cents, status, due_date)
			VALUES ($1, $2, $3, 'pending', $4)
		`, subID, userID, priceCents, dueDate)

		if err != nil {
			tx.Rollback()
			continue
		}

		// 2. Advance the subscription dates
		_, err = tx.Exec(`
			UPDATE subscriptions 
			SET current_period_start = $1, current_period_end = $2
			WHERE id = $3
		`, newStart, newEnd, subID)

		if err != nil {
			tx.Rollback()
			continue
		}

		tx.Commit()
		log.Printf("SUCCESS: Renewed subscription %s and generated invoice.\n", subID.String())
	}
}