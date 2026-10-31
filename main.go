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
	_ "github.com/lib/pq"
)

// --- Router Setup ---

func setupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Simple health-check for the browser
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Syncro Billing Engine API is running!"))
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