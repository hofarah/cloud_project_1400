package global

import (
	"cloudProject/middlewares"
	"github.com/gofiber/fiber/v2"
)

var routes = map[string]string{}

func setUpRoutes(app *fiber.App) {
	app.Use(middlewares.Auth)
	global := app.Group("/global")
	global.Get("")
}
