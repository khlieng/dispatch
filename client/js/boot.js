/* eslint-disable no-underscore-dangle */

window.__init__ = fetch('/init', {
  credentials: 'same-origin'
}).then(res => {
  if (res.ok) {
    return res.json();
  }

  throw new Error(res.statusText);
});
