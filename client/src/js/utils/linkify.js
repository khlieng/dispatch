import Autolinker from 'autolinker';
import React from 'react';

const autolinker = new Autolinker({
  stripPrefix: false,
  stripTrailingSlash: false
});

export default function linkify(text) {
  let matches = autolinker.parseText(text);

  if (matches.length === 0) {
    return text;
  }

  const result = [];
  let pos = 0;
  matches = autolinker.compactMatches(matches);

  for (let i = 0; i < matches.length; i++) {
    const match = matches[i];

    if (match.getType() === 'url') {
      if (match.offset > pos) {
        if (typeof result[result.length - 1] === 'string') {
          result[result.length - 1] += text.slice(pos, match.offset);
        } else {
          result.push(text.slice(pos, match.offset));
        }
      }

      result.push(
        <a
          target="_blank"
          rel="noopener noreferrer"
          href={match.getAnchorHref()}
          key={i}
        >
          {match.matchedText}
        </a>
      );
    } else if (typeof result[result.length - 1] === 'string') {
      result[result.length - 1] += text.slice(
        pos,
        match.offset + match.matchedText.length
      );
    } else {
      result.push(text.slice(pos, match.offset + match.matchedText.length));
    }

    pos = match.offset + match.matchedText.length;
  }

  if (pos < text.length) {
    if (typeof result[result.length - 1] === 'string') {
      result[result.length - 1] += text.slice(pos);
    } else {
      result.push(text.slice(pos));
    }
  }

  if (result.length === 1) {
    return result[0];
  }

  return result;
}
