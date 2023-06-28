package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // the underscore is to tell the compiler that we are using the package but not directly
	"github.com/sandepten/go-rss-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load() // we can provide the env filename as a parameter but the default is .env
	port := os.Getenv("PORT")
	if port == "" {
		// log.Fatal will immediately exit the program with code 1 and the message
		log.Fatal("PORT not found in the environment")
	}

	//? database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not found in the environment")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	apiConfig := apiConfig{ // now we can pass this apiConfig to our handlers, so that they can access the database
		DB: database.New(db),
	}

	//? router
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)

	// user routes
	v1Router.Post("/createUser", apiConfig.handlerCreateUser)
	v1Router.Get("/getUserByEmail", apiConfig.handlerGetUserByEmail)
	v1Router.Get("/getUserByAPIKey", apiConfig.middlewareAuth(apiConfig.handlerGetUserByAPIKey))
	v1Router.Get("/users", apiConfig.handlerGetAllUsers)

	// feed routes
	v1Router.Post("/addFeed", apiConfig.middlewareAuth(apiConfig.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfig.handlerGetAllFeeds)

	// feed follow routes
	v1Router.Post("/feed_follow", apiConfig.middlewareAuth(apiConfig.handlerCreateFeedFollow))
	v1Router.Get("/feed_follow", apiConfig.middlewareAuth(apiConfig.handlerGetFeedFollows))
	v1Router.Delete("/feed_follow/{feedFollowId}", apiConfig.middlewareAuth(apiConfig.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Println("Server is running on port", port)

	// listen and serve will block the program until it is terminated
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
