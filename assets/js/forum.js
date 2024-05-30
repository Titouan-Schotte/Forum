let currentBeautyRating = 0;
let currentDangerRating = 0;

function openTab(evt, tabName) {
    let i, tabcontent, tablinks;
    tabcontent = document.getElementsByClassName('tabcontent');
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = 'none';
    }
    tablinks = document.getElementsByClassName('tablink');
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(' active', '');
    }
    document.getElementById(tabName).style.display = 'block';
    evt.currentTarget.className += ' active';
}

function filterPosts() {
    const searchQuery = document.getElementById('search-bar').value.toLowerCase();
    const selectedCategories = Array.from(document.querySelectorAll('.category:checked')).map(cb => cb.value);
    const posts = document.querySelectorAll('.post');

    posts.forEach(post => {
        const title = post.textContent.toLowerCase();
        const categories = post.getAttribute('data-categories').split(',');
        const postBeauty = parseInt(post.getAttribute('data-beauty'), 10);
        const postDanger = parseInt(post.getAttribute('data-danger'), 10);

        const matchesSearch = title.includes(searchQuery);
        const matchesCategory = selectedCategories.length === 0 || selectedCategories.some(category => categories.includes(category));
        const matchesBeauty = postBeauty >= currentBeautyRating;
        const matchesDanger = postDanger >= currentDangerRating;

        if (matchesSearch && matchesCategory && matchesBeauty && matchesDanger) {
            post.style.display = 'block';
        } else {
            post.style.display = 'none';
        }
    });
}

function setBeautyRating(rating) {
    currentBeautyRating = rating;
    const beautyIcons = document.querySelectorAll('.beauty-icon');
    beautyIcons.forEach(icon => {
        icon.classList.remove('active');
        if (parseInt(icon.getAttribute('data-value'), 10) <= rating) {
            icon.classList.add('active');
        }
    });
    filterPosts();
}

function setDangerRating(rating) {
    currentDangerRating = rating;
    const dangerIcons = document.querySelectorAll('.danger-icon');
    dangerIcons.forEach(icon => {
        icon.classList.remove('active');
        if (parseInt(icon.getAttribute('data-value'), 10) <= rating) {
            icon.classList.add('active');
        }
    });
    filterPosts();
}

function openNotifications() {
    document.getElementById('notifications-panel').style.display = 'block';
}

function closeNotifications() {
    document.getElementById('notifications-panel').style.display = 'none';
}

document.querySelectorAll('.tablink').forEach(tab => {
    tab.addEventListener('click', function() {
        openTab(event, this.getAttribute('data-tab'));
    });
});

document.getElementById('search-bar').addEventListener('input', filterPosts);
document.querySelectorAll('.category').forEach(category => {
    category.addEventListener('change', filterPosts);
});
document.querySelectorAll('.beauty-icon').forEach(icon => {
    icon.addEventListener('click', function() {
        setBeautyRating(parseInt(this.getAttribute('data-value'), 10));
    });
});
document.querySelectorAll('.danger-icon').forEach(icon => {
    icon.addEventListener('click', function() {
        setDangerRating(parseInt(this.getAttribute('data-value'), 10));
    });
});

// Définit l'onglet par défaut (Derniers Posts)
document.querySelector('.tablink.active').click();
