package main

import (
	. "cloudProject/controllers/authentication"
	"cloudProject/middlewares"
	"github.com/gofiber/fiber/v2"
)

var routes = map[string]string{
	"check":  "/authentication/check",
	"signup": "/authentication/signup",
}

func setUpRoutes(app *fiber.App) {
	app.Use(middlewares.CheckRequiredHeaders)
	app.Post(routes["signup"], SignUp)
	app.Get(routes["check"], Check)
}
