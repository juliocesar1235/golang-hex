package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	apiHandler "github.com/juliocesar1235/golang-hex/api"
	mongoRepository "github.com/juliocesar1235/golang-hex/repository/mongo"
	redisRepository "github.com/juliocesar1235/golang-hex/repository/redis"

	"github.com/juliocesar1235/golang-hex/shortener"
)

// https://www.google.com -> 98sj1-293
// http://localhost:8000/98sj1-293 -> https://www.google.com

func main() {
	repository := getRepository()
	service := shortener.NewRedirectServiceInstance
	handler := apiHandler.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != ""{
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}


func getRepository() shortener.RedirectRepository {
	dbUrl := os.Getenv("URL_DB")
	switch dbUrl {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redisRepository.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongoRepository.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}