package global

import (
	"cloudProject/controllers/mainController"
	"cloudProject/pkg/createService"
)

func main() {
	app := createService.New("global")
	mainController.StartPrometheus(routes)
	setUpRoutes(app)
	panic(app.Listen(":7575"))
}
