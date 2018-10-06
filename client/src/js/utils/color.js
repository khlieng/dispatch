import fnv1a from '@sindresorhus/fnv1a';

const colors = [];

for (let i = 0; i < 72; i++) {
  colors[i] = `hsl(${i * 5}, 30%, 40%)`;
  colors[i + 72] = `hsl(${i * 5}, 60%, 40%)`;
  colors[i + 144] = `hsl(${i * 5}, 90%, 40%)`;
}

const cache = {};

export default function stringToHSL(str) {
  if (cache[str]) {
    return cache[str];
  }

  const color = colors[fnv1a(str) % 216];
  cache[str] = color;
  return color;
}
