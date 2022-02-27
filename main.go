package main

import (
	"log"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/sonda2208/guardrails-challenge/api"
	"github.com/sonda2208/guardrails-challenge/app"
	"github.com/sonda2208/guardrails-challenge/model"
)

// @title Server API
func main() {
	configFilePath := os.Getenv("CONFIG")
	if len(configFilePath) == 0 {
		configFilePath = "./config.json"
	}

	conf, err := model.ConfigFromFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	a, err := app.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	_, err = api.New(a, e)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
