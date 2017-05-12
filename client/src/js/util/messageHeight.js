import FontFaceObserver from 'fontfaceobserver';
import { stringWidth, measureScrollBarWidth } from './index';
import { updateMessageHeight } from '../actions/message';

const lineHeight = 24;
const menuWidth = 200;
const userListWidth = 200;
const messagePadding = 30;
const smallScreen = 600;
let windowWidth;

function init(store, charWidth, done) {
  window.messageIndent = 6 * charWidth;
  const scrollBarWidth = measureScrollBarWidth();
  let prevWrapWidth;

  function updateWidth() {
    windowWidth = window.innerWidth;
    let wrapWidth = windowWidth - scrollBarWidth - messagePadding;
    if (windowWidth > smallScreen) {
      wrapWidth -= menuWidth;
    }

    if (wrapWidth !== prevWrapWidth) {
      prevWrapWidth = wrapWidth;
      store.dispatch(updateMessageHeight(wrapWidth, charWidth));
    }
  }

  let resizeRAF;

  function resize() {
    if (resizeRAF) {
      window.cancelAnimationFrame(resizeRAF);
    }
    resizeRAF = window.requestAnimationFrame(updateWidth);
  }

  updateWidth();
  done();
  window.addEventListener('resize', resize);
}

export function initWidthUpdates(store, done) {
  let charWidth = localStorage.charWidth;
  if (charWidth) {
    init(store, parseFloat(charWidth), done);
  }

  new FontFaceObserver('Roboto Mono').load().then(() => {
    if (!charWidth) {
      charWidth = stringWidth(' ', '16px Roboto Mono');
      init(store, charWidth, done);
      localStorage.charWidth = charWidth;
    }
  });

  new FontFaceObserver('Montserrat').load();
  new FontFaceObserver('Montserrat', { weight: 700 }).load();
  new FontFaceObserver('Roboto Mono', { weight: 700 }).load();
}

export function findBreakpoints(text) {
  const breakpoints = [];

  for (let i = 0; i < text.length; i++) {
    const char = text.charAt(i);

    if (char === ' ') {
      breakpoints.push({ end: i, next: i + 1 });
    } else if (char === '-' && i !== text.length - 1) {
      breakpoints.push({ end: i + 1, next: i + 1 });
    }
  }

  return breakpoints;
}

export function messageHeight(message, width, charWidth, indent = 0) {
  let pad = (6 + (message.from ? message.from.length + 1 : 0)) * charWidth;
  let height = lineHeight + 8;

  if (message.channel && windowWidth > smallScreen) {
    width -= userListWidth;
  }

  if (pad + (message.length * charWidth) < width) {
    return height;
  }

  const breaks = message.breakpoints;
  let prevBreak = 0;
  let prevPos = 0;

  for (let i = 0; i < breaks.length; i++) {
    if (pad + ((breaks[i].end - prevBreak) * charWidth) >= width) {
      prevBreak = prevPos;
      pad = indent;
      height += lineHeight;
    }

    prevPos = breaks[i].next;
  }

  if (pad + ((message.length - prevBreak) * charWidth) >= width) {
    height += lineHeight;
  }

  return height;
}
