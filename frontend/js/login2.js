// Constants
const SERVER_URL = 'http://ip-10-113-60-162:8080';
const RANCHER_URL = 'https://ip-10-113-61-179/dashboard/home';

// Elements
const loginButton = document.getElementById('loginButton');
const closeSessionButton = document.getElementById('closeSession');
const countdownTimer = document.getElementById('countdownTimer');

// Event Listeners
document.addEventListener('DOMContentLoaded', checkUserExists);
loginButton.addEventListener('click', doLogin);
closeSessionButton.addEventListener('click', deleteUser);

// Open external link
function openLink() {
  window.open(RANCHER_URL, '_blank');
}

// User login/logout display logic
function updateUIForLoggedInUser(expirationTime) {
  toggleElement(loginButton, false);
  toggleElement(closeSessionButton, true);
  startCountdown(expirationTime);
}

function resetLogoutState() {
  toggleElement(countdownTimer, false);
  toggleElement(loginButton, true);
  toggleElement(closeSessionButton, false);
}

function toggleElement(element, show) {
  element.style.display = show ? 'block' : 'none';
}

// Check user existence
function checkUserExists() {
  fetch(`${SERVER_URL}/api/v1/user`, { method: 'GET', headers: { 'UserDN': 'wawrig2' } })
    .then(handleResponse)
    .then(data => data && data.userId !== "" ? updateUIForLoggedInUser(data.expirationTime) : null)
    .catch(error => {
        console.error('Error checking user:', error);
        showMaintenanceModePopup(); // Show popup on network error
      });
}

// Handle fetch response
function handleResponse(response) {
  if (!response.ok) {
    if (response.status === 401) showUnauthorizedPopup();
    throw new Error('Network response was not ok');
  }
  return response.json();
}

// Show popup
function showPopup(content, id, customStyle = '') {
  // Common popup creation logic
  const popup = document.createElement('div');
  popup.id = id;
  popup.style.cssText = `position: absolute; left: 50%; top: 50%; transform: translate(-50%, -50%); padding: 20px; background-color: white; border: 1px solid black; z-index: 10001; ${customStyle}`;
  popup.innerHTML = content;
  document.body.appendChild(popup);
  return popup;
}

// Unauthorized popup
function showUnauthorizedPopup() {
  const content = '<p>Unable To Authenticate With PKI</p><button onclick="reloadPage()">Reload Page</button>';
  const popup = showPopup(content, 'unauthorizedPopup');
  disableButtons();
  centerPopupButton(popup);
}

function showServiceUnavailablePopup() {
    // Create backdrop
    const backdrop = document.createElement('div');
    backdrop.classList.add('backdrop');

    // Create popup
    const popup = document.createElement('div');
    popup.classList.add('popup');
    popup.innerHTML = `
        <h2>Play With CK8S Is Undergoing Matinance</h2>
        <p>The service is currently unavailable. Please try again later.</p>
        <button onclick="window.location.reload();">Reload Page</button>
    `;

    backdrop.appendChild(popup);
    document.body.appendChild(backdrop);
}

function showMaintenanceModePopup() {
    // Create backdrop
    const backdrop = document.createElement('div');
    backdrop.classList.add('backdrop');

    // Create popup
    const popup = document.createElement('div');
    popup.classList.add('popup');
    popup.innerHTML = `
        <h2>PWCK8S Maintenance Mode</h2>
        <p>We're currently performing scheduled maintenance.</p>
        <p>Please check back later.</p>
        <button onclick="window.location.reload();">Reload Page</button>
    `;

    backdrop.appendChild(popup);
    document.body.appendChild(backdrop);
}

function reloadPage() {
  enableButtons();
  window.location.reload();
}

// Disable/Enable all buttons
function disableButtons() {
  document.querySelectorAll('button').forEach(button => button.disabled = true);
}

function enableButtons() {
  document.querySelectorAll('button').forEach(button => button.disabled = false);
}

// Center button in popup
function centerPopupButton(popup) {
  const buttonContainer = popup.querySelector('button').parentNode;
  buttonContainer.style.textAlign = 'center';
}

// Countdown timer
function startCountdown(expirationTime) {
  toggleElement(countdownTimer, true);
  const countDownDate = new Date(expirationTime).getTime();
  const intervalId = setInterval(() => {
    const now = new Date().getTime();
    const distance = countDownDate - now;
    if (distance < 0) {
      clearInterval(intervalId);
      countdownTimer.innerHTML = "EXPIRED";
      return;
    }
    countdownTimer.innerHTML = formatTime(distance);
  }, 1000);
}

function formatTime(distance) {
  const hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((distance % (1000 * 60)) / 1000);
  return `${hours}h ${minutes}m ${seconds}s`;
}

// User deletion
function deleteUser() {
  showProvisioningPopup('Deleting...');
  fetch(`${SERVER_URL}/api/v1/user`, { method: 'DELETE', headers: { 'UserDN': 'wawrig2' } })
    .then(response => response.ok ? setTimeout(resetLogoutState, 5000) : Promise.reject(new Error('Failed to delete user')))
    .catch(console.error);
}

// Login logic
function doLogin() {
  showProvisioningPopup();
  sendLoginRequest().then(data => {
    if (data && data.userId !== "") setTimeout(() => updateUIForLoggedInUser(data.expirationTime), 5000);
  }).catch(console.error);
  sendProjectRequest().catch(console.error);
}

function sendLoginRequest() {
  return fetch(`${SERVER_URL}/api/v1/user`, { method: 'POST', headers: { 'UserDN': 'wawrig2' } })
    .then(handleResponse);
}

function sendProjectRequest() {
  return fetch(`${SERVER_URL}/api/v1/project`, { method: 'POST', headers: { 'UserDN': 'wawrig2' } })
    .then(response => {
      if (response.status !== 201) throw new Error('Project was not successful');
      return response.json();
    });
}

// Provisioning popup
function showProvisioningPopup(message = 'Provisioning...') {
  const content = `<div class="loading-circle"></div><p style="color: white; text-align: center; margin-top: 10px;">${message}</p>`;
  showPopup(content, 'provisioningBackdrop', 'display: flex; justify-content: center; align-items: center;');
  setTimeout(() => document.body.removeChild(document.getElementById('provisioningBackdrop')), 5000);
}
