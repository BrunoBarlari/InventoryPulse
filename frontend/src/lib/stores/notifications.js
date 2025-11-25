import { writable } from 'svelte/store';

function createNotificationStore() {
  const { subscribe, update } = writable([]);

  let id = 0;

  function addNotification(message, type = 'info', duration = 4000) {
    const notification = {
      id: ++id,
      message,
      type, // 'success', 'error', 'info', 'warning'
    };

    update(notifications => [...notifications, notification]);

    if (duration > 0) {
      setTimeout(() => {
        removeNotification(notification.id);
      }, duration);
    }

    return notification.id;
  }

  function removeNotification(id) {
    update(notifications => notifications.filter(n => n.id !== id));
  }

  return {
    subscribe,
    success: (msg, duration) => addNotification(msg, 'success', duration),
    error: (msg, duration) => addNotification(msg, 'error', duration),
    info: (msg, duration) => addNotification(msg, 'info', duration),
    warning: (msg, duration) => addNotification(msg, 'warning', duration),
    remove: removeNotification,
  };
}

export const notifications = createNotificationStore();

