document.addEventListener('DOMContentLoaded', function() {
    // Exemple de notifications
    const notifications = [
        { message: 'Votre post a été créé avec succès.', type: 'success' },
        { message: 'Vous avez reçu un like sur votre post.', type: 'like' },
        { message: 'Vous avez reçu un dislike sur votre post.', type: 'dislike' },
        { message: 'Un utilisateur a commenté votre post.', type: 'comment' },
    ];

    const notificationsList = document.getElementById('notifications-list');
    const notificationCount = document.getElementById('notification-count');

    const icons = {
        success: '/assets/images/success.png',
        like: '/assets/images/like.png',
        dislike: '/assets/images/dislike.png',
        comment: '/assets/images/message.png'
    };

    notifications.forEach((notification, index) => {
        const notificationItem = document.createElement('div');
        notificationItem.className = 'notification-item';

        const img = document.createElement('img');
        img.src = icons[notification.type];
        img.alt = notification.type;

        const message = document.createElement('span');
        message.innerText = notification.message;

        notificationItem.appendChild(img);
        notificationItem.appendChild(message);
        notificationsList.appendChild(notificationItem);
    });

    notificationCount.innerText = notifications.length;
});
