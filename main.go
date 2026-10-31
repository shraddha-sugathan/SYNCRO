package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"syncro/auth"
	"syncro/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// --- Router Setup ---

func setupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public routes
	r.Post("/api/register", handlers.HandleRegister(db))
	r.Post("/api/login", handlers.HandleLogin(db))

	// Protected Admin routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Use(auth.RequireRole("admin"))
		
		r.Post("/api/products", handlers.HandleCreateProduct(db))
	})

	// Protected Customer routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Use(auth.RequireRole("customer"))

		// New subscription endpoint
		r.Post("/api/subscribe", handlers.HandleSubscribe(db))
		
		// Existing placeholder for fetching a user's subscriptions
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

func main() {
	db := initDB("postgres")
	defer db.Close()

	r := setupRouter(db)
	
	fmt.Println("Syncro server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}