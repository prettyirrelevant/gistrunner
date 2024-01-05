console.log('[gistrunner]:background.js loaded.');

function getDefaultSettings() {
  return { baseServerURL: 'https://gistrunner.fly.dev' };
}

function initializeExtension() {
  chrome.storage.local.get(['baseServerURL']).then(({ baseServerURL }) => {
    const defaultSettings = getDefaultSettings();

    defaultSettings.baseServerURL = baseServerURL || defaultSettings.baseServerURL;
    chrome.storage.local
      .set(defaultSettings)
      .then(() => console.log('[gistrunner]: default values for extension set!'));
  });
}

chrome.runtime.onInstalled.addListener(() => {
  initializeExtension();
});

chrome.runtime.onStartup.addListener(() => {
  initializeExtension();
});
