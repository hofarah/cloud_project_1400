package main

import (
	"cloudProject/controllers/dataAnalyticController"
	"github.com/gofiber/fiber/v2"
)

var routes = map[string]string{
	"sellsCompare2Game": "/chart/sellCompare2Game",
	"sumSellAnnually":   "/chart/sumSellAnnually",
	"producerSells":     "/chart/producerSells",
	"genreSells":        "/chart/genreSells",
}

func setupRoute(app *fiber.App) {
	app.Get(routes["genreSells"], dataAnalyticController.GenreSells)
	app.Get(routes["producerSells"], dataAnalyticController.ProducerSells)
	app.Get(routes["sellsCompare2Game"], dataAnalyticController.SellsCompare2Game)
	app.Get(routes["sumSellAnnually"], dataAnalyticController.SumSellAnnually)
}
