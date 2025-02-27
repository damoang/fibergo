package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"
)

var db *sql.DB

// InitDB initializes the database connection for routes
func InitDB(database *sql.DB) {
	db = database
}

// HandleBoardSSR handles server-side rendering for board pages
func HandleBoardSSR(c *fiber.Ctx) error {
	boardType := c.Params("type")
	postId := c.Params("id")
	
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

	// 상세 페이지도 SSR로 처리
	if postId != "" {
		// 게시글 데이터 조회
		tableName := "g5_write_" + boardType
		query := `
			SELECT wr_id, wr_subject, wr_name, wr_datetime, wr_hit, wr_good, wr_content
			FROM ` + tableName + `
			WHERE wr_id = ?
		`
		
		var post struct {
			ID       int    `json:"id"`
			제목     string `json:"제목"`
			이름     string `json:"이름"`
			날짜     string `json:"날짜"`
			조회     int    `json:"조회"`
			추천     int    `json:"추천"`
			내용     string `json:"내용"`
		}
		
		err := db.QueryRow(query, postId).Scan(
			&post.ID, &post.제목, &post.이름, &post.날짜, 
			&post.조회, &post.추천, &post.내용,
		)
		
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(404).SendFile("templates/404.html")
			}
			return c.Status(500).SendFile("templates/500.html")
		}

		// 조회수 증가
		updateQuery := `UPDATE ` + tableName + ` SET wr_hit = wr_hit + 1 WHERE wr_id = ?`
		_, err = db.Exec(updateQuery, postId)
		if err != nil {
			log.Printf("조회수 증가 실패: %v", err)
		}

		// SSR로 상세 페이지 렌더링
		return c.Render("board_view", fiber.Map{
			"Title": getBoardTitle(boardType),
			"BoardType": boardType,
			"Post": post,
		})
	}

	// 목록 페이지는 SSR로 처리
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	// DB에서 게시글 목록 조회
	tableName := "g5_write_" + boardType
	query := `
		SELECT wr_id, wr_subject, wr_name, wr_datetime, wr_hit, wr_good
		FROM ` + tableName + `
		WHERE wr_is_comment = 0
		ORDER BY wr_num DESC
		LIMIT ? OFFSET ?
	`

	// 전체 게시글 수 조회
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM ` + tableName + ` WHERE wr_is_comment = 0`
	err := db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return c.Status(500).SendFile("templates/500.html")
	}

	// 게시글 목록 조회
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return c.Status(500).SendFile("templates/500.html")
	}
	defer rows.Close()

	var posts []map[string]interface{}
	for rows.Next() {
		var (
			wr_id                             int
			wr_subject, wr_name, wr_datetime string
			wr_hit, wr_good                  int
		)

		err := rows.Scan(&wr_id, &wr_subject, &wr_name, &wr_datetime, &wr_hit, &wr_good)
		if err != nil {
			continue
		}

		posts = append(posts, map[string]interface{}{
			"id":       wr_id,
			"제목":      wr_subject,
			"작성자":     wr_name,
			"작성일":     wr_datetime,
			"조회수":     wr_hit,
			"추천수":     wr_good,
			"BoardType": boardType,
		})
	}

	// SSR 템플릿 렌더링
	return c.Render("board_list", fiber.Map{
		"BoardType": boardType,
		"Title":     getBoardTitle(boardType),
		"Posts":     posts,
		"Total":     totalCount,
		"Page":      page,
	})
}

func getBoardTitle(boardType string) string {
	titles := map[string]string{
		"free":    "자유게시판",
		"notice":  "공지사항",
		"gallery": "갤러리",
	}
	return titles[boardType]
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

	// 페이지 정보
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	offset := (page - 1) * limit

	// DB에서 게시글 목록 조회
	tableName := "g5_write_" + boardType
	query := `
		SELECT wr_id, wr_subject, wr_name, wr_datetime, wr_hit, wr_good
		FROM ` + tableName + `
		WHERE wr_is_comment = 0
		ORDER BY wr_num DESC
		LIMIT ? OFFSET ?
	`

	// 전체 게시글 수 조회
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM ` + tableName + ` WHERE wr_is_comment = 0`
	err := db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "게시글 수 조회 중 오류가 발생했습니다",
		})
	}

	// 게시글 목록 조회
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "게시글 목록 조회 중 오류가 발생했습니다",
		})
	}
	defer rows.Close()

	var posts []fiber.Map
	for rows.Next() {
		var (
			wr_id                                     int
			wr_subject, wr_name, wr_datetime         string
			wr_hit, wr_good                          int
		)

		err := rows.Scan(&wr_id, &wr_subject, &wr_name, &wr_datetime, &wr_hit, &wr_good)
		if err != nil {
			continue
		}

		posts = append(posts, fiber.Map{
			"id":     wr_id,
			"제목":    wr_subject,
			"작성자":   wr_name,
			"작성일":   wr_datetime,
			"조회수":   wr_hit,
			"추천수":   wr_good,
		})
	}

	return c.JSON(fiber.Map{
		"게시판":   boardType,
		"현재페이지": page,
		"전체개수":  totalCount,
		"게시글":   posts,
	})
}

// HandleCommentsAPI handles comment loading for a post
func HandleCommentsAPI(c *fiber.Ctx) error {
	boardType := c.Params("type")
	postId := c.Params("id")

	// 게시판 타입 검증
	if !isValidBoardType(boardType) {
		return c.Status(400).JSON(fiber.Map{
			"error": "유효하지 않은 게시판입니다",
		})
	}

	tableName := "g5_write_" + boardType
	query := `
		SELECT 
			wr_id,
			wr_content,
			wr_name,
			wr_datetime,
			wr_parent
		FROM ` + tableName + `
		WHERE wr_is_comment = 1 
		AND wr_parent = ?
		ORDER BY wr_comment, wr_comment_reply
	`

	rows, err := db.Query(query, postId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "댓글 조회 중 오류가 발생했습니다",
		})
	}
	defer rows.Close()

	var comments []fiber.Map
	for rows.Next() {
		var (
			wr_id                             int
			wr_content, wr_name, wr_datetime string
			wr_parent                        int
		)

		err := rows.Scan(&wr_id, &wr_content, &wr_name, &wr_datetime, &wr_parent)
		if err != nil {
			continue
		}

		comments = append(comments, fiber.Map{
			"id":     wr_id,
			"내용":    wr_content,
			"작성자":   wr_name,
			"날짜":    wr_datetime,
			"부모글ID": wr_parent,
		})
	}

	return c.JSON(fiber.Map{
		"count":    len(comments),
		"comments": comments,
	})
}