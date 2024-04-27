const btn = document.querySelector('.btn-toggle');
const prefersDarkScheme = window.matchMedia('(prefers-color-scheme: dark)');
if (prefersDarkScheme.matches) {
  document.body.classList.add('dark-theme');
}

// Listen for a click on the button
btn.addEventListener('click', function () {
  // Then toggle (add/remove) the .dark-theme class to the body
  document.body.classList.toggle('dark-theme');
});
