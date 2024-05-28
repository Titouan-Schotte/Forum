function setRating(element, category) {
    const index = parseInt(element.getAttribute('data-index'));
    const emojis = element.parentElement.children;
    for (let i = 0; i < emojis.length; i++) {
        if (i < index) {
            emojis[i].classList.add('selected');
            emojis[i].textContent = category === 'beauty' ? '🌸' : '💀';
        } else {
            emojis[i].classList.remove('selected');
            emojis[i].textContent = category === 'beauty' ? '☢️' : '😎';
        }
    }
    document.getElementById(`${category}-rating`).value = index;
}

function hoverEmoji(element, category) {
    const index = parseInt(element.getAttribute('data-index'));
    const emojis = element.parentElement.children;
    for (let i = 0; i < emojis.length; i++) {
        if (i < index) {
            emojis[i].textContent = category === 'beauty' ? '🌸' : '💀';
        }
    }
}

function unhoverEmoji(element, category) {
    const index = parseInt(element.getAttribute('data-index'));
    const emojis = element.parentElement.children;
    for (let i = 0; i < emojis.length; i++) {
        if (!emojis[i].classList.contains('selected')) {
            emojis[i].textContent = category === 'beauty' ? '☢️' : '😎';
        }
    }
}
