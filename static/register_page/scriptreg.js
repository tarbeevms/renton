document.getElementById('registerForm').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const phone = document.getElementById('phone').value;
    const password = document.getElementById('password').value;
    const account = document.getElementById('account').value;
    const surname = document.getElementById('surname').value;
    const name = document.getElementById('firstname').value;
    
    const data = {
        phone: phone,
        password: password,
        account: account,
        surname: surname,
        firstname: firstname
    };
    
    fetch('/api/register/', {
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