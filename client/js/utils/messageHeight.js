const lineHeight = 24;
const userListWidth = 200;
const smallScreen = 600;

function findBreakpointsString(text, breakpoints, index) {
  for (let i = 0; i < text.length; i++) {
    const char = text.charAt(i);

    if (char === ' ') {
      breakpoints.push({ end: i + index, next: i + 1 + index });
    } else if (i !== text.length - 1 && (char === '-' || char === '?')) {
      breakpoints.push({ end: i + 1 + index, next: i + 1 + index });
    }
  }
}

export function findBreakpoints(text) {
  const breakpoints = [];
  let length = 0;

  if (typeof text === 'string') {
    findBreakpointsString(text, breakpoints, length);
    length = text.length;
  } else if (Array.isArray(text)) {
    for (let i = 0; i < text.length; i++) {
      const node = text[i];

      if (typeof node === 'string') {
        findBreakpointsString(node, breakpoints, length);
        length += node.length;
      } else {
        findBreakpointsString(node.props.children, breakpoints, length);
        length += node.props.children.length;
      }
    }
  }

  return [breakpoints, length];
}

export function messageHeight(
  message,
  wrapWidth,
  charWidth,
  indent = 0,
  windowWidth
) {
  let pad = (6 + (message.from ? message.from.length + 1 : 0)) * charWidth;
  let height = lineHeight + 8;

  if (message.channel && windowWidth > smallScreen) {
    wrapWidth -= userListWidth;
  }

  if (pad + message.length * charWidth < wrapWidth) {
    return height;
  }

  const breaks = message.breakpoints;
  let prevBreak = 0;
  let prevPos = 0;

  for (let i = 0; i < breaks.length; i++) {
    if (pad + (breaks[i].end - prevBreak) * charWidth >= wrapWidth) {
      prevBreak = prevPos;
      pad = indent;
      height += lineHeight;
    }

    prevPos = breaks[i].next;
  }

  if (pad + (message.length - prevBreak) * charWidth >= wrapWidth) {
    height += lineHeight;
  }

  return height;
}
