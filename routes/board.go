// SSR 라우트 핸들러
func HandleFreeBoardSSR(c *fiber.Ctx) error {
    // 게시글 데이터 조회
    posts, err := GetPosts(/* 필요한 파라미터 */)
    if err != nil {
        return c.Status(500).SendString("서버 오류가 발생했습니다")
    }

    // 템플릿 렌더링
    return c.Render("free", fiber.Map{
        "Title": "자유게시판",
        "Posts": posts,
    })
}

// API 라우트 핸들러 (기존 코드 유지)
func HandleFreeBoardAPI(c *fiber.Ctx) error {
    // 기존 API 로직
} 