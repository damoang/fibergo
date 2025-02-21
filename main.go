package main

import (
    "database/sql"
    "log"
    "os"
    "strings"
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
        EnableGzip:   true,
        Prefork:      true,
        // 에러 핸들러 추가
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
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

    // 📌 (1) 정적 파일 제공 (HTML, JS)
    app.Static("/", "./static")

    // SSR 라우트
    app.Get("/:type", HandleBoardSSR)
    
    // API 라우트
    app.Get("/api/:type", HandleBoardAPI)

    // 게시글 상세 조회 API
    app.Get("/board/:type/:id", func(c *fiber.Ctx) error {
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
            "free": true,
            "notice": true,
            "gallery": true,
            // 필요한 게시판 타입 추가
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
            "id":    wr_id,
            "추천":  wr_good,
            "제목":  wr_subject,
            "이름":  wr_name,
            "날짜":  formattedTime,
            "조회":  wr_hit,
            "내용":  wr_content,
        })
    })

    // 댓글 조회 API
    app.Get("/board/:type/:id/comments", func(c *fiber.Ctx) error {
        boardType := c.Params("type")
        wrParentID := c.Params("id")
        
        if wrParentID == "" {
            return c.Status(400).JSON(fiber.Map{
                "error": "잘못된 게시글 ID입니다",
            })
        }

        // 게시판 타입 검증
        allowedBoards := map[string]bool{
            "free": true,
            "notice": true,
            "gallery": true,
        }
        
        if !allowedBoards[boardType] {
            return c.Status(400).JSON(fiber.Map{
                "error": "유효하지 않은 게시판입니다",
            })
        }

        tableName := "g5_write_" + boardType
        
        // Prepared Statement 사용
        query := strings.Replace(
            `SELECT wr_id, wr_parent, wr_content, wr_name, wr_datetime 
             FROM ?? 
             WHERE wr_parent = ? AND wr_is_comment = 1  
             ORDER BY wr_datetime ASC`,
            "??",
            tableName,
            1,
        )

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

        // 에러 처리 개선
        if err := rows.Err(); err != nil {
            return c.Status(500).JSON(fiber.Map{
                "error": "댓글 조회 중 오류가 발생했습니다",
            })
        }

        return c.JSON(fiber.Map{
            "comments": comments,
            "count":    len(comments),
        })
    })

    // 404 에러 핸들러 개선
    app.Use(func(c *fiber.Ctx) error {
        return c.Status(404).JSON(fiber.Map{
            "error": "요청하신 페이지를 찾을 수 없습니다",
        })
    })

    log.Printf("🚀 서버가 http://localhost:%s 에서 실행 중...", apiPort)
    log.Fatal(app.Listen(":" + apiPort))
}
