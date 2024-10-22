const params = new URLSearchParams(document.location.search);
const emailParam = params.get('email');
if (emailParam && emailParam != '') {
  document.getElementById('email').value = emailParam;
}

const form = document.querySelector('form');

form.addEventListener('submit', async (e) => {
  e.preventDefault();

  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  if (email === '' || password === '') {
    // TODO: show error
    return;
  }

  const body = {
    email,
    password,
  };

  try {
    const response = await fetch(`http://localhost:3000/api/sign-up`, {
      method: `POST`,
      body: JSON.stringify(body),
    });

    const data = await response.json();

    // TODO: send sessionId to extension
    document.location = '/passes';
  } catch (error) {
    console.log(error);
  }
});
