


const SERVER_URL = 'http://ip-10-113-60-162:8080';

document.addEventListener('DOMContentLoaded', (event) => {
  checkUserExists();
});

function startTimer(expirationTime){
    // User exists, hide login button and show delete button
    document.getElementById('loginButton').style.display = 'none';
    document.getElementById('deleteProject').style.display = 'block';
    document.getElementById('deleteProject').onclick = deleteUser; // Attach event handler
    // Set up the countdown timer
    startCountdown(expirationTime);
}


function checkUserExists() {
  fetch(SERVER_URL+'/api/v1/user',{
    method: 'GET',
    headers: {
      'UserDN': 'wawrig2'
      },
  })
    .then(response => {
      if (response.status === 401) {
        // Handle 401 Unauthorized
        showUnauthorizedPopup();
        throw new Error('Unauthorized');
      }
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      return response.json();
    })
    .then(data => {
      if (data && data.userId !== "") {
        startTimer(data.expirationTime)
      }
    })
    .catch(error => {
      console.error('Error checking user:', error);
    });
}


function showUnauthorizedPopup() {
  // Disable all buttons on the page
  document.querySelectorAll('button').forEach(button => {
    button.disabled = true;
  });

  // Create backdrop
  let backdrop = document.createElement('div');
  backdrop.style.position = 'fixed';
  backdrop.style.top = '0';
  backdrop.style.left = '0';
  backdrop.style.width = '100%';
  backdrop.style.height = '100%';
  backdrop.style.backgroundColor = 'rgba(0, 0, 0, 0.5)';
  backdrop.style.backdropFilter = 'blur(5px)';

  // Create popup div
  let popup = document.createElement('div');
  popup.setAttribute('id', 'unauthorizedPopup');
  popup.style.position = 'absolute';
  popup.style.left = '50%';
  popup.style.top = '50%';
  popup.style.transform = 'translate(-50%, -50%)';
  popup.style.padding = '20px';
  popup.style.backgroundColor = 'white';
  popup.style.border = '1px solid black';
  popup.style.zIndex = '10001'; // Ensure popup is above backdrop

  // Create popup text
  let popupText = document.createElement('p');
  popupText.textContent = 'Unable To Authenticate With PKI';
  
  // Create reload button
  let reloadButton = document.createElement('button');
  reloadButton.textContent = 'Reload Page';
  reloadButton.style.marginTop = '20px'; // Add some space above the button
  reloadButton.style.padding = '10px 20px'; // Make the button bigger
  reloadButton.onclick = function() {
    // Re-enable all buttons before reloading
    document.querySelectorAll('button').forEach(button => {
      button.disabled = false;
    });
    window.location.reload();
  };

  // Center button inside the popup
  let buttonContainer = document.createElement('div');
  buttonContainer.style.textAlign = 'center'; // Center-align the button container
  buttonContainer.appendChild(reloadButton);

  // Append text and button container to popup
  popup.appendChild(popupText);
  popup.appendChild(buttonContainer);

  // Add backdrop and popup to body
  document.body.appendChild(backdrop);
  document.body.appendChild(popup);
}


function startCountdown(expirationTime) {
  const countdownTimer = document.getElementById('countdownTimer');
  countdownTimer.style.display = 'block'; // Show the countdown timer

  const countDownDate = new Date(expirationTime).getTime();

  const x = setInterval(() => {
    let now = new Date().getTime();
    let distance = countDownDate - now;

    let hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    let minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    let seconds = Math.floor((distance % (1000 * 60)) / 1000);

    countdownTimer.innerHTML = hours + "h " + minutes + "m " + seconds + "s ";

    if (distance < 0) {
      clearInterval(x);
      countdownTimer.innerHTML = "EXPIRED";
    }
  }, 1000);
}

function deleteUser() {
  fetch(SERVER_URL+'/api/v1/user',{
    method: 'DELETE',
    headers: {
      'UserDN': 'wawrig2'
      },
  })
    .then(response => {
      if (!response.ok) {
        throw new Error('Error deleting user');
      }
      return response.json();
    })
    .then(data => {
      console.log('User deleted:', data);
      // Handle successful deletion, maybe redirect to a login page
    })
    .catch(error => {
      console.error('Error during deletion:', error);
      // Handle deletion error
    });
}

// Add event listener to login button
document.getElementById('loginButton').addEventListener('click', sendLoginRequest);

function sendLoginRequest() {
    // Send the POST request
    fetch(SERVER_URL+'/api/v1/user', {
      method: 'POST',
      headers: {
        'UserDN': 'wawrig2'
        },
    })
    .then(response => {
      if (response.status === 201) {

        return response.json();
      } else {
        throw new Error('Login was not successful');
      }
    })
    .then(data => {
              // Login was successful, now check if the user/project exists
              if (data && data.userId !== "") {
                startTimer(data.expirationTime)
              }
    })
    .catch(error => {
      console.error('There was an error logging in:', error);
      // Handle login errors, possibly display a message to the user
    });
  }
