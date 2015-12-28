import { padLeft } from 'lodash';

export wrapMessages from './wrapMessages';

export function normalizeChannel(channel) {
  if (channel.indexOf('#') !== 0) {
    return channel;
  }

  return channel.split('#').join('').toLowerCase();
}

export function createUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

export function timestamp(date = new Date()) {
  const h = padLeft(date.getHours(), 2, '0');
  const m = padLeft(date.getMinutes(), 2, '0');
  return h + ':' + m;
}

const canvas = document.createElement('canvas');
const ctx = canvas.getContext('2d');

export function stringWidth(str, font) {
  ctx.font = font;
  return ctx.measureText(str).width;
}

export function scrollbarWidth() {
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
