/* eslint-disable no-underscore-dangle */

// This entrypoint gets inlined in the index page cached by service workers
// and is responsible for fetching the data we would otherwise embed

window.__env__ = fetch(`/data${window.location.pathname}`, {
  credentials: 'same-origin'
}).then(res => {
  if (res.ok) {
    return res.json();
  }

  throw new Error(res.statusText);
});
