 <canvas id="confetti-canvas" style="position: fixed; top: 0; left: 0; width: 100%; height: 100%; pointer-events: none; z-index: 9999;"></canvas>
    <script src="https://cdn.jsdelivr.net/npm/canvas-confetti@1.6.1/dist/confetti.browser.min.js"></script>
    <script>
        document.addEventListener('DOMContentLoaded', function () {
            const canvas = document.getElementById('confetti-canvas');
            if (!canvas) return;

            // canvas 크기 설정
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;

            const confettiInstance = confetti.create(canvas, {
                resize: true,
                useWorker: true
            });

            // 자동 실행 (상단 중앙에서 1회)
            // confettiInstance({
            //     particleCount: 200,
            //     spread: 160,
            //     origin: { x: 0.5, y: 0 },
            //     colors: ['#ff6347', '#ffa500', '#4caf50', '#1e90ff', '#9400d3'],
            //     scalar: 1.2
            // });

            confettiInstance({
                particleCount: 200,
                angle: 90,
                spread: 360,
                startVelocity: 50,
                decay: 0.92,
                gravity: 1,
                ticks: 180,
                origin: { x: 0.5, y: 0.3 }, // 정중앙보다 약간 위쪽
                colors: ['#ff6347', '#ffa500', '#4caf50', '#1e90ff', '#9400d3'],
                scalar: 1.4
            });


            // 클릭 시 클릭한 위치에서 폭죽
            window.addEventListener('click', function (e) {
                const x = e.clientX / window.innerWidth;
                const y = e.clientY / window.innerHeight;

                confettiInstance({
                    particleCount: 150,
                    spread: 120,
                    origin: { x: x, y: y },
                    colors: ['#ff6347', '#ffa500', '#4caf50', '#1e90ff', '#9400d3'],
                    scalar: 1.2
                });
            });

            // 5초 후 자동 애니메이션 종료 (클릭 효과는 계속됨)
            setTimeout(() => {
                confettiInstance.reset(); // 이건 자동 효과만 초기화
            }, 5000);
        });
    </script>
