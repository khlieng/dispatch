/* eslint-disable no-bitwise */
import { hsluvToHex } from 'hsluv';

//
// github.com/sindresorhus/fnv1a
//
const OFFSET_BASIS_32 = 2166136261;

const fnv1a = string => {
  let hash = OFFSET_BASIS_32;

  for (let i = 0; i < string.length; i++) {
    hash ^= string.charCodeAt(i);

    // 32-bit FNV prime: 2**24 + 2**8 + 0x93 = 16777619
    // Using bitshift for accuracy and performance. Numbers in JS suck.
    hash +=
      (hash << 1) + (hash << 4) + (hash << 7) + (hash << 8) + (hash << 24);
  }

  return hash >>> 0;
};

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
