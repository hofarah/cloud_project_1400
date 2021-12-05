package main

import (
	. "cloudProject/controllers/global"
	"cloudProject/middlewares"
	"github.com/gofiber/fiber/v2"
)

var routes = map[string]string{
	"getByRank":                  "/global/get/byRank",
	"getByName":                  "/global/get/byName",
	"getBestOnPlatform":          "/global/get/bestOnPlatform",
	"getBestOnYear":              "/global/get/bestOnYear",
	"getBestOnGenre":             "/global/get/bestOnGenre",
	"getBest5YearPlatform":       "/global/get/best5SellsByYearAndPlatform",
	"europeMoreThanNorthAmerica": "/global/get/europeMoreThanNorthAmerica",
}

func setUpRoutes(app *fiber.App) {
	app.Use(middlewares.Auth())
	app.Get(routes["getByRank"], GetByRank)
	app.Get(routes["getByName"], GetByName)
	app.Get(routes["getBestOnPlatform"], GetBestOnPlatform)
	app.Get(routes["getBestOnYear"], GetBestOnYear)
	app.Get(routes["getBestOnGenre"], GetBestOnGenre)
	app.Get(routes["getBest5YearPlatform"], GetBest5SellsByYearAndPlatform)
	app.Get(routes["europeMoreThanNorthAmerica"], GetEuropeMoreThanNorthAmerica)
}
