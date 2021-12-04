package auth

import (
	"cloudProject/pkg/validate"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func (req *SignUpRequest) Validate(ctx *fiber.Ctx) (string, int, error) {
	errStrName, err := validate.Struct(req)
	if err != nil {
		switch errStrName {
		case "UserName,Required":
			return "03", 400, errors.New("userNameRequired")
		}
	}
	return "", 200, nil
}
