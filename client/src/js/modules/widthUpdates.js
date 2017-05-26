import { when } from '../util/observe';
import { measureScrollBarWidth } from '../util';
import { getCharWidth } from '../state/environment';
import { updateMessageHeight } from '../state/messages';

const menuWidth = 200;
const messagePadding = 30;
const smallScreen = 600;

export default function widthUpdates({ store }) {
  when(store, getCharWidth, charWidth => {
    window.messageIndent = 6 * charWidth;
    const scrollBarWidth = measureScrollBarWidth();
    let prevWrapWidth;

    function updateWidth() {
      const windowWidth = window.innerWidth;
      let wrapWidth = windowWidth - scrollBarWidth - messagePadding;
      if (windowWidth > smallScreen) {
        wrapWidth -= menuWidth;
      }

      if (wrapWidth !== prevWrapWidth) {
        prevWrapWidth = wrapWidth;
        store.dispatch(updateMessageHeight(wrapWidth, charWidth, windowWidth));
      }
    }

    let resizeRAF;

    function resize() {
      if (resizeRAF) {
        window.cancelAnimationFrame(resizeRAF);
      }
      resizeRAF = window.requestAnimationFrame(updateWidth);
    }

    updateWidth();
    window.addEventListener('resize', resize);
  });
}
