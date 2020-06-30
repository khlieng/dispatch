import Autolinker from 'autolinker';

const autolinker = new Autolinker({
  stripPrefix: false,
  stripTrailingSlash: false
});

function pushText(arr, text) {
  const last = arr[arr.length - 1];
  if (last?.type === 'text') {
    last.text += text;
  } else {
    arr.push({
      type: 'text',
      text
    });
  }
}

function pushLink(arr, url, text) {
  arr.push({
    type: 'link',
    url,
    text
  });
}

export default function linkify(text) {
  if (typeof text !== 'string') {
    return text;
  }

  let matches = autolinker.parseText(text);

  if (matches.length === 0) {
    return [
      {
        type: 'text',
        text
      }
    ];
  }

  const result = [];
  let pos = 0;
  matches = autolinker.compactMatches(matches);

  for (let i = 0; i < matches.length; i++) {
    const match = matches[i];

    if (match.getType() === 'url') {
      if (match.offset > pos) {
        pushText(result, text.slice(pos, match.offset));
      }

      pushLink(result, match.getAnchorHref(), match.matchedText);
    } else {
      pushText(
        result,
        text.slice(pos, match.offset + match.matchedText.length)
      );
    }

    pos = match.offset + match.matchedText.length;
  }

  if (pos < text.length) {
    if (result[result.length - 1]?.type === 'text') {
      result[result.length - 1].text += text.slice(pos);
    } else {
      pushText(result, text.slice(pos));
    }
  }

  return result;
}
