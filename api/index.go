package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var app *fiber.App

func init() {
	_ = godotenv.Load() // .env 파일 로드 (로컬 개발용, Vercel에서는 필요 없음)
	app = fiber.New()

	// 환경 변수 로드
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	apiPort := os.Getenv("API_PORT")
	allowedBoards := os.Getenv("ALLOWED_BOARDS")

	log.Println("DB Connection:", dbUser, dbPassword, dbHost, dbPort, dbName)
	log.Println("API_PORT:", apiPort)
	log.Println("ALLOWED_BOARDS:", allowedBoards)

	// API 라우트 추가
	app.Get("/free", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Free Board"})
	})
	app.Get("/notice", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Notice Board"})
	})
}

// ✅ Vercel 서버리스 함수 진입점 (Handler)
func Handler(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = r.URL.RequestURI()
	app.Handler()(w, r)
}
