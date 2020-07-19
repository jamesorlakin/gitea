import {registerRoute} from 'workbox-routing';
import {StaleWhileRevalidate} from 'workbox-strategies';

const cacheName = 'static-cache-v2';

// disable workbox debug logging in development, remove when debugging the service worker
self.__WB_DISABLE_DEV_LOGS = true;

// see https://developer.mozilla.org/en-US/docs/Web/API/RequestDestination for possible values
const cachedDestinations = new Set([
  'font',
  'manifest',
  'paintworklet',
  'script',
  'sharedworker',
  'style',
  'worker',
]);

registerRoute(
  ({request}) => cachedDestinations.has(request.destination),
  new StaleWhileRevalidate({cacheName}),
);

// Push notification event handling
self.addEventListener('push', (event) => {
  const eventPayload = event.data.json();
  const options = {
    body: eventPayload.text,
    vibrate: [100, 50, 100],
    data: {
      url: eventPayload.url
    }
  };

  event.waitUntil(self.registration.showNotification(eventPayload.title, options));
});

self.addEventListener('notificationclick', (event) => {
  const notification = event.notification;
  const url = notification.data.url;

  clients.openWindow(url);
  notification.close();
});