package main

import (
	"log"
	"net/http"
	"usr_mngmnt/pkg/model"
	"usr_mngmnt/pkg/rest"
	"usr_mngmnt/pkg/usecase"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init config.
	type config struct {
		MysqlDsn string `env:"MYSQL_DSN"`
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err.Error())
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Rest controller.
	userRepo, err := model.NewSqlcRepo(cfg.MysqlDsn)
	if err != nil {
		log.Fatalf(err.Error())
	}

	userCntrl := rest.NewUserCntrl(usecase.NewUserUsecases(userRepo))

	// Routes
	e.GET("/", hello)
	e.GET("/users/:offset/:limit", userCntrl.ViewAll)
	e.POST("/users", userCntrl.Create)
	e.PUT("/users", userCntrl.Update, middleware.BasicAuth(userCntrl.BasicAuthMiddleware))
	e.GET("/user/:uuid", userCntrl.View)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
