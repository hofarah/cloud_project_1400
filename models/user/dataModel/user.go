package dataModel

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"username"`
	Secret   string `json:"secret"`
}
