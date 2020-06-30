export const formatChars = {
  bold: 0x02,
  italic: 0x1d,
  underline: 0x1f,
  strikethrough: 0x1e,
  color: 0x03,
  reverseColor: 0x16,
  reset: 0x0f
};

export const colors = {
  0: 'white',
  1: 'black',
  2: 'blue',
  3: 'green',
  4: 'red',
  5: 'brown',
  6: 'magenta',
  7: 'orange',
  8: 'yellow',
  9: 'lightgreen',
  10: 'cyan',
  11: 'lightcyan',
  12: 'lightblue',
  13: 'pink',
  14: 'gray',
  15: 'lightgray',
  16: '#470000',
  17: '#472100',
  18: '#474700',
  19: '#324700',
  20: '#004700',
  21: '#00472c',
  22: '#004747',
  23: '#002747',
  24: '#000047',
  25: '#2e0047',
  26: '#470047',
  27: '#47002a',
  28: '#740000',
  29: '#743a00',
  30: '#747400',
  31: '#517400',
  32: '#007400',
  33: '#007449',
  34: '#007474',
  35: '#004074',
  36: '#000074',
  37: '#4b0074',
  38: '#740074',
  39: '#740045',
  40: '#b50000',
  41: '#b56300',
  42: '#b5b500',
  43: '#7db500',
  44: '#00b500',
  45: '#00b571',
  46: '#00b5b5',
  47: '#0063b5',
  48: '#0000b5',
  49: '#7500b5',
  50: '#b500b5',
  51: '#b5006b',
  52: '#ff0000',
  53: '#ff8c00',
  54: '#ffff00',
  55: '#b2ff00',
  56: '#00ff00',
  57: '#00ffa0',
  58: '#00ffff',
  59: '#008cff',
  60: '#0000ff',
  61: '#a500ff',
  62: '#ff00ff',
  63: '#ff0098',
  64: '#ff5959',
  65: '#ffb459',
  66: '#ffff71',
  67: '#cfff60',
  68: '#6fff6f',
  69: '#65ffc9',
  70: '#6dffff',
  71: '#59b4ff',
  72: '#5959ff',
  73: '#c459ff',
  74: '#ff66ff',
  75: '#ff59bc',
  76: '#ff9c9c',
  77: '#ffd39c',
  78: '#ffff9c',
  79: '#e2ff9c',
  80: '#9cff9c',
  81: '#9cffdb',
  82: '#9cffff',
  83: '#9cd3ff',
  84: '#9c9cff',
  85: '#dc9cff',
  86: '#ff9cff',
  87: '#ff94d3',
  88: '#000000',
  89: '#131313',
  90: '#282828',
  91: '#363636',
  92: '#4d4d4d',
  93: '#656565',
  94: '#818181',
  95: '#9f9f9f',
  96: '#bcbcbc',
  97: '#e2e2e2',
  98: '#ffffff'
};

function tokenize(str) {
  const tokens = [];

  let colorBuffer = '';
  let color = false;
  let background = false;
  let colorToken;

  let start = 0;
  let end = 0;

  const pushText = () => {
    if (end > start) {
      tokens.push({
        type: 'text',
        content: str.slice(start, end)
      });
      start = end;
    }
  };

  const pushToken = token => {
    pushText();
    tokens.push(token);
  };

  for (let i = 0; i < str.length; i++) {
    const charCode = str.charCodeAt(i);

    if (color) {
      if (charCode >= 48 && charCode <= 57 && colorBuffer.length < 2) {
        colorBuffer += str[i];
      } else if (charCode === 44 && !background) {
        colorToken.color = colors[parseInt(colorBuffer, 10)];
        colorBuffer = '';
        background = true;
      } else {
        if (background) {
          if (colorBuffer.length > 0) {
            colorToken.background = colors[parseInt(colorBuffer, 10)];
          } else {
            // Trailing comma
            start--;
          }
        } else {
          colorToken.color = colors[parseInt(colorBuffer, 10)];
        }

        start--;
        colorBuffer = '';
        color = false;
        tokens.push(colorToken);
      }
    } else {
      switch (charCode) {
        case formatChars.bold:
          pushToken({
            type: 'bold'
          });
          break;

        case formatChars.italic:
          pushToken({
            type: 'italic'
          });
          break;

        case formatChars.underline:
          pushToken({
            type: 'underline'
          });
          break;

        case formatChars.strikethrough:
          pushToken({
            type: 'strikethrough'
          });
          break;

        case formatChars.color:
          pushText();

          colorToken = {
            type: 'color'
          };
          color = true;
          background = false;
          break;

        case formatChars.reverseColor:
          pushToken({
            type: 'reverse'
          });
          break;

        case formatChars.reset:
          pushToken({
            type: 'reset'
          });
          break;

        default:
          start--;
      }
    }

    start++;
    end++;
  }

  if (start === 0) {
    return str;
  }

  pushText();

  return tokens;
}

function colorifyString(str, state = {}) {
  const tokens = tokenize(str);

  if (tokens === str) {
    return [tokens, state];
  }

  const result = [];
  let style = state.style || {};
  let reverse = state.reverse || false;

  const toggle = (prop, value, multiple) => {
    if (style[prop]) {
      if (multiple) {
        const props = style[prop].split(' ');
        const i = props.indexOf(value);
        if (i !== -1) {
          props.splice(i, 1);
        } else {
          props.push(value);
        }
        style[prop] = props.join(' ');
      } else {
        delete style[prop];
      }
    } else {
      style[prop] = value;
    }
  };

  for (let i = 0; i < tokens.length; i++) {
    const token = tokens[i];

    switch (token.type) {
      case 'bold':
        toggle('fontWeight', 700);
        break;

      case 'italic':
        toggle('fontStyle', 'italic');
        break;

      case 'underline':
        toggle('textDecoration', 'underline', true);
        break;

      case 'strikethrough':
        toggle('textDecoration', 'line-through', true);
        break;

      case 'color':
        if (!token.color) {
          delete style.color;
          delete style.background;
        } else if (reverse) {
          style.color = token.background;
          style.background = token.color;
        } else {
          style.color = token.color;
          style.background = token.background;
        }
        break;

      case 'reverse':
        reverse = !reverse;
        if (style.color) {
          const bg = style.background;
          style.background = style.color;
          style.color = bg;
        }
        break;

      case 'reset':
        style = {};
        break;

      case 'text':
        if (Object.keys(style).length > 0) {
          result.push({
            type: 'format',
            style,
            text: token.content
          });
          style = { ...style };
        } else {
          result.push({
            type: 'text',
            text: token.content
          });
        }
        break;

      default:
    }
  }

  return [result, { style, reverse }];
}

export default function colorify(blocks) {
  if (!blocks) {
    return blocks;
  }

  const result = [];
  let colored;
  let state;

  for (let i = 0; i < blocks.length; i++) {
    const block = blocks[i];

    if (block.type === 'text') {
      [colored, state] = colorifyString(block.text, state);
      if (colored !== block.text) {
        result.push(...colored);
      } else {
        result.push(block);
      }
    } else {
      result.push(block);
    }
  }

  return result;
}
