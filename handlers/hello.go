package handlers

import "github.com/gofiber/fiber/v2"

// 예시 핸들러. 이 틀을 사용하여 핸들러를 작성하세요.
//
//	@summary		Hello API
//	@description	예시로 만든 API입니다. 해당 핸들러를 틀로 삼아 다른 핸들러를 작성하세요.
//	@tags			demo
//	@router			/api/v1/hello [get]
func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
