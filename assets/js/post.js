let currentCarouselIndex = 0;

function toggleComments() {
    const commentsSection = document.getElementById('comments-section');
    commentsSection.style.display = (commentsSection.style.display === 'none' || commentsSection.style.display === '') ? 'block' : 'none';
}

function toggleLike(button) {
    const img = button.querySelector('img');
    img.src = img.src.includes('likevide') ? '/assets/images/like.png' : '/assets/images/likevide.png';
}

function toggleDislike(button) {
    const img = button.querySelector('img');
    img.src = img.src.includes('dislikevide') ? '/assets/images/dislike.png' : '/assets/images/dislikevide.png';
}

function changeCarouselImage(direction) {
    const images = document.querySelectorAll('.carousel-image');
    images[currentCarouselIndex].style.display = 'none';
    currentCarouselIndex = (currentCarouselIndex + direction + images.length) % images.length;
    if (currentCarouselIndex < 0) {
        currentCarouselIndex = images.length - 1;
    }
    images[currentCarouselIndex].style.display = 'block';
}

document.addEventListener('DOMContentLoaded', () => {
    const images = document.querySelectorAll('.carousel-image');
    images[currentCarouselIndex].style.display = 'block';
});
