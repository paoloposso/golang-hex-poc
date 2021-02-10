package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	h "github.com/paoloposso/poc/hex-ms/api"
	mr "github.com/paoloposso/poc/hex-ms/repository/mongo"
	rr "github.com/paoloposso/poc/hex-ms/repository/redis"
	"github.com/paoloposso/poc/hex-ms/shortener"
)

func main() {
	repo := chooseRepo()
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{code}", handler.Get)
	router.Post("/", handler.Post)

	errs := make(chan error, 2)

	go func() {
		fmt.Println("Listening on port ", httpPort())
		errs <- http.ListenAndServe(httpPort(), router)
	}()

	go func(){
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
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

func chooseRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisUrl := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository(redisUrl)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoUrl := os.Getenv("MONGO_URL")
		mongoDb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepository(mongoUrl, mongoDb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	default:
		redisUrl := "redis://localhost:6379"
		repo, err := rr.NewRedisRepository(redisUrl)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
}