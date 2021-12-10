package main

import (
	"cloudProject/controllers/mainController"
	"cloudProject/pkg/createService"
)

func main() {
	app := createService.New("dataAnalytic")
	mainController.StartPrometheus(routes)
	setupRoute(app)
	panic(app.Listen(":7575"))
}
