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
    _ = godotenv.Load()

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    apiPort := os.Getenv("API_PORT")

    dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    app := fiber.New()

    // 자유게시판 목록 조회 API
    app.Get("/free", func(c *fiber.Ctx) error {
        query := `SELECT wr_id, wr_subject, wr_name, wr_datetime, wr_hit, wr_good 
                  FROM g5_write_free 
                  ORDER BY wr_datetime DESC 
                  LIMIT 10`

        rows, err := db.Query(query)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        defer rows.Close()

        var posts []map[string]interface{}

        for rows.Next() {
            var wr_id, wr_hit, wr_good int
            var wr_subject, wr_name, wr_datetime string

            if err := rows.Scan(&wr_id, &wr_subject, &wr_name, &wr_datetime, &wr_hit, &wr_good); err != nil {
                return c.Status(500).JSON(fiber.Map{"error": err.Error()})
            }

            // 날짜 형식 변환 (YYYY-MM-DD HH:MM:SS)
            parsedTime, _ := time.Parse("2006-01-02 15:04:05", wr_datetime)
            formattedTime := parsedTime.Format("2006-01-02 15:04:05")

            posts = append(posts, fiber.Map{
                "id":    wr_id, // 상세 조회를 위해 ID 포함
                "추천":  wr_good,
                "제목":  wr_subject,
                "이름":  wr_name,
                "날짜":  formattedTime,
                "조회":  wr_hit,
            })
        }

        return c.JSON(posts)
    })

    // 게시글 상세 조회 API
    app.Get("/free/:id", func(c *fiber.Ctx) error {
        wrID := c.Params("id")

        query := `SELECT wr_id, wr_subject, wr_name, wr_datetime, wr_hit, wr_good, wr_content 
                  FROM g5_write_free 
                  WHERE wr_id = ?`

        var wr_id, wr_hit, wr_good int
        var wr_subject, wr_name, wr_datetime, wr_content string

        err := db.QueryRow(query, wrID).Scan(&wr_id, &wr_subject, &wr_name, &wr_datetime, &wr_hit, &wr_good, &wr_content)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "게시글을 찾을 수 없습니다."})
        }

        // 날짜 형식 변환 (YYYY-MM-DD HH:MM:SS)
        parsedTime, _ := time.Parse("2006-01-02 15:04:05", wr_datetime)
        formattedTime := parsedTime.Format("2006-01-02 15:04:05")

        return c.JSON(fiber.Map{
            "id":    wr_id,
            "추천":  wr_good,
            "제목":  wr_subject,
            "이름":  wr_name,
            "날짜":  formattedTime,
            "조회":  wr_hit,
            "내용":  wr_content, // 상세 조회에서 내용 포함
        })
    })

    log.Printf("서버가 http://localhost:%s/free 에서 실행 중...", apiPort)
    log.Fatal(app.Listen(":" + apiPort))
}
