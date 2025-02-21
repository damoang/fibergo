package routes

import (
	"github.com/gofiber/fiber/v2"
)

// HandleBoardSSR handles server-side rendering for board pages
func HandleBoardSSR(c *fiber.Ctx) error {
	boardType := c.Params("type")
	
	// 게시판 타입 검증
	allowedBoards := map[string]bool{
		"free":    true,
		"notice":  true,
		"gallery": true,
	}
	
	if !allowedBoards[boardType] {
		return c.Status(400).JSON(fiber.Map{
			"error": "유효하지 않은 게시판입니다",
		})
	}

	// board.html 파일을 렌더링 (경로 수정)
	return c.SendFile("./templates/board.html")
}

// HandleBoardAPI handles API requests for board data
func HandleBoardAPI(c *fiber.Ctx) error {
	boardType := c.Params("type")
	
	// 게시판 타입 검증
	allowedBoards := map[string]bool{
		"free":    true,
		"notice":  true,
		"gallery": true,
	}
	
	if !allowedBoards[boardType] {
		return c.Status(400).JSON(fiber.Map{
			"error": "유효하지 않은 게시판입니다",
		})
	}

	// TODO: 게시판 데이터 조회 로직 구현
	return c.JSON(fiber.Map{
		"message": "게시판 API 구현 예정",
		"type": boardType,
	})
}