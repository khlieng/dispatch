import { hsluvToHex } from 'hsluv';
import fnv1a from '@sindresorhus/fnv1a';

const colors = [];

for (let i = 0; i < 72; i++) {
  colors[i] = hsluvToHex([i * 5, 40, 50]);
  colors[i + 72] = hsluvToHex([i * 5, 70, 50]);
  colors[i + 144] = hsluvToHex([i * 5, 100, 50]);
}

const cache = {};

export default function stringToRGB(str) {
  if (cache[str]) {
    return cache[str];
  }

  const color = colors[fnv1a(str) % colors.length];
  cache[str] = color;
  return color;
}
