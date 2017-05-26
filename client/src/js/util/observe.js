function subscribeArray(store, selectors, handler, init) {
  let state = store.getState();
  let prev = selectors.map(selector => selector(state));
  if (init) {
    handler(...prev);
  }

  return store.subscribe(() => {
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

function subscribe(store, selector, handler, init) {
  if (Array.isArray(selector)) {
    return subscribeArray(store, selector, handler, init);
  }

  let prev = selector(store.getState());
  if (init) {
    handler(prev);
  }

  return store.subscribe(() => {
    const next = selector(store.getState());
    if (next !== prev) {
      handler(next);
      prev = next;
    }
  });
}

//
// Handler gets called every time the selector(s) change
//
export function observe(store, selector, handler) {
  return subscribe(store, selector, handler, true);
}

//
// Handler gets called once the next time the selector(s) change
//
export function once(store, selector, handler) {
  let done = false;
  const unsubscribe = subscribe(store, selector, (...args) => {
    if (!done) {
      done = true;
      handler(...args);
    }
    unsubscribe();
  });
}

//
// Handler gets called once when the predicate returns true, the predicate gets passed
// the result of the selector(s), if no predicate is set it defaults to checking if the
// selector(s) return something truthy
//
export function when(store, selector, predicate, handler) {
  if (arguments.length === 3) {
    handler = predicate;

    if (Array.isArray(selector)) {
      predicate = (...args) => {
        for (let i = 0; i < args.length; i++) {
          if (!args[i]) {
            return false;
          }
        }
        return true;
      };
    } else {
      predicate = o => o;
    }
  }

  const state = store.getState();
  if (Array.isArray(selector)) {
    const val = selector.map(s => s(state));
    if (predicate(...val)) {
      return handler(...val);
    }
  } else {
    const val = selector(state);
    if (predicate(val)) {
      return handler(val);
    }
  }

  let done = false;
  const unsubscribe = subscribe(store, selector, (...args) => {
    if (!done && predicate(...args)) {
      done = true;
      handler(...args);
    }
    unsubscribe();
  });
}
