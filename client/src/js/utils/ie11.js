if (Object.keys) {
  try {
    Object.keys('');
  } catch (e) {
    Object.keys = function keys(o, k, r) {
      r = [];
      // eslint-disable-next-line
      for (k in o) r.hasOwnProperty.call(o, k) && r.push(k);
      return r;
    };
  }
}
