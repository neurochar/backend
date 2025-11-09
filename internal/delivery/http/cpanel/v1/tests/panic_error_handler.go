package tests

import (
	"github.com/gofiber/fiber/v2"
)

func (ctrl *Controller) PanicErrorHandler(c *fiber.Ctx) error {
	panic("this is test panic")
}
