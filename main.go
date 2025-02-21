package main

import (
    "database/sql"
    "log"
    "os"
    "time"

    "github.com/gofiber/fiber/v2"
    _ "github.com/go-sql-driver/mysql"
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

    // DB 커넥션 풀 설정
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // DB 연결 확인
    if err := db.Ping(); err != nil {
        log.Fatal("데이터베이스 연결 확인 실패:", err)
    }

    app := fiber.New(fiber.Config{
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
        EnableGzip: true,
        Prefork: true,
    })

    app.Use(func(c *fiber.Ctx) error {
        c.Set("Access-Control-Allow-Origin", "*")
        c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        return c.Next()
    })

    // 📌 (1) 정적 파일 제공 (HTML, JS)
    app.Static("/", "./static")

    // SSR 라우트
    app.Get("/:type", HandleBoardSSR)
    
    // API 라우트
    app.Get("/api/:type", HandleBoardAPI)

    // 📌 (3) 게시글 상세 조회 API
    app.Get("/free/:id", func(c *fiber.Ctx) error {
        wrID := c.Params("id")

        query := `SELECT wr_id, wr_subject, wr_name, wr_datetime, wr_hit, wr_good, wr_content 
                  FROM g5_write_free 
                  WHERE wr_id = ? AND wr_is_comment = 0`  /* 게시글만 조회 */

        var wr_id, wr_hit, wr_good int
        var wr_subject, wr_name, wr_datetime, wr_content string

        err := db.QueryRow(query, wrID).Scan(&wr_id, &wr_subject, &wr_name, &wr_datetime, &wr_hit, &wr_good, &wr_content)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "게시글을 찾을 수 없습니다."})
        }

        // 날짜 변환
        parsedTime, _ := time.Parse("2006-01-02 15:04:05", wr_datetime)
        formattedTime := parsedTime.Format("2006-01-02 15:04:05")

        return c.JSON(fiber.Map{
            "id":    wr_id,
            "추천":  wr_good,
            "제목":  wr_subject,
            "이름":  wr_name,
            "날짜":  formattedTime,
            "조회":  wr_hit,
            "내용":  wr_content,
        })
    })

    // 📌 (4) 특정 게시글의 댓글 조회 API
    app.Get("/free/:id/comments", func(c *fiber.Ctx) error {
        wrParentID := c.Params("id")

        query := `SELECT wr_id, wr_parent, wr_content, wr_name, wr_datetime 
                  FROM g5_write_free 
                  WHERE wr_parent = ? AND wr_is_comment = 1  
                  ORDER BY wr_datetime ASC`

        rows, err := db.Query(query, wrParentID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        defer rows.Close()

        var comments []map[string]interface{}

        for rows.Next() {
            var wr_id, wr_parent int
            var wr_content, wr_name, wr_datetime string

            if err := rows.Scan(&wr_id, &wr_parent, &wr_content, &wr_name, &wr_datetime); err != nil {
                return c.Status(500).JSON(fiber.Map{"error": err.Error()})
            }

            comments = append(comments, fiber.Map{
                "댓글ID":  wr_id,
                "부모ID":  wr_parent,
                "내용":    wr_content,
                "작성자":  wr_name,
                "날짜":    wr_datetime,
            })
        }

        return c.JSON(comments)
    })

    // 에러 핸들러 추가
    app.Use(func(c *fiber.Ctx) error {
        return c.Status(404).SendString("페이지를 찾을 수 없습니다")
    })

    log.Printf("🚀 서버가 http://localhost:%s 에서 실행 중...", apiPort)
    log.Fatal(app.Listen(":" + apiPort))
}
