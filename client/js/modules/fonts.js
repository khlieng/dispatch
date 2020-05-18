import { setCharWidth } from 'state/app';
import { stringWidth } from 'utils';

export default async function fonts({ store }) {
  let { charWidth } = localStorage;
  if (charWidth) {
    store.dispatch(setCharWidth(parseFloat(charWidth)));
  } else {
    await document.fonts.load('16px Roboto Mono');

    charWidth = stringWidth(' ', '16px Roboto Mono');
    store.dispatch(setCharWidth(charWidth));
    localStorage.charWidth = charWidth;
  }
}
