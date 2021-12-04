package global

import (
	. "cloudProject/controllers/global"
	"cloudProject/middlewares"
	"github.com/gofiber/fiber/v2"
)

var routes = map[string]string{
	"getByRank": "/global/get/byRank",
	"getByName": "/global/get/byName",
}

func setUpRoutes(app *fiber.App) {
	app.Use(middlewares.Auth)
	app.Get(routes["getByRank"], GetByRank)
	app.Get(routes["getByName"], GetByName)
}
