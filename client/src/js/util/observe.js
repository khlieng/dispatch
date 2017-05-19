function subscribe(store, selector, handler) {
  let prev = selector(store.getState());
  handler(prev);

  store.subscribe(() => {
    const next = selector(store.getState());
    if (next !== prev) {
      handler(next);
      prev = next;
    }
  });
}

function subscribeArray(store, selectors, handler) {
  let state = store.getState();
  let prev = selectors.map(selector => selector(state));
  handler(...prev);

  store.subscribe(() => {
    state = store.getState();
    const next = [];
    let changed = false;

    for (let i = 0; i < selectors.length; i++) {
      next[i] = selectors[i](state);
      if (next[i] !== prev[i]) {
        changed = true;
      }
    }

    if (changed) {
      handler(...next);
      prev = next;
    }
  });
}

export default function observe(store, selector, handler) {
  if (Array.isArray(selector)) {
    subscribeArray(store, selector, handler);
  } else {
    subscribe(store, selector, handler);
  }
}
