function openNotifications() {
    document.getElementById('notifications-panel').style.display = 'block';
}

function closeNotifications() {
    document.getElementById('notifications-panel').style.display = 'none';
}


// Mock notifications for demonstration
document.addEventListener('DOMContentLoaded', () => {
    const notificationsList = document.getElementById('notifications-list');

    const notifications = [
        { type: 'like', message: 'Votre post a reçu un like' },
        { type: 'dislike', message: 'Votre post a reçu un dislike' },
        { type: 'comment', message: 'Un commentaire a été posté sous votre post' },
        { type: 'success', message: 'Votre post a été créé avec succès' }
    ];

    notifications.forEach(notification => {
        const notificationElement = document.createElement('div');
        notificationElement.classList.add('notification');

        const icon = document.createElement('img');
        switch (notification.type) {
            case 'like':
                icon.src = '/assets/images/like.png';
                break;
            case 'dislike':
                icon.src = '/assets/images/dislike.png';
                break;
            case 'comment':
                icon.src = '/assets/images/comment.png';
                break;
            case 'success':
                icon.src = '/assets/images/success.png';
                break;
        }

        notificationElement.appendChild(icon);
        notificationElement.appendChild(document.createTextNode(notification.message));
        notificationsList.appendChild(notificationElement);
    });
});
