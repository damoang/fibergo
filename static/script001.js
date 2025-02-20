document.addEventListener("DOMContentLoaded", async function() {
    const postContainer = document.getElementById("post-list");
    const postDetailContainer = document.getElementById("post-detail");

    // 📌 게시판 목록 (`free.html`)
    if (postContainer) {
        const apiUrl = "https://staging.damoang.net/go/api/free"; // API 호출 경로

        try {
            const response = await fetch(apiUrl);
            const posts = await response.json();

            console.log(posts); // API 응답 확인

            postContainer.innerHTML = "";
            posts.forEach(post => {
                const title = post.제목 || "제목 없음";
                const name = post.이름 || "익명";
                const date = post.날짜 ? new Date(post.날짜).toLocaleString() : "날짜 없음";

                const postElement = document.createElement("div");
                postElement.className = "post";
                postElement.innerHTML = `
                    <h3>
                        <a href="view.html?id=${post.id}" class="post-link">${title}</a>
                    </h3>
                    <p>작성자: ${name} | 추천: ${post.추천} | 조회: ${post.조회}</p>
                    <small>${date}</small>
                `;
                postContainer.appendChild(postElement);
            });
        } catch (error) {
            console.error("API 호출 오류:", error);
            postContainer.innerHTML = "<p>게시글을 불러오는 데 실패했습니다.</p>";
        }
    }

    // 📌 게시글 상세보기 (`view.html`)
    if (postDetailContainer) {
        const params = new URLSearchParams(window.location.search);
        const wr_id = params.get("id"); // URL에서 id 값 가져오기

        if (!wr_id) {
            postDetailContainer.innerHTML = "<p>잘못된 접근입니다.</p>";
            return;
        }

        const apiUrl = `/go/api/free/${wr_id}`; // 게시글 상세 API 호출
        try {
            const response = await fetch(apiUrl);
            if (!response.ok) throw new Error("게시글을 불러올 수 없습니다.");

            const post = await response.json();

            postDetailContainer.innerHTML = `
                <h2>${post.제목}</h2>
                <p><strong>작성자:</strong> ${post.이름} | <strong>날짜:</strong> ${new Date(post.날짜).toLocaleString()}</p>
                <p><strong>추천:</strong> ${post.추천} | <strong>조회:</strong> ${post.조회}</p>
                <hr>
                <p>${post.내용.replace(/\n/g, "<br>")}</p>
            `;
        } catch (error) {
            console.error("API 호출 오류:", error);
            postDetailContainer.innerHTML = "<p>게시글을 불러오는 데 실패했습니다.</p>";
        }
    }
});

async function loadPosts() {
    try {
        const response = await fetch('/free');
        const posts = await response.json();
        
        const postsHTML = posts.map(post => {
            // 날짜 포맷팅
            const formattedDate = formatDate(post.날짜);
            
            return `
                <tr>
                    <td class="text-center">${post.id}</td>
                    <td>
                        <a href="/view.html?id=${post.id}" class="title-link">
                            ${post.제목}
                            ${post.댓글수 > 0 ? `<span class="comment-count">[${post.댓글수}]</span>` : ''}
                        </a>
                    </td>
                    <td class="text-center">${post.이름}</td>
                    <td class="text-center">${formattedDate}</td>
                    <td class="text-center">${post.조회}</td>
                    <td class="text-center">${post.추천}</td>
                </tr>
            `;
        }).join('');
        
        document.getElementById('posts').innerHTML = postsHTML;
    } catch (error) {
        console.error('게시글 로딩 실패:', error);
        document.getElementById('posts').innerHTML = `
            <tr>
                <td colspan="6" class="text-center">게시글을 불러오는데 실패했습니다.</td>
            </tr>
        `;
    }
}

function formatDate(dateString) {
    const date = new Date(dateString);
    const today = new Date();
    
    // 오늘 날짜인 경우 시간만 표시
    if (date.toDateString() === today.toDateString()) {
        return date.toLocaleTimeString('ko-KR', { 
            hour: '2-digit', 
            minute: '2-digit' 
        });
    }
    
    // 다른 날짜의 경우 년-월-일 표시
    return date.toLocaleDateString('ko-KR', { 
        year: '2-digit', 
        month: '2-digit', 
        day: '2-digit' 
    }).replace(/\./g, '-').replace(/\s/g, '');
}

// 페이지 로드 시 게시글 목록 불러오기
document.addEventListener('DOMContentLoaded', loadPosts);
