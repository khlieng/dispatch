const lineHeight = 24;
let prevWidth;
let windowWidth;

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
  let height = lineHeight + 4;

  if (message.channel) {
    if (width !== prevWidth) {
      prevWidth = width;
      windowWidth = window.innerWidth;
    }

    if (windowWidth > 600) {
      width -= 200;
    }
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
