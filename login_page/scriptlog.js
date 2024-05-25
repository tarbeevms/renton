document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const phone = document.getElementById('phone').value;
    const password = document.getElementById('password').value;
    
    const data = { phone: phone, password: password };
    
    fetch('/api/login/', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        document.getElementById('result').style.display = 'block';
        document.getElementById('result').textContent = JSON.stringify(data, null, 2);
    })
    .catch((error) => {
        console.error('Ошибка:', error);
    });
});