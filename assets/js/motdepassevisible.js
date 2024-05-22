document.addEventListener('DOMContentLoaded', () => {
    const passwordFields = [
        { input: document.getElementById('password'), checkbox: document.getElementById('show-password') },
        { input: document.getElementById('confirm-password'), checkbox: document.getElementById('show-confirm-password') }
    ].filter(field => field.input && field.checkbox);

    passwordFields.forEach(({ input, checkbox }) => {
        checkbox.addEventListener('change', () => {
            input.type = checkbox.checked ? 'text' : 'password';
        });
    });
});
