// assets/js/components/toasts.ts
document.addEventListener("alpine:init", () => {
  Alpine.data("toasts", () => ({
    notifications: [],
    displayDuration: 5e3,
    addNotification(notification) {
      const id = Date.now();
      if (this.notifications.length >= 20) {
        this.notifications.splice(0, this.notifications.length - 19);
      }
      this.notifications.push({ ...notification, id });
    },
    removeNotification(id) {
      setTimeout(() => {
        this.notifications = this.notifications.filter(
          (notification) => notification.id !== id
        );
      }, 400);
    }
  }));
});
function showToast(toast) {
  window.dispatchEvent(new CustomEvent("notify", { detail: toast }));
}

// assets/js/app/myCodes.ts
var copy = async (code) => {
  try {
    await navigator.clipboard.writeText(code);
    showToast({
      variant: "error",
      title: "Copiado!",
      message: "El codigo ha sido copiado al portapapeles."
    });
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};
var share = async (code) => {
  const data = {
    title: "Comparte la alegria!",
    text: `Utiliza este codigo para obtener un regalo: ${code}`,
    url: `${window.location.origin}/catalog?code=${code}`
  };
  try {
    await navigator.share(data);
  } catch (err) {
    console.error("Share failed: ", err);
  }
};
globalThis.copyCode = copy;
globalThis.shareCode = async (code) => {
  if (typeof navigator.share === "undefined") {
    copy(code);
  } else {
    share(code);
  }
};
//# sourceMappingURL=volunteer.js.map
