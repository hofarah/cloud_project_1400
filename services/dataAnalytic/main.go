package main

import (
	"cloudProject/controllers/mainController"
	"cloudProject/pkg/createService"
)

func main() {
	app := createService.New("authentication")
	mainController.StartPrometheus(routes)
	setupRoute(app)
	panic(app.Listen(":7575"))
}
