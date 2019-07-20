DOM = document.querySelector('.history-link');
if (DOM) {
  DOM.addEventListener('click', () => chrome.runtime.openOptionsPage());
}
