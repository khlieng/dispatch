import FontFaceObserver from 'fontfaceobserver';
import { setCharWidth } from 'state/app';
import { stringWidth } from 'util';

export default function fonts({ store }) {
  let charWidth = localStorage.charWidth;
  if (charWidth) {
    store.dispatch(setCharWidth(parseFloat(charWidth)));
  }

  new FontFaceObserver('Roboto Mono').load().then(() => {
    if (!charWidth) {
      charWidth = stringWidth(' ', '16px Roboto Mono');
      store.dispatch(setCharWidth(charWidth));
      localStorage.charWidth = charWidth;
    }
  });

  new FontFaceObserver('Montserrat').load();
  new FontFaceObserver('Montserrat', { weight: 700 }).load();
  new FontFaceObserver('Roboto Mono', { weight: 700 }).load();
}
