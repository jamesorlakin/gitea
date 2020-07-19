// Handles the opt in and out options on the notifications page
// Notifications themselves are produced in serviceworker.js

/**
 * @param {ServiceWorkerRegistration} registration
 */
export default async function initPushNotificationsOptIn(registration) {
  if (!window.location.pathname.startsWith('/notifications')) return;
  if (!('PushManager' in window && window.config.PushNotificationsEnabled)) {
    return;
  }

  const subscription = await registration.pushManager.getSubscription();
  const button = $('#pushnotifications_opt_in');
  if (subscription) {
    button.text(button.attr('data-opt-out'));
    button.on('click', () => unsubscribe(registration));
  } else {
    button.text(button.attr('data-opt-in'));
    button.on('click', () => subscribe(registration));
  }
  button.css('display', '');
}

/**
 * @param {ServiceWorkerRegistration} registration
 */
async function subscribe(registration) {
  const canNotify = await hasNotificationPermission();
  if (!canNotify) return false;

  const subscriptionResults = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: window.config.PushNotificationsPublicKey
  });
  await createGiteaServerSubscription(subscriptionResults.toJSON());
  initPushNotificationsOptIn(registration); // Update the button
}

/**
 * @param {ServiceWorkerRegistration} registration
 */
async function unsubscribe(registration) {
  // TODO: A delete API...
  const subscription = await registration.pushManager.getSubscription();
  await subscription.unsubscribe();
  initPushNotificationsOptIn(registration); // Update the button
}

async function createGiteaServerSubscription(subscriptionJSON) {
  try {
    const request = await fetch(`${window.config.AppSubUrl}/api/v1/notifications/subscription`, {
      credentials: 'include',
      method: 'POST',
      headers: {
        'Content-Type': 'application/json; charset=utf-8',
      },
      body: JSON.stringify({
        endpoint: subscriptionJSON.endpoint,
        auth: subscriptionJSON.keys.auth,
        p256dh: subscriptionJSON.keys.p256dh
      })
    });
    if (request.status === 201) return true;
  } catch (error) {
    console.error(error);
  }
  return false;
}

async function hasNotificationPermission() {
  const requestResult = await requestNotificationPermission();
  if (requestResult === 'granted') {
    return true;
  }
  return false;
}

function requestNotificationPermission() {
  return new Promise((resolve, reject) => {
    // This used to be callback-based instead of a Promise. We account for that here:
    const permissionResult = Notification.requestPermission((result) => {
      return resolve(result);
    });

    if (permissionResult) {
      permissionResult.then(resolve, reject);
    }
  });
}