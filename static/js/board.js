// 게시글 동적 로딩
async function loadPosts(page = 1) {
    try {
        const boardType = window.location.pathname.split('/')[1];
        const response = await fetch(`/api/${boardType}?page=${page}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('게시글 로딩 실패:', error);
        throw error;
    }
}

// 게시글 상세 정보 로딩
async function loadPostDetail(boardType, postId) {
    try {
        const response = await fetch(`/api/${boardType}/${postId}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('게시글 상세 정보 로딩 실패:', error);
        throw error;
    }
}

document.addEventListener('DOMContentLoaded', function() {
    const path = window.location.pathname;
    const pathParts = path.split('/').filter(Boolean);
    const boardType = pathParts[0];
    const postId = pathParts[1];

    // 상세 페이지일 때만 CSR 처리
    if (postId) {
        const boardContent = document.getElementById('board-content');
        
        try {
            fetch(`/api/${boardType}/${postId}`)
                .then(response => response.json())
                .then(post => {
                    const date = new Date(post.날짜).toLocaleString();
                    boardContent.innerHTML = `
                        <div class="post-view">
                            <h1>${getBoardTitle(boardType)}</h1>
                            <div class="post-header">
                                <h2 class="post-title">${post.제목}</h2>
                                <div class="post-info">
                                    <span class="author">작성자: ${post.이름}</span>
                                    <span class="date">작성일: ${date}</span>
                                    <span class="views">조회: ${post.조회}</span>
                                    <span class="likes">추천: ${post.추천}</span>
                                </div>
                            </div>
                            <div class="post-content">${post.내용}</div>
                            <div class="post-actions">
                                <button onclick="location.href='/${boardType}'">목록</button>
                            </div>
                        </div>
                    `;
                })
                .catch(error => {
                    console.error('Error:', error);
                    boardContent.innerHTML = `
                        <div class="error">게시글을 불러오는 중 오류가 발생했습니다.</div>
                    `;
                });
        } catch (error) {
            console.error('Error:', error);
            boardContent.innerHTML = `
                <div class="error">게시글을 불러오는 중 오류가 발생했습니다.</div>
            `;
        }
    } else {
        // 목록 페이지는 이미 SSR로 렌더링됨
        // 클릭 이벤트만 처리
        const postList = document.getElementById('postList');
        if (postList) {
            postList.addEventListener('click', (e) => {
                const titleCell = e.target.closest('.title');
                if (titleCell) {
                    const postId = titleCell.dataset.id;
                    window.location.href = `/${boardType}/${postId}`;
                }
            });
        }
    }
});

function getBoardTitle(boardType) {
    const titles = {
        'free': '자유게시판',
        'notice': '공지사항',
        'gallery': '갤러리'
    };
    return titles[boardType] || '게시판';
} 