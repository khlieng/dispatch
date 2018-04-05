const positions = {};

export function getScrollPos(key) {
  if (key in positions) {
    return positions[key];
  }
  return -1;
}

export function saveScrollPos(key, pos) {
  positions[key] = pos;
}
