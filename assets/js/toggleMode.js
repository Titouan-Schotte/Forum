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
});