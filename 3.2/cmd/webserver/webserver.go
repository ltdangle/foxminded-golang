package main

import (
	"fmt"
	"jwt/pkg/cache"
	"jwt/pkg/model"
	"jwt/pkg/rest"
	"jwt/pkg/usecase"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Init config.
	type config struct {
		MysqlDsn        string `env:"MYSQL_DSN" validate:"required"`
		JwtKey          string `env:"JWT_KEY" validate:"required"`
		RedisAddr       string `env:"REDIS_ADDR" validate:"required"`
		RedisPass       string `env:"REDIS_PASS"`
		RedisDb         int    `env:"REDIS_DB" `
		CacheTimeoutSec int    `env:"CACHE_TIMEOUT_SEC"`
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	validator := validator.New()
	err = validator.Struct(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(cfg)

	db, err := sqlx.Open("mysql", cfg.MysqlDsn)
	if err != nil {
		log.Fatal("failed to connect database")
	}

	userRepo := model.NewWriteUserRepo(db)
	uCase := usecase.NewUserUsecase(userRepo)

	voteRepo := model.NewWriteVoteRepo(db)
	voteCase := usecase.NewVoteUsecase(voteRepo)

	logger := logrus.New()
	rspondr := rest.NewResponder()
	jwt := rest.NewJwtService([]byte(cfg.JwtKey))
	cache := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDb)
	userCntrl := rest.NewUserController(uCase, voteCase, validator, rspondr, logger, jwt, cache, cfg.CacheTimeoutSec)
	auth := rest.NewAuth(jwt, logger)

	router := rest.SetupRouter(auth, userCntrl)

	log.Info("Starting server on localhost:8080")
	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
