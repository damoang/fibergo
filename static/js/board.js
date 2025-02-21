// 게시글 동적 로딩
async function loadPosts(page = 1) {
    try {
        const response = await fetch(`/api/board/${BOARD_TYPE}?page=${page}`);
        const posts = await response.json();
        updatePostList(posts);
    } catch (error) {
        console.error('게시글 로딩 실패:', error);
    }
}

// 게시글 클릭 이벤트
document.getElementById('postList').addEventListener('click', (e) => {
    const postItem = e.target.closest('.post-item');
    if (postItem) {
        const postId = postItem.dataset.id;
        window.location.href = `/go/board.html/${BOARD_TYPE}/view/${postId}`;
    }
}); 