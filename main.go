package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"regexp"

	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// --- Models ---

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
	UpdatedAt      time.Time `json:"updated_at"`
}

// Use an environment variable for this in production
var jwtSecret = []byte("syncro-hackathon-super-secret-key")

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
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

func validatePassword(password string) bool {
	if len(password) <= 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasUpper && hasLower && hasSpecial
}

// --- JWT Utilities ---

func GenerateToken(userID uuid.UUID, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// --- Middleware ---

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store claims in request context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("role").(string)
			if !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			for _, role := range allowedRoles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

// --- Router Setup ---

func setupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Post("/api/register", handleRegister(db))
	r.Post("/api/login", handleLogin(db))
	r.Get("/api/products", handleGetProducts(db))

	//r.Post("/api/login", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("Login handler goes here"))
	//})

	// Protected Admin routes
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Use(RequireRole("admin"))

		r.Post("/api/products", handleCreateProduct(db))
		r.Get("/api/admin/revenue", handleGetRevenue(db))
	})

	// Authenticated invoice routes
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Get("/api/invoices", handleGetInvoices(db))
	})

	// Protected Customer routes
	r.Group(func(r chi.Router) {
		r.Use(RequireAuth)
		r.Use(RequireRole("customer"))

		r.Get("/api/my-subscriptions", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("user_id").(uuid.UUID)
			w.Write([]byte(fmt.Sprintf("Fetch subscriptions for %s", userID)))
		})
	})

	return r
}

// --- Database Connection ---

func initDB(password string) *sql.DB {
	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5432/syncro?sslmode=disable", password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the syncro database.")
	return db
}

func handleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if !validatePassword(req.Password) {
			http.Error(w, "Password must be >8 chars, with uppercase, lowercase, and a special character", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var newID uuid.UUID
		err = db.QueryRow(
			"INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id",
			req.Email, string(hashedPassword), "customer",
		).Scan(&newID)

		if err != nil {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Account created successfully", "id": newID.String()})
	}
}

func handleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user User
		err := db.QueryRow("SELECT id, password_hash, role FROM users WHERE email = $1", req.Email).
			Scan(&user.ID, &user.PasswordHash, &user.Role)

		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := GenerateToken(user.ID, user.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

func handleCreateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		if req.Name == "" {
			http.Error(w, "Product name is required", http.StatusBadRequest)
			return
		}
		if req.PriceCents < 0 {
			http.Error(w, "Price cents must be non-negative", http.StatusBadRequest)
			return
		}
		if req.BillingInterval != "monthly" && req.BillingInterval != "yearly" {
			http.Error(w, "Billing interval must be 'monthly' or 'yearly'", http.StatusBadRequest)
			return
		}

		var newID uuid.UUID
		err := db.QueryRow(
			"INSERT INTO products (name, price_cents, billing_interval) VALUES ($1, $2, $3) RETURNING id",
			req.Name, req.PriceCents, req.BillingInterval,
		).Scan(&newID)

		if err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Product created", "id": newID.String()})
	}
}

func handleGetProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, price_cents, billing_interval, created_at FROM products ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		products := make([]Product, 0)
		for rows.Next() {
			var product Product
			if err := rows.Scan(&product.ID, &product.Name, &product.PriceCents, &product.BillingInterval, &product.CreatedAt); err != nil {
				http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(products); err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		}
	}
}

func handleGetInvoices(db *sql.DB) http.HandlerFunc {
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

		query := `SELECT id, subscription_id, user_id, amount_due_cents, status, due_date, created_at, updated_at
			FROM invoices`
		args := []interface{}{}
		if role == "customer" {
			query += " WHERE user_id = $1"
			args = append(args, userID)
		}
		query += " ORDER BY created_at DESC"

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Failed to fetch invoices", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		invoices := make([]Invoice, 0)
		for rows.Next() {
			var invoice Invoice
			if err := rows.Scan(&invoice.ID, &invoice.SubscriptionID, &invoice.UserID, &invoice.AmountDueCents, &invoice.Status, &invoice.DueDate, &invoice.CreatedAt, &invoice.UpdatedAt); err != nil {
				http.Error(w, "Failed to fetch invoices", http.StatusInternalServerError)
				return
			}
			invoices = append(invoices, invoice)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "Failed to fetch invoices", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(invoices); err != nil {
			http.Error(w, "Failed to fetch invoices", http.StatusInternalServerError)
		}
	}
}

func handleGetRevenue(db *sql.DB) http.HandlerFunc {
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
			JOIN products p ON p.id = s.product_id
			WHERE s.status = 'active'`).Scan(&mrrCents)
		if err != nil {
			http.Error(w, "Failed to calculate revenue", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]int64{"mrr_cents": mrrCents}); err != nil {
			http.Error(w, "Failed to calculate revenue", http.StatusInternalServerError)
		}
	}
}

func main() {
	db := initDB("postgres")
	defer db.Close()

	r := setupRouter(db)

	fmt.Println("Syncro server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
