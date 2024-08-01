package handler

import "github.com/gofiber/fiber/v2"

// HandlerContext is the struct that holds the context of the handler
type HandlerContext struct {
	RunMode string
	Fiber   *fiber.App
	Topic   interface{}
}
