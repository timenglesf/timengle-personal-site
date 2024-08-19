document.documentElement.setAttribute('data-theme', 'mytheme');




document.addEventListener('DOMContentLoaded', function() {
  // Post form validation
  const titleInput = document.getElementById('title_input');
  const contentInput = document.querySelector("textarea")
  const titleWarning = document.getElementById('title_warning');
  const contentWarning = document.getElementById('content_warning');

  if (titleWarning) {
    console.log(titleInput)
    titleInput.focus()
  }

  if (contentWarning) {
    contentInput.focus();
  }

  // Sign up / Log in form validation

  const emailInput = document.getElementById('email');
  const emailConfirmInput = document.getElementById('confirm_email');
  const displayNameInput = document.getElementById('display_name');
  const displayNameWarning = document.getElementById('display_name_warning');
  const passwordInput = document.getElementById('password');
  const passwordConfirmInput = document.getElementById('confirm_password');
  const emailWarning = document.getElementById('email_warning');
  const emailConfirmWarning = document.getElementById('confirm_email_warning');
  const passwordWarning = document.getElementById('password_warning');
  const passwordConfirmWarning = document.getElementById('confirm_password_warning');


  if (emailWarning) {
    emailInput.focus();
  } else if (emailConfirmWarning) {
    emailConfirmInput.focus();
  } else if (displayNameWarning) {
    displayNameInput.focus();
  } else if (passwordWarning) {
    passwordInput.focus();
  } else if (passwordConfirmWarning) {
    passwordInput.focus();
  }

  // Prevent form submission on Enter key for title input
  if (titleInput) {
    titleInput.addEventListener("keydown", e => {
      if (e.key === "Enter") {
        e.preventDefault();
        contentInput.focus();
      }
    })
  }

});



