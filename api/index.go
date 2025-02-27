package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var app *fiber.App
var db *sql.DB

func init() {
	// DB 연결
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// DB 커넥션 풀 설정
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 템플릿 엔진 설정
	engine := html.New("./templates", ".html")
	engine.Reload(true)
	engine.Debug(true)
	engine.Layout("layouts/main")

	// Fiber 앱 초기화
	app = fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
		Prefork:     false,
		// 캐시 설정 추가
		CacheControl: true,
		// 압축 설정
		Compression: true,
	})

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

// Handler is the Vercel serverless function entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	app.Handler()(w, r)
}
