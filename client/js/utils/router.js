import history from 'history/browser';
import UrlPattern from 'url-pattern';

export const LOCATION_CHANGED = 'ROUTER_LOCATION_CHANGED';
export const PUSH = 'ROUTER_PUSH';
export const REPLACE = 'ROUTER_REPLACE';

export function locationChanged(route, params, location) {
  Object.keys(params).forEach(key => {
    params[key] = decodeURIComponent(params[key]);
  });

  const query = {};
  new URLSearchParams(location.search).forEach((value, key) => {
    query[key] = value;
  });

  return {
    type: LOCATION_CHANGED,
    route,
    params,
    query,
    path: decodeURIComponent(location.pathname)
  };
}

export function push(path) {
  return {
    type: PUSH,
    path
  };
}

export function replace(path) {
  return {
    type: REPLACE,
    path
  };
}

export function routeReducer(state = {}, { type, ...action }) {
  if (type === LOCATION_CHANGED) {
    return action;
  }

  return state;
}

export function routeMiddleware() {
  return next => action => {
    switch (action.type) {
      case PUSH:
        history.push(`${action.path}`);
        break;
      case REPLACE:
        history.replace(action.path);
        break;
      default:
        return next(action);
    }
  };
}

function match(routes, location) {
  for (let i = 0; i < routes.length; i++) {
    const params = routes[i].pattern.match(location.pathname);
    if (params !== null) {
      return locationChanged(routes[i].name, params, location);
    }
  }
}

export default function initRouter(routes, store) {
  const patterns = [];
  const opts = {
    segmentValueCharset: 'a-zA-Z0-9-_.%'
  };

  Object.keys(routes).forEach(name =>
    patterns.push({
      name,
      pattern: new UrlPattern(routes[name], opts)
    })
  );

  let matched = match(patterns, history.location);
  if (matched) {
    store.dispatch(matched);
  }

  history.listen(({ location }) => {
    const nextMatch = match(patterns, location);
    if (nextMatch && nextMatch.path !== matched?.path) {
      matched = nextMatch;
      store.dispatch(matched);
    }
  });
}
