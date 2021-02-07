package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	mr "github.com/paoloposso/poc/hex-ms/repository/mongo"
	rr "github.com/paoloposso/poc/hex-ms/repository/redis"
	"github.com/paoloposso/poc/hex-ms/shortener"
)

func main() {

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
	}
	return nil
}