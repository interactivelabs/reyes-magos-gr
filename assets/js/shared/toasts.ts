export interface Toast {
  variant: "info" | "success" | "warning" | "error";
  title: string | null;
  message: string | null;
}

document.addEventListener("alpine:init", () => {
  Alpine.data("toasts", () => ({
    notifications: [],
    displayDuration: 5000,

    addNotification(notification: Toast) {
      const id = Date.now();

      // Keep only the most recent 20 notifications
      if (this.notifications.length >= 20) {
        this.notifications.splice(0, this.notifications.length - 19);
      }

      // Add the new notification to the notifications stack
      this.notifications.push({ ...notification, id });
    },
    removeNotification(id) {
      setTimeout(() => {
        this.notifications = this.notifications.filter(
          (notification) => notification.id !== id
        );
      }, 400);
    },
  }));
});

export function showToast(toast: Toast) {
  window.dispatchEvent(new CustomEvent("notify", { detail: toast }));
}
