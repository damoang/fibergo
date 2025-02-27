package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"fibergo/routes" // 이 부분이 go.mod의 모듈명과 일치해야 함

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func main() {
	// 환경 변수 로드
	err := godotenv.Load()
	if err != nil {
		log.Fatal("환경 변수를 로드할 수 없습니다:", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	apiPort := os.Getenv("API_PORT")

	// 데이터베이스 연결
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("데이터베이스 연결 실패:", err)
	}
	defer db.Close()

	// DB 연결을 routes 패키지에 전달
	routes.InitDB(db)

	// DB 커넥션 풀 설정
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// DB 연결 확인
	if err := db.Ping(); err != nil {
		log.Fatal("데이터베이스 연결 확인 실패:", err)
	}

	// 템플릿 엔진 설정을 더 자세하게
	engine := html.New("./templates", ".html")
	engine.Reload(true)           // 개발 환경에서 템플릿 자동 리로드
	engine.Debug(true)            // 디버그 모드 활성화
	engine.Layout("layouts/main") // 기본 레이아웃 설정

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main", // 기본 레이아웃 설정
		Prefork:     false,
		// 캐시 설정 추가
		CacheControl: true,
		// 압축 설정
		Compression: true,
		// 에러 핸들링
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			// Accept 헤더에 따라 JSON 또는 HTML 응답
			if c.Accepts("json") {
				return c.Status(code).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return c.Status(code).Render("error", fiber.Map{
				"Title": "오류가 발생했습니다",
				"Error": err.Error(),
			})
		},
	})

	// CORS 미들웨어 개선
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")

		// 보안 헤더 추가
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(204)
		}
		return c.Next()
	})

	// 정적 파일 제공
	app.Static("/static", "./static")

	// 루트 경로 처리
	app.Get("/", func(c *fiber.Ctx) error {
		// DB 상태 확인
		dbStatus := "연결 성공"
		if err := db.Ping(); err != nil {
			dbStatus = "연결 실패: " + err.Error()
		}

		// 서버 정보 수집
		serverInfo := map[string]interface{}{
			"status":    "정상 작동 중",
			"version":   "1.0.0",
			"startTime": time.Now().Format("2006-01-02 15:04:05"),
			"database": map[string]string{
				"status": dbStatus,
				"host":   os.Getenv("DB_HOST"),
				"name":   os.Getenv("DB_NAME"),
			},
			"boards": map[string]string{
				"free":    "자유게시판",
				"notice":  "공지사항",
				"gallery": "갤러리",
			},
			"endpoints": map[string]string{
				"boards":   "/api/:type",
				"post":     "/api/:type/:id",
				"comments": "/api/:type/:id/comments",
			},
		}

		return c.JSON(fiber.Map{
			"message": "Board API Server",
			"server":  serverInfo,
		})
	})

	// API 라우트
	apiGroup := app.Group("/api")

	// 게시판 목록 API
	apiGroup.Get("/:type", routes.HandleBoardAPI)

	// 게시글 상세 조회 API
	apiGroup.Get("/:type/:id", func(c *fiber.Ctx) error {
		boardType := c.Params("type")
		wrID := c.Params("id")

		// 입력값 검증 추가
		if wrID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "잘못된 게시글 ID입니다",
			})
		}

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

		tableName := "g5_write_" + boardType

		// Prepared Statement 사용
		query := `SELECT wr_id, wr_subject, wr_name, wr_datetime, wr_hit, wr_good, wr_content 
                  FROM ?? 
                  WHERE wr_id = ? AND wr_is_comment = 0`

		// 실제 쿼리 생성 (더 안전한 방식)
		query = strings.Replace(query, "??", tableName, 1)

		var wr_id, wr_hit, wr_good int
		var wr_subject, wr_name, wr_datetime, wr_content string

		err := db.QueryRow(query, wrID).Scan(&wr_id, &wr_subject, &wr_name, &wr_datetime, &wr_hit, &wr_good, &wr_content)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(404).JSON(fiber.Map{
					"error": "게시글을 찾을 수 없습니다",
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"error": "서버 오류가 발생했습니다",
			})
		}

		// 조회수 증가 쿼리도 같은 방식으로 수정
		updateQuery := strings.Replace("UPDATE ?? SET wr_hit = wr_hit + 1 WHERE wr_id = ?", "??", tableName, 1)
		_, err = db.Exec(updateQuery, wrID)
		if err != nil {
			log.Printf("조회수 증가 실패: %v", err)
		}

		// 날짜 변환
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", wr_datetime)
		formattedTime := parsedTime.Format("2006-01-02 15:04:05")

		return c.JSON(fiber.Map{
			"id": wr_id,
			"추천": wr_good,
			"제목": wr_subject,
			"이름": wr_name,
			"날짜": formattedTime,
			"조회": wr_hit,
			"내용": wr_content,
		})
	})

	// 댓글 API
	apiGroup.Get("/:type/:id/comments", routes.HandleCommentsAPI)

	// 웹 페이지 라우트
	app.Get("/:type", routes.HandleBoardSSR)
	app.Get("/:type/:id", routes.HandleBoardSSR)

	// 404 에러 핸들러
	app.Use(func(c *fiber.Ctx) error {
		// Accept 헤더 확인
		accepts := c.Accepts("html", "json")
		if accepts == "json" {
			return c.Status(404).JSON(fiber.Map{
				"error": "요청하신 페이지를 찾을 수 없습니다",
			})
		}
		// HTML 응답
		return c.Status(404).SendFile("templates/404.html")
	})

	log.Printf("🚀 서버가 http://localhost:%s 에서 실행 중...", apiPort)
	log.Fatal(app.Listen(":" + apiPort))
}
