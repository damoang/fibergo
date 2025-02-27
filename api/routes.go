package api

import (
	"github.com/gofiber/fiber/v2"
)

func setupRoutes() {
	// 정적 파일 제공
	app.Static("/static", "./static")

	// API 라우트
	apiGroup := app.Group("/api")
	apiGroup.Get("/:type", HandleBoardAPI)
	apiGroup.Get("/:type/:id/comments", HandleCommentsAPI)

	// 웹 페이지 라우트
	app.Get("/:type", HandleBoardSSR)
	app.Get("/:type/:id", HandleBoardSSR)

	// 404 에러 핸들러
	app.Use(func(c *fiber.Ctx) error {
		if c.Accepts("json") {
			return c.Status(404).JSON(fiber.Map{
				"error": "요청하신 페이지를 찾을 수 없습니다",
			})
		}
		return c.Status(404).SendFile("templates/404.html")
	})
} 