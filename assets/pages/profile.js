document.addEventListener('DOMContentLoaded', () => {
    const toggleModeButton = document.getElementById('toggleMode');
    const body = document.body;

    toggleModeButton.addEventListener('click', () => {
        if (body.classList.contains('dark-mode')) {
            body.classList.remove('dark-mode');
            body.classList.add('light-mode');
        } else {
            body.classList.remove('light-mode');
            body.classList.add('dark-mode');
        }
    });

    // Simuler le remplissage des données de profil
    document.getElementById('username').textContent = 'JohnDoe';
    document.getElementById('followers').textContent = '120';
    document.getElementById('following').textContent = '200';
    document.getElementById('posts').textContent = '15';
    document.getElementById('likes').textContent = '350';

    // Simuler les posts créés
    const createdPosts = document.getElementById('created-posts');
    for (let i = 0; i < 5; i++) {
        const post = document.createElement('div');
        post.textContent = `Post créé ${i + 1}`;
        post.classList.add('post');
        createdPosts.appendChild(post);
    }

    // Simuler les posts likés
    const likedPosts = document.getElementById('liked-posts');
    for (let i = 0; i < 3; i++) {
        const post = document.createElement('div');
        post.textContent = `Post liké ${i + 1}`;
        post.classList.add('post');
        likedPosts.appendChild(post);
    }
});
