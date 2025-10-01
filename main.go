package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dmitriy-zverev/blog-api/internal/db"
	"github.com/dmitriy-zverev/blog-api/internal/handlers"
	"github.com/dmitriy-zverev/blog-api/internal/models"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	config, err := loadEnv()
	if err != nil {
		log.Fatal("error loading config:", err)
	}

	dbQueries, db, err := initDatabase(config.DBUrl)
	if err != nil {
		log.Fatal("error initializing database: ", err)
	}
	defer db.Close()

	mux := setupRoutes(createApiConfig(config, dbQueries))

	if err := startServer(mux, config.Port); err != nil {
		log.Fatal("server error:", err)
	}
}

func loadEnv() (*models.Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error: .env file not found, using system variables")
	}

	build := os.Getenv("BUILD")
	if build == "" {
		return nil, errors.New("BUILD string is requiered")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("PORT string is requiered")
	}

	var dbUrl string

	if build == "dev" {
		dbUrl = os.Getenv("DB_URL")
		if dbUrl == "" {
			return nil, errors.New("DB_URL string is requiered")
		}
	} else {
		dbHost := os.Getenv("DB_HOST")
		if dbHost == "" {
			return nil, errors.New("DB_HOST string is requiered")
		}
		dbPort := os.Getenv("DB_PORT")
		if dbPort == "" {
			return nil, errors.New("DB_PORT string is requiered")
		}
		dbUser := os.Getenv("DB_USER")
		if dbUser == "" {
			return nil, errors.New("DB_USER string is requiered")
		}
		dbPassword := os.Getenv("DB_PASSWORD")
		if dbPassword == "" {
			return nil, errors.New("DB_PASSWORD string is requiered")
		}
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			return nil, errors.New("DB_NAME string is requiered")
		}

		dbUrl = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		)
	}

	return &models.Config{
		Build: build,
		Port:  port,
		DBUrl: dbUrl,
	}, nil
}

func setupRoutes(cfg *handlers.ApiConfig) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET "+API_VERSION+BASE_ROUTE, cfg.BaseHandler)
	mux.HandleFunc("POST "+API_VERSION+POSTS_ROUTE, cfg.PostsPostHandler)
	mux.HandleFunc("PUT "+API_VERSION+POST_ROUTE, cfg.PostsPutHandler)
	mux.HandleFunc("DELETE "+API_VERSION+POST_ROUTE, cfg.PostsDeleteHandler)
	mux.HandleFunc("GET "+API_VERSION+POST_ROUTE, cfg.PostsGetOneHandler)
	mux.HandleFunc("GET "+API_VERSION+POSTS_ROUTE, cfg.PostsGetManyHandler)

	return mux
}

func createApiConfig(cfg *models.Config, dbQueries *db.Queries) *handlers.ApiConfig {
	return &handlers.ApiConfig{
		Build: cfg.Build,
		Port:  cfg.Port,
		DB:    dbQueries,
	}
}

func startServer(mux *http.ServeMux, port string) error {
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Println("Starting server...")
	log.Printf("Running on http://localhost:%s\n", port)

	return server.ListenAndServe()
}

func initDatabase(dbUrl string) (*db.Queries, *sql.DB, error) {
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbQueries := db.New(database)
	return dbQueries, database, nil
}
