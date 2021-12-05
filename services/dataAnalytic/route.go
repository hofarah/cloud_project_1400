package main

import (
	"cloudProject/controllers/dataAnalytic"
	"cloudProject/middlewares"
	"github.com/gofiber/fiber/v2"
)

var routes = map[string]string{
	"sellsCompare2Game": "/chart/sellCompare2Game",
	"sumSellAnnually":   "/chart/sumSellAnnually",
	"producerSells":     "/chart/producerSells",
	"genreSells":        "/chart/genreSells",
}

func setupRoute(app *fiber.App) {
	app.Use(middlewares.Auth())
	app.Get(routes["genreSells"], dataAnalytic.GenreSells)
	app.Get(routes["producerSells"], dataAnalytic.ProducerSells)
	app.Get(routes["sellsCompare2Game"], dataAnalytic.SellsCompare2Game)
	app.Get(routes["sumSellAnnually"], dataAnalytic.SumSellAnnually)
}
