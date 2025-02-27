package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/go-sql-driver/mysql"
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
	})

	// 라우트 설정
	setupRoutes()
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.Handler()(w, r)
}
