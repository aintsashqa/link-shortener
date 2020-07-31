package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/aintsashqa/link-shortener/src/api"
	"github.com/aintsashqa/link-shortener/src/repository/mongodb"
	"github.com/aintsashqa/link-shortener/src/repository/redis"
	"github.com/aintsashqa/link-shortener/src/shortener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	repository := createRepository()
	service := shortener.NewRedirectService(repository)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		port := httpPort()
		fmt.Printf("Listening on port %s\n", port)
		errs <- http.ListenAndServe(port, router)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func createRepository() shortener.RedirectRepositoryInterface {
	switch os.Getenv("DB_DRIVER") {
	case "redis":
		rURL := os.Getenv("REDIS_URL")
		repository, err := redis.NewRedisRepository(rURL)
		if err != nil {
			log.Fatal(err)
		}
		return repository
	case "mongo":
		mURL := os.Getenv("MONGO_URL")
		mdb := os.Getenv("MONGO_DATABASE")
		timeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repository, err := mongodb.NewMongoRepository(mURL, mdb, timeout)
		if err != nil {
			log.Fatal(err)
		}
		return repository
	}
	return nil
}
