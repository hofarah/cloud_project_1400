package common

type Response struct {
	Status string      `json:"status"`
	Error  Error       `json:"error"`
	Data   interface{} `json:"data"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
