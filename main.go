package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danny-ell77/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

type Middleware func(h http.HandlerFunc) http.HandlerFunc

func Adapt(h http.HandlerFunc, middlewares []Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
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

	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	middlewares := []Middleware{apiCfg.newAuthMiddleware, loggingMiddleware}

	go startScraping(db, 10, time.Minute)
	
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
	v1Router.Get("/feeds", apiCfg.getAllFeeds)
	v1Router.Post("/feed_follows", apiCfg.authMiddleware(apiCfg.followFeed))
	v1Router.Delete("/feed_follows/{feedFollowId}", apiCfg.authMiddleware(apiCfg.unfollowFeed))
	v1Router.Get("/feed_follows", apiCfg.authMiddleware(apiCfg.getFeedFollow))
	v1Router.Get("/posts", Adapt(apiCfg.getPostsForUser, middlewares))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: 	router,
		Addr: 		":" + port,
	}

	log.Printf("Server starting on %v", port)
	log.Fatal(server.ListenAndServe())
	// fmt.Println("Port:", port)
	// fmt.Println("Hello World")
}