package api

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
)

// HandleBoardSSR handles server-side rendering for board pages
func HandleBoardSSR(c *fiber.Ctx) error {
	boardType := c.Params("type")
	postId := c.Params("id")
	
	// 게시판 타입 검증
	if !isValidBoardType(boardType) {
		return c.Status(400).JSON(fiber.Map{
			"error": "유효하지 않은 게시판입니다",
		})
	}

	// ... (routes/board.go에서 복사) ...
}

// HandleBoardAPI handles API requests for board data
func HandleBoardAPI(c *fiber.Ctx) error {
	boardType := c.Params("type")
	
	// ... (routes/board.go에서 복사) ...
}

// HandleCommentsAPI handles comment loading for a post
func HandleCommentsAPI(c *fiber.Ctx) error {
	boardType := c.Params("type")
	postId := c.Params("id")

	// ... (routes/board.go에서 복사) ...
}

// 유틸리티 함수들
func isValidBoardType(boardType string) bool {
	allowedBoards := map[string]bool{
		"free":    true,
		"notice":  true,
		"gallery": true,
	}
	return allowedBoards[boardType]
}

func getBoardTitle(boardType string) string {
	titles := map[string]string{
		"free":    "자유게시판",
		"notice":  "공지사항",
		"gallery": "갤러리",
	}
	return titles[boardType]
} 