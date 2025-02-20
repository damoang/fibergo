package main
import (
    "database/sql"
    "log"
    "os"
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
    app.Get("/free", func(c *fiber.Ctx) error {
        rows, err := db.Query("SELECT wr_id, wr_subject, wr_content, wr_datetime FROM g5_write_free ORDER BY wr_datetime DESC LIMIT 10")
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": err.Error()})
        }
        defer rows.Close()
        var posts []map[string]interface{}
        for rows.Next() {
            var wr_id int
            var wr_subject, wr_content, wr_datetime string
            if err := rows.Scan(&wr_id, &wr_subject, &wr_content, &wr_datetime); err != nil {
                return c.Status(500).JSON(fiber.Map{"error": err.Error()})
            }
            posts = append(posts, fiber.Map{
                "id": wr_id,
                "title": wr_subject,
                "content": wr_content,
                "datetime": wr_datetime,
            })
        }
        return c.JSON(posts)
    })
    log.Printf("서버가 http://localhost:%s/free 에서 실행 중...", apiPort)
    log.Fatal(app.Listen(":" + apiPort))
}
