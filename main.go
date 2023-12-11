package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/danny-ell77/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}


func main()  {
	godotenv.Load() 

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB url not found")
	}

	if port == "" {
		log.Fatal("Value for PORT not found")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("could not connect to database:", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
  	}))

	v1Router := chi.NewRouter() 

	v1Router.Get("/health", readinessHandler)
	v1Router.Get("/error", errorHandler)
	v1Router.Post("/users", apiCfg.createUser)
	v1Router.Get("/users", apiCfg.authMiddleware(apiCfg.getUser))
	v1Router.Post("/feeds", apiCfg.authMiddleware(apiCfg.createFeed))
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: 	router,
		Addr: 		":" + port,
	}

	log.Printf("Server starting on %v", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Port:", port)
	// fmt.Println("Hello World")
}