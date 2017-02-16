const lineHeight = 24;

export default function messageHeight(message, width, charWidth, indent = 0) {
  let pad = (6 + (message.from ? message.from.length + 1 : 0)) * charWidth;
  let height = lineHeight + 4;

  if (message.channel) {
    width -= 200;
  }

  if (pad + (message.message.length * charWidth) < width) {
    return height;
  }

  let prevBreak = 0;
  let prevPos = 0;

  for (let i = 0, len = message.message.length; i < len; i++) {
    const c = message.message.charAt(i);

    if (c === ' ' || c === '-') {
      const end = c === ' ' ? i : i + 1;

      if (pad + ((end - prevBreak) * charWidth) >= width) {
        prevBreak = prevPos;
        pad = indent;
        height += lineHeight;
      }

      prevPos = i + 1;
    } else if (i === len - 1) {
      if (pad + ((len - prevBreak) * charWidth) >= width) {
        height += lineHeight;
      }
    }
  }

  return height;
}
