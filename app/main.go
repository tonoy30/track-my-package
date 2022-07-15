package main

import (
	"log"

	"track-my-package/app/package/client"
	delivery "track-my-package/app/package/delivery/http"
	"track-my-package/app/package/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.File("/", "app/static/index.html")
	rmc, err := client.NewRabbitMqClient("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicln(err)
	}
	defer rmc.Close()
	pu := usecase.NewPackageUseCase(rmc)
	delivery.NewPackageHandler(e, pu)
	e.Logger.Fatal(e.Start(":1323"))
}
