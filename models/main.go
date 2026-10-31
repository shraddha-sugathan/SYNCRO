package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Omitted from JSON responses
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Product struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	PriceCents      int       `json:"price_cents"`
	BillingInterval string    `json:"billing_interval"`
	CreatedAt       time.Time `json:"created_at"`
}

type Subscription struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	ProductID          uuid.UUID `json:"product_id"`
	Status             string    `json:"status"`
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
	CreatedAt          time.Time `json:"created_at"`
}

type Invoice struct {
	ID             uuid.UUID `json:"id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	UserID         uuid.UUID `json:"user_id"`
	AmountDueCents int       `json:"amount_due_cents"`
	Status         string    `json:"status"`
	DueDate        time.Time `json:"due_date"`
	CreatedAt      time.Time `json:"created_at"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProductRequest struct {
	Name            string `json:"name"`
	PriceCents      int    `json:"price_cents"`
	BillingInterval string `json:"billing_interval"`
}

type SubscribeRequest struct {
	ProductID uuid.UUID `json:"product_id"`
}