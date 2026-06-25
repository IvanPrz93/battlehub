package main

import (
	"database/sql"
    "log"
    "net/http"
	"os"

	"github.com/IvanPrz93/battlehub/internal/database"
	"github.com/joho/godotenv"

	 _ "github.com/lib/pq"
)

type apiConfig struct {
	db			*database.Queries
	dbConn		*sql.DB
	jwtSecret	string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load() 
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		db:         dbQueries,
		dbConn:		dbConn,
		jwtSecret:  jwtSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	
	mux.HandleFunc("GET /api/users/me", apiCfg.handlerMyProfile)
	mux.HandleFunc("GET /api/users/{username}", apiCfg.handlerGetUser)

	mux.HandleFunc("GET /api/users/me/matches", apiCfg.handlerMyMatchHistory)
	mux.HandleFunc("GET /api/users/{username}/matches", apiCfg.handlerGetMatchHistory)

	mux.HandleFunc("POST /api/matches", apiCfg.createNewMatch)
	mux.HandleFunc("POST /api/matches/{id}/result/", apiCfg.endMatch)

	mux.HandleFunc("GET /api/leaderboard", apiCfg.handlerLeaderboard)


	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
