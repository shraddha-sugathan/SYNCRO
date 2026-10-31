package workers

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

// StartBillingEngine starts the background billing worker.
func StartBillingEngine(db *sql.DB) {
	log.Println("Starting background billing engine...")

	// Process any overdue subscriptions immediately on startup.
	go processRenewals(db)

	// Then check for renewals every minute.
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			processRenewals(db)
		}
	}()
}

func processRenewals(db *sql.DB) {
	now := time.Now()

	rows, err := db.Query(`
		SELECT id
		FROM subscriptions
		WHERE status = 'active'
		  AND current_period_end <= $1
	`, now)

	if err != nil {
		log.Printf("Billing engine error fetching renewals: %v", err)
		return
	}

	var subscriptionIDs []uuid.UUID

	for rows.Next() {
		var subID uuid.UUID

		if err := rows.Scan(&subID); err != nil {
			log.Printf("Error scanning subscription ID: %v", err)
			continue
		}

		subscriptionIDs = append(subscriptionIDs, subID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error reading subscriptions: %v", err)
	}

	rows.Close()

	// Process each subscription separately.
	for _, subID := range subscriptionIDs {
		processSingleRenewal(db, subID)
	}
}

func processSingleRenewal(db *sql.DB, subID uuid.UUID) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to start transaction for %s: %v", subID, err)
		return
	}

	defer tx.Rollback()

	var userID uuid.UUID
	var priceCents int
	var interval string
	var currentEnd time.Time
	var status string

	err = tx.QueryRow(`
		SELECT s.user_id,
		       p.price_cents,
		       p.billing_interval,
		       s.current_period_end,
		       s.status
		FROM subscriptions s
		JOIN products p ON s.product_id = p.id
		WHERE s.id = $1
		FOR UPDATE OF s
	`, subID).Scan(
		&userID,
		&priceCents,
		&interval,
		&currentEnd,
		&status,
	)

	if err == sql.ErrNoRows {
		return
	}

	if err != nil {
		log.Printf("Failed to lock subscription %s: %v", subID, err)
		return
	}

	now := time.Now()

	if status != "active" || currentEnd.After(now) {
		return
	}

	overdueDuration := now.Sub(currentEnd)
	overdueDays := int(overdueDuration.Hours() / 24)

	newStart := now
	var newEnd time.Time

	switch interval {
	case "monthly":
		newEnd = newStart.AddDate(0, 1, 0)

	case "yearly":
		newEnd = newStart.AddDate(1, 0, 0)

	default:
		log.Printf(
			"Unsupported billing interval '%s' for subscription %s",
			interval,
			subID,
		)
		return
	}

	dueDate := now.AddDate(0, 0, 7)

	_, err = tx.Exec(`
		INSERT INTO invoices
			(subscription_id, user_id, amount_due_cents, status, due_date)
		VALUES ($1, $2, $3, 'pending', $4)
	`, subID, userID, priceCents, dueDate)

	if err != nil {
		log.Printf("Failed to generate invoice for %s: %v", subID, err)
		return
	}

	_, err = tx.Exec(`
		UPDATE subscriptions
		SET current_period_start = $1,
		    current_period_end = $2
		WHERE id = $3
	`, newStart, newEnd, subID)

	if err != nil {
		log.Printf("Failed to update subscription %s: %v", subID, err)
		return
	}

	// Commit both operations together.
	if err := tx.Commit(); err != nil {
		log.Printf("Failed to commit renewal for %s: %v", subID, err)
		return
	}

	log.Printf(
		"SUCCESS: Subscription %s renewed. Overdue by %d days. Invoice due: %s",
		subID,
		overdueDays,
		dueDate.Format(time.RFC3339),
	)
}
