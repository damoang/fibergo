# fibergo

패키지 설치
```
go get -u github.com/gofiber/fiber/v2
go get -u github.com/go-sql-driver/mysql
go get -u github.com/joho/godotenv
```
실행 방법
```
go mod init project_name
go mod tidy
go run main.go
```


/fibergo
 ├── main.go           # Go API 서버 (API만 처리)
 ├── .env              # 환경 변수 (DB 정보 저장)
 ├── static/           # 정적 파일 (HTML, JS, CSS)
 │   ├── index.html    # 자유게시판 목록 (프론트엔드)
 │   ├── view.html     # 게시글 상세 페이지
 │   ├── script.js     # JavaScript (API 호출)
✅ static/ 폴더에 HTML & JavaScript를 배치하여 Go 서버와 분리
✅ Go 서버는 /free, /free/{id} API만 제공
✅ 정적 파일을 제공하여 클라이언트에서 직접 HTML & JS를 로드
