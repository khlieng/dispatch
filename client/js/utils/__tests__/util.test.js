import React from 'react';
import TestRenderer from 'react-test-renderer';
import {
  trimPrefixChar,
  isChannel,
  isValidNick,
  isValidChannel,
  isValidUsername
} from '..';
import linkify from '../linkify';

const render = el => TestRenderer.create(el).toJSON();

describe('trimPrefixChar()', () => {
  it('trims prefix characters', () => {
    expect(trimPrefixChar('##chan', '#')).toBe('chan');
    expect(trimPrefixChar('#chan', '#')).toBe('chan');
    expect(trimPrefixChar('chan', '#')).toBe('chan');
  });
});

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

describe('isValidNick()', () => {
  it('validates nicks', () =>
    Object.entries({
      bob: true,
      'bob likes cake': false,
      '-bob': false,
      'bob.': false,
      'bob-': true,
      '1bob': false,
      '[bob}': true,
      '': false,
      '   ': false
    }).forEach(([input, expected]) =>
      expect(isValidNick(input)).toBe(expected)
    ));
});

describe('isValidChannel()', () => {
  it('validates channels', () =>
    Object.entries({
      '#chan': true,
      '#cak e': false,
      '#cake:': false,
      '#[cake]': true,
      '#ca,ke': false,
      '': false,
      '   ': false,
      cake: false
    }).forEach(([input, expected]) =>
      expect(isValidChannel(input)).toBe(expected)
    ));

  it('handles requirePrefix', () =>
    Object.entries({
      chan: true,
      'cak e': false,
      '#cake:': false,
      '#[cake]': true,
      '#ca,ke': false
    }).forEach(([input, expected]) =>
      expect(isValidChannel(input, false)).toBe(expected)
    ));
});

describe('isValidUsername()', () => {
  it('validates usernames', () =>
    Object.entries({
      bob: true,
      'bob likes cake': false,
      '-bob': true,
      'bob.': true,
      'bob-': true,
      '1bob': true,
      '[bob}': true,
      '': false,
      '   ': false,
      'b@b': false
    }).forEach(([input, expected]) =>
      expect(isValidUsername(input)).toBe(expected)
    ));
});

describe('linkify()', () => {
  const proto = href => (href.indexOf('http') !== 0 ? `http://${href}` : href);
  const linkTo = href =>
    render(
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
      expect(render(linkify(input))).toEqual(expected)
    ));
});
