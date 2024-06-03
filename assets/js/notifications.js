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