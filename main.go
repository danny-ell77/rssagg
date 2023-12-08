package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main()  {
	godotenv.Load() 

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Value for PORT not found")
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

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: 	router,
		Addr: 		":" + port,
	}

	log.Printf("Server starting on %v", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Port:", port)
	// fmt.Println("Hello World")
}