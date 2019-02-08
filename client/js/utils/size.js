let width;
let height;
const listeners = [];

function update() {
  width = window.innerWidth;
  height = window.innerHeight;

  for (let i = 0; i < listeners.length; i++) {
    listeners[i](width, height);
  }
}

let resizeRAF;

function resize() {
  if (resizeRAF) {
    window.cancelAnimationFrame(resizeRAF);
  }
  resizeRAF = window.requestAnimationFrame(update);
}

update();
window.addEventListener('resize', resize);

export function windowWidth() {
  return width;
}

export function windowHeight() {
  return height;
}

export function addResizeListener(f, init) {
  listeners.push(f);
  if (init) {
    f(width, height);
  }
}

export function removeResizeListener(f) {
  const i = listeners.indexOf(f);
  if (i > -1) {
    listeners.splice(i, 1);
  }
}
