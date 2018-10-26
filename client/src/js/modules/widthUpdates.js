import { getCharWidth } from 'state/app';
import { updateMessageHeight } from 'state/messages';
import { when } from 'utils/observe';
import { measureScrollBarWidth } from 'utils';
import { addResizeListener } from 'utils/size';

const menuWidth = 200;
const messagePadding = 30;
const smallScreen = 600;

export default function widthUpdates({ store }) {
  when(store, getCharWidth, charWidth => {
    window.messageIndent = 6 * charWidth;
    const scrollBarWidth = measureScrollBarWidth();
    let prevWrapWidth;

    function updateWidth(windowWidth) {
      let wrapWidth = windowWidth - scrollBarWidth - messagePadding;
      if (windowWidth > smallScreen) {
        wrapWidth -= menuWidth;
      }

      if (wrapWidth !== prevWrapWidth) {
        prevWrapWidth = wrapWidth;
        store.dispatch(updateMessageHeight(wrapWidth, charWidth, windowWidth));
      }
    }

    addResizeListener(updateWidth, true);
  });
}
