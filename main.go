package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"syncro/auth"
	"syncro/handlers"
	"syncro/workers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func initSchema(db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role VARCHAR(50) DEFAULT 'customer',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		price_cents INT NOT NULL,
		billing_interval VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS subscriptions (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID REFERENCES users(id),
		product_id UUID REFERENCES products(id),
		status VARCHAR(50) DEFAULT 'active',
		current_period_start TIMESTAMP NOT NULL,
		current_period_end TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS invoices (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		subscription_id UUID REFERENCES subscriptions(id),
		user_id UUID REFERENCES users(id),
		amount_due_cents INT NOT NULL,
		status VARCHAR(50) DEFAULT 'pending',
		due_date TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	log.Println("Database schema initialized successfully.")
}

// --- Router Setup ---

func setupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve the Vue UI
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	// Public routes
	r.Post("/api/register", handlers.HandleRegister(db))
	r.Post("/api/login", handlers.HandleLogin(db))

	// Anyone needs to be able to see the catalog to subscribe!
	r.Get("/api/products", handlers.HandleGetProducts(db))

	// Protected Admin routes (Team Member 2 & 3)
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Use(auth.RequireRole("admin"))

		r.Post("/api/products", handlers.HandleCreateProduct(db))

		// Member 2's reporting endpoint
		r.Get("/api/revenue", handlers.HandleGetRevenue(db))
	})

	// Protected Customer routes (Team Member 2 & 3)
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Use(auth.RequireRole("customer"))

		// Member 3's subscription endpoints
		r.Post("/api/subscribe", handlers.HandleSubscribe(db))
		r.Get("/api/my-subscriptions", handlers.HandleGetSubscriptions(db))

		// Member 2's invoice endpoint
		r.Get("/api/invoices", handlers.HandleGetInvoices(db))
	})

	return r
}

// --- Database Connection ---

func initDB(connStr string) *sql.DB {
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

func main() {
	// Fallback to defaults if env vars aren't set (great for local hackathon dev)
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres" // your local docker password
	}

	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5433/syncro?sslmode=disable", dbPassword)
	db := initDB(connStr)
	defer db.Close()

	// Member 1's final touch: Build the tables automatically!
	initSchema(db)

	r := setupRouter(db)

	workers.StartBillingEngine(db)

	fmt.Println("Syncro server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
