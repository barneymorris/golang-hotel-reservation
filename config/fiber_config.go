package config

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var Config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        code := fiber.StatusInternalServerError

		var e *fiber.Error
        if errors.As(err, &e) {
            code = e.Code
        }

		ctx.Status(code)

		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}
