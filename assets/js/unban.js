const users = [
    { name: 'Utilisateur 1', reason: 'Violation des règles' },
    { name: 'Utilisateur 2', reason: 'Spam' },
    { name: 'Utilisateur 3', reason: 'Langage inapproprié' },
    { name: 'Utilisateur 4', reason: 'Comportement toxique' },
    { name: 'Utilisateur 5', reason: 'Multiples infractions' },
    { name: 'Utilisateur 6', reason: 'Contenu inapproprié' },
    { name: 'Utilisateur 7', reason: 'Harcelement' },
    { name: 'Utilisateur 8', reason: 'Cheating' },
    { name: 'Utilisateur 9', reason: 'Impersonation' }
];

const itemsPerPage = 6;
let currentPage = 0;

function renderUsers() {
    const unbanUsersContainer = document.getElementById('unban-users');
    unbanUsersContainer.innerHTML = '';

    const start = currentPage * itemsPerPage;
    const end = start + itemsPerPage;
    const paginatedUsers = users.slice(start, end);

    let column;
    paginatedUsers.forEach((user, index) => {
        if (index % 3 === 0) {
            column = document.createElement('div');
            column.classList.add('user-column');
            unbanUsersContainer.appendChild(column);
        }
        const userCard = document.createElement('div');
        userCard.classList.add('user-card');
        userCard.innerHTML = `<p><strong>${user.name}</strong></p><p>Raison : ${user.reason}</p><button class="unban-button">Débannir</button>`;
        column.appendChild(userCard);
    });
}

function searchUser() {
    const input = document.getElementById('search-bar').value.toLowerCase();
    const userCards = document.querySelectorAll('.user-card');

    userCards.forEach(card => {
        const username = card.querySelector('p strong').textContent.toLowerCase();
        if (username.includes(input)) {
            card.style.display = '';
        } else {
            card.style.display = 'none';
        }
    });
}

function nextPage() {
    if ((currentPage + 1) * itemsPerPage < users.length) {
        currentPage++;
        transitionEffect('next');
    }
}

function prevPage() {
    if (currentPage > 0) {
        currentPage--;
        transitionEffect('prev');
    }
}

function transitionEffect(direction) {
    const unbanUsersContainer = document.getElementById('unban-users');
    unbanUsersContainer.style.opacity = 0;

    setTimeout(() => {
        renderUsers();
        if (direction === 'next') {
            unbanUsersContainer.style.transform = 'translateX(100%)';
        } else {
            unbanUsersContainer.style.transform = 'translateX(-100%)';
        }
        unbanUsersContainer.style.opacity = 1;
        setTimeout(() => {
            unbanUsersContainer.style.transform = 'translateX(0)';
        }, 50);
    }, 400);
}

document.addEventListener('DOMContentLoaded', renderUsers);
