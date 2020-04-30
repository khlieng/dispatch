import { skipWaiting, clientsClaim } from 'workbox-core';
import { precacheAndRoute, createHandlerBoundToURL } from 'workbox-precaching';
import { NavigationRoute, registerRoute } from 'workbox-routing';

skipWaiting();
clientsClaim();

precacheAndRoute(self.__WB_MANIFEST, {
  ignoreUrlParametersMatching: [/.*/]
});

const handler = createHandlerBoundToURL('/');
registerRoute(new NavigationRoute(handler));
