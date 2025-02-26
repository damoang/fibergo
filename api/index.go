// 파일: api/index.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/joho/godotenv"
    "net/http"
)

// 환경 변수 로드 등 초기화
var app *fiber.App

func init() {
    // (선택) .env 파일 로드 - 배포환경에선 .env 없고 Vercel 환경변수 사용
    _ = godotenv.Load()  
    app = fiber.New()
    app.Use(logger.New()) // 예: 로거 미들웨어
    // 정적 파일 서빙 설정 (필요한 경우)
    app.Static("/", "./static")

    // API 라우트 설정
    app.Get("/free", listFreeHandler)       // /free 목록 조회 핸들러
    app.Get("/free/:id", viewFreeHandler)   // /free/{id} 상세 조회 핸들러

    // 추가 라우트나 미들웨어 설정 가능
}

// Vercel이 호출하는 Handler 함수 - fiber 앱에 요청을 전달
func Handler(w http.ResponseWriter, r *http.Request) {
    // Fiber의 Context 경로 설정을 위해 RequestURI 조정
    r.RequestURI = r.URL.RequestURI()
    app.Handler()(w, r)  // Fiber 앱으로 요청 처리
}
