document.querySelectorAll('.beauty-icon').forEach(icon => {
    icon.addEventListener('click', function() {
        const rating = parseInt(this.getAttribute('data-value'), 10);
        setBeautyRating(rating);
        const beautyIcons = document.querySelectorAll('.beauty-icon');
        beautyIcons.forEach((icon, index) => {
            if (index < rating) {
                icon.classList.add('selected');
            } else {
                icon.classList.remove('selected');
            }
        });
    });
});

document.querySelectorAll('.danger-icon').forEach(icon => {
    icon.addEventListener('click', function() {
        const rating = parseInt(this.getAttribute('data-value'), 10);
        setDangerRating(rating);
        const dangerIcons = document.querySelectorAll('.danger-icon');
        dangerIcons.forEach((icon, index) => {
            if (index < rating) {
                icon.classList.add('selected');
            } else {
                icon.classList.remove('selected');
            }
        });
    });
});
