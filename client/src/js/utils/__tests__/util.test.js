import React from 'react';
import { isChannel } from '..';
import linkify from '../linkify';

describe('isChannel()', () => {
  it('it handles strings', () => {
    expect(isChannel('#cake')).toBe(true);
    expect(isChannel('cake')).toBe(false);
  });

  it('handles tab objects', () => {
    expect(isChannel({ name: '#cake' })).toBe(true);
    expect(isChannel({ name: 'cake' })).toBe(false);
  });
});

describe('linkify()', () => {
  const proto = href => (href.indexOf('http') !== 0 ? `http://${href}` : href);
  const linkTo = href => (
    <a href={proto(href)} rel="noopener noreferrer" target="_blank">
      {href}
    </a>
  );

  it('returns the arg when no matches are found', () =>
    [null, undefined, 10, false, true, 'just some text', ''].forEach(input =>
      expect(linkify(input)).toBe(input)
    ));

  it('linkifies text', () =>
    Object.entries({
      'google.com': linkTo('google.com'),
      'google.com stuff': [linkTo('google.com'), ' stuff'],
      'cake google.com stuff': ['cake ', linkTo('google.com'), ' stuff'],
      'cake google.com stuff https://google.com': [
        'cake ',
        linkTo('google.com'),
        ' stuff ',
        linkTo('https://google.com')
      ],
      'cake google.com stuff pie https://google.com  ': [
        'cake ',
        linkTo('google.com'),
        ' stuff pie ',
        linkTo('https://google.com'),
        '  '
      ],
      ' google.com': [' ', linkTo('google.com')],
      'google.com ': [linkTo('google.com'), ' '],
      '/google.com?': ['/', linkTo('google.com'), '?']
    }).forEach(([input, expected]) =>
      expect(linkify(input)).toEqual(expected)
    ));
});
