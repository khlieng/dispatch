export default function wrapMessages(messages, width, charWidth, indent = 0) {
  return messages.withMutations(m => {
    for (let j = 0, llen = messages.size; j < llen; j++) {
      const message = messages.get(j);
      let lineWidth = (6 + (message.from ? message.from.length + 1 : 0)) * charWidth;

      if (lineWidth + message.message.length * charWidth < width) {
        m.setIn([j, 'lines'], [message.message]);
        continue;
      }

      const words = message.message.split(' ');
      const wrapped = [];
      let line = '';
      let wordCount = 0;
      let hasWrapped = false;

      // Add empty line if first word after timestamp + sender wraps
      if (words.length > 0 && message.from && lineWidth + words[0].length * charWidth >= width) {
        wrapped.push(line);
        lineWidth = 0;
      }

      for (let i = 0, wlen = words.length; i < wlen; i++) {
        const word = words[i];

        if (hasWrapped) {
          hasWrapped = false;
          lineWidth += indent;
        }

        lineWidth += word.length * charWidth;
        wordCount++;

        if (lineWidth >= width) {
          if (wordCount !== 1) {
            wrapped.push(line);

            if (i !== wlen - 1) {
              line = word + ' ';
              lineWidth = (word.length + 1) * charWidth;
              wordCount = 1;
            } else {
              wrapped.push(word);
            }
          } else {
            wrapped.push(word);
            lineWidth = 0;
            wordCount = 0;
          }

          hasWrapped = true;
        } else if (i !== wlen - 1) {
          line += word + ' ';
          lineWidth += charWidth;
        } else {
          line += word;
          wrapped.push(line);
        }
      }

      m.setIn([j, 'lines'], wrapped);
    }
  });
}
