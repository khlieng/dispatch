import padStart from 'lodash/padStart';

export { findBreakpoints, messageHeight } from './messageHeight';
export { default as linkify } from './linkify';

export function normalizeChannel(channel) {
  if (channel.indexOf('#') !== 0) {
    return channel;
  }

  return channel.split('#').join('').toLowerCase();
}

export function isChannel(name) {
  // TODO: Handle other channel types
  return typeof name === 'string' && name[0] === '#';
}

export function timestamp(date = new Date()) {
  const h = padStart(date.getHours(), 2, '0');
  const m = padStart(date.getMinutes(), 2, '0');
  return `${h}:${m}`;
}

const canvas = document.createElement('canvas');
const ctx = canvas.getContext('2d');

export function stringWidth(str, font) {
  ctx.font = font;
  return ctx.measureText(str).width;
}

export function measureScrollBarWidth() {
  const outer = document.createElement('div');
  outer.style.visibility = 'hidden';
  outer.style.width = '100px';

  document.body.appendChild(outer);

  const widthNoScroll = outer.offsetWidth;
  outer.style.overflow = 'scroll';

  const inner = document.createElement('div');
  inner.style.width = '100%';
  outer.appendChild(inner);

  const widthWithScroll = inner.offsetWidth;

  outer.parentNode.removeChild(outer);

  return widthNoScroll - widthWithScroll;
}

export function find(arr, pred) {
  if (!arr) {
    return null;
  }

  for (let i = 0; i < arr.length; i++) {
    if (pred(arr[i])) {
      return arr[i];
    }
  }
}
