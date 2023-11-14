


const SERVER_URL = 'http://ip-10-113-60-162:8080';

const RANCHER_URL = 'https://ip-10-113-61-179/dashboard/home';


document.addEventListener('DOMContentLoaded', (event) => {
  checkUserExists();
});

function openLink() {
  window.open(RANCHER_URL, '_blank'); // '_blank' is used to open the link in a new tab
}

function startTimer(expirationTime){
    // User exists, hide login button and show delete button
    document.getElementById('loginButton').style.display = 'none';
    document.getElementById('closeSession').style.display = 'block';
    document.getElementById('closeSession').onclick = deleteUser; // Attach event handler
    // Set up the countdown timer
    startCountdown(expirationTime);
}
function logoutreset() {
  const countdownTimer = document.getElementById('countdownTimer');
  countdownTimer.style.display = 'none'; 
  document.getElementById('loginButton').style.display = 'block';
  document.getElementById('closeSession').style.display = 'none';
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
        bindLink(data.clusterId);
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
  // Display the provisioning popup with the loading animation
  showProvisioningPopup('Deleting...'); // You can pass a custom message if you want

  fetch(SERVER_URL+'/api/v1/user', {
    method: 'DELETE',
    headers: {
      'UserDN': 'wawrig2'
    },
  })
  .then(response => {
    if (response.ok) {
      // Wait for 5 seconds, then remove the provisioning popup and reset the UI
      setTimeout(() => {
        logoutreset();
        // Optionally, you can also reload the page here or redirect the user
      }, 5000);
    } else {
      throw new Error('Failed to delete user');
    }
  })
  .catch(error => {
    console.error('Error during deletion:', error);
    // Handle deletion error
  });
}

// Add event listener to login button
document.getElementById('loginButton').addEventListener('click', doLogin);



function doLogin(){
  // Display the provisioning popup
  showProvisioningPopup();
  sendLoginRequest()
  sendProjectRequest()
}

function sendProjectRequest(){
// Send the POST request
fetch(SERVER_URL + '/api/v1/project', {
  method: 'POST',
  headers: {
    'UserDN': 'wawrig2'
  },
})
.then(response => {
  if (response.status === 201) {
    return response.json();
  } else {
    throw new Error('Project was not successful');
  }
})
.catch(error => {
  console.error('There was an error creating project', error);
});
}

function sendLoginRequest() {
  // Send the POST request
  fetch(SERVER_URL + '/api/v1/user', {
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
    // Check if the user/project exists after 5 seconds
    setTimeout(() => {
      if (data && data.userId !== "") {
        startTimer(data.expirationTime);
        // bindLink(data.clusterId);
      }
    }, 5000);
  })
  .catch(error => {
    console.error('There was an error logging in:', error);
    // Handle login errors, possibly display a message to the user
  });
}


 function showProvisioningPopup(message) {
  // Create backdrop
  let backdrop = document.createElement('div');
  backdrop.setAttribute('id', 'provisioningBackdrop');
  backdrop.style.cssText = `
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(5px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 10002;`;

  // Create loading circle div
  let loadingCircle = document.createElement('div');
  loadingCircle.setAttribute('class', 'loading-circle');

  // Create provisioning text
  let provisioningText = document.createElement('p');
  provisioningText.textContent = message || 'Provisioning...'; // Default message is 'Provisioning...'
  provisioningText.style.color = 'white';
  provisioningText.style.textAlign = 'center';
  provisioningText.style.marginTop = '10px';
  provisioningText.style.marginLeft = '10px';

  // Append loading circle and text to backdrop
  backdrop.appendChild(loadingCircle);
  backdrop.appendChild(provisioningText);

  // Add backdrop to body
  document.body.appendChild(backdrop);

  // Remove the provisioning popup after 5 seconds
  setTimeout(() => {
    document.body.removeChild(backdrop);
  }, 5000);
}