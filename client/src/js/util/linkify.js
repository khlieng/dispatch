import Autolinker from 'autolinker';
import React from 'react';

const autolinker = new Autolinker({
  stripPrefix: false,
  doJoin: false,
  replaceFn: (linker, match) => {
    if (match.getType() === 'url') {
      return <a target="_blank" rel="noopener noreferrer" href={match.getAnchorHref()}>{match.getAnchorText()}</a>;
    }

    return null;
  },
  React
});

export default function linkify(text) {
  return autolinker.link(text);
}
