import createHistory from 'history/createBrowserHistory';
import UrlPattern from 'url-pattern';

const history = createHistory();

export const LOCATION_CHANGED = 'ROUTER_LOCATION_CHANGED';
export const PUSH = 'ROUTER_PUSH';
export const REPLACE = 'ROUTER_REPLACE';

export function locationChanged(route, params, location) {
  return {
    type: LOCATION_CHANGED,
    route,
    params,
    location
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

export function routeReducer(state = {}, action) {
  if (action.type === LOCATION_CHANGED) {
    return {
      route: action.route,
      params: action.params,
      location: action.location
    };
  }

  return state;
}

export function routeMiddleware() {
  return next => action => {
    switch (action.type) {
      case PUSH:
        history.push(action.path);
        break;
      case REPLACE:
        history.replace(action.path);
        break;
      default:
        return next(action);
    }
  };
}

function decode(location) {
  location.pathname = decodeURIComponent(location.pathname);
  return location;
}

function match(routes, location) {
  let params;
  for (let i = 0; i < routes.length; i++) {
    params = routes[i].pattern.match(location.pathname);
    if (params !== null) {
      const keys = Object.keys(params);
      for (let j = 0; j < keys.length; j++) {
        params[keys[j]] = decodeURIComponent(params[keys[j]]);
      }
      return locationChanged(routes[i].name, params, decode(location));
    }
  }
  return null;
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
  } else {
    matched = { location: {} };
  }

  history.listen(location => {
    const nextMatch = match(patterns, location);
    if (
      nextMatch &&
      nextMatch.location.pathname !== matched.location.pathname
    ) {
      matched = nextMatch;
      store.dispatch(matched);
    }
  });
}
