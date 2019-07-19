const DOM = document.querySelector('.setting-link');
if (DOM) {
  DOM.addEventListener('click', () => chrome.runtime.openOptionsPage());
}
