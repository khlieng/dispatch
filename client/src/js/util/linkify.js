import Autolinker from 'autolinker';
import React from 'react';

const autolinker = new Autolinker({
  stripPrefix: false,
  stripTrailingSlash: false
});

export default function linkify(text) {
  const matches = autolinker.parseText(text);
  const result = [];
  let pos = 0;

  for (let i = 0; i < matches.length; i++) {
    const match = matches[i];

    if (match.offset > pos) {
      result.push(text.slice(pos, match.offset));
      pos = match.offset;
    }

    if (match.getType() === 'url') {
      result.push(
        <a target="_blank" rel="noopener noreferrer" href={match.getAnchorHref()}>
          {match.matchedText}
        </a>
      );
    } else {
      result.push(match.matchedText);
    }

    pos += match.matchedText.length;
  }

  if (pos < text.length) {
    result.push(text.slice(pos));
  }

  return result;
}
