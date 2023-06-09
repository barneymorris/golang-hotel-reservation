package middleware

import (
	"github.com/betelgeusexru/golang-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fiber.NewError(401, "not authorized")
	}

	if !user.IsAdmin {
		return fiber.NewError(401, "not authorized")
	}
	
	return c.Next()
}