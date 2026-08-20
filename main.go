package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"web123/config"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreateUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func main() {
	// Загрузка конфига из .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	//Формирование строки подключения для pgx
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSslMode,
	)

	// Создаю контекст
	ctx := context.Background()

	//Подключаюсь к PostgreSQL через pgx
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	defer conn.Close(ctx)

	fmt.Printf("Connected to PostgreSQL: %s\n", cfg.DBName)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsersHandler(conn)(w, r)
		case http.MethodPost:
			createUserHandler(conn)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Порт из конфига
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port :%s", port)
	// main не выполнится пока работает сервер
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func getUsersHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "application/json")

		rows, err := conn.Query(ctx, "SELECT id, email, role FROM users")
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			err = rows.Scan(&u.ID, &u.Email, &u.Role)
			if err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "json encode error", http.StatusInternalServerError)
			return
		}
	}
}

func createUserHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		// Читаем JSON из тела запроса
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		// Поля не должны быть пустыми
		req.Email = strings.TrimSpace(req.Email)
		req.Role = strings.TrimSpace(req.Role)

		if req.Email == "" || req.Role == "" {
			http.Error(w, "email and role are required", http.StatusBadRequest)
			return
		}

		// Защита от дублей по UNIQUE на email
		var id int
		err := conn.QueryRow(
			ctx,
			`INSERT INTO users (email, role)
			 VALUES ($1, $2)
			 ON CONFLICT (email) DO NOTHING
			 RETURNING id`,
			req.Email,
			req.Role,
		).Scan(&id)

		if err != nil {

			if err.Error() == "no rows in result set" {
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "user already exists",
					"email": req.Email,
				})
				return
			}

			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		createdUser := User{
			ID:    id,
			Email: req.Email,
			Role:  req.Role,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdUser)
	}
}
