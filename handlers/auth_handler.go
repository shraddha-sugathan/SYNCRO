package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"

	"syncro/auth"
	"syncro/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

func GoogleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Replace with your actual Client ID
		payload, err := idtoken.Validate(context.Background(), req.Token, "YOUR_GOOGLE_CLIENT_ID")
		if err != nil {
			http.Error(w, "Invalid Google token", http.StatusUnauthorized)
			return
		}

		email := payload.Claims["email"].(string)
		var user models.User

		err = db.QueryRow("SELECT id, email, role FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.Role)

		if err == sql.ErrNoRows {
			user = models.User{
				ID:           uuid.New(),
				Email:        email,
				Role:         "customer",
				PasswordHash: "",
			}
			_, dbErr := db.Exec("INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)", user.ID, user.Email, user.PasswordHash, user.Role)
			if dbErr != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		}

		token, err := auth.GenerateToken(user.ID, user.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
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

func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AuthRequest
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

func HandleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user models.User
		err := db.QueryRow("SELECT id, password_hash, role FROM users WHERE email = $1", req.Email).
			Scan(&user.ID, &user.PasswordHash, &user.Role)

		if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}