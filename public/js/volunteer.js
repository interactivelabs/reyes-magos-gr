// assets/js/components/toasts.ts
var DISPLAY_DURATION = 5e3;
var LEAVE_DURATION = 300;
var MAX_TOASTS = 20;
var ENTER_START_CLASSES = ["translate-y-8"];
var ENTER_END_CLASSES = ["transition", "duration-300", "ease-out", "translate-y-0"];
var LEAVE_CLASSES = [
  "transition",
  "duration-300",
  "ease-in",
  "-translate-x-24",
  "opacity-0",
  "md:translate-x-24"
];
var activeToasts = [];
function buildToastElement(toast) {
  const template = document.getElementById(
    `toast-template-${toast.variant}`
  );
  if (!template) return null;
  const el = template.content.firstElementChild?.cloneNode(true);
  if (!el) return null;
  const titleEl = el.querySelector("[data-toast-title]");
  const messageEl = el.querySelector("[data-toast-message]");
  if (toast.title) {
    titleEl.textContent = toast.title;
  } else {
    titleEl?.remove();
  }
  if (toast.message) {
    messageEl.textContent = toast.message;
  } else {
    messageEl?.remove();
  }
  return el;
}
function mountToast(el, viewport) {
  let timeout;
  let dismissed = false;
  const dismiss = () => {
    if (dismissed) return;
    dismissed = true;
    window.clearTimeout(timeout);
    const index = activeToasts.findIndex((toast) => toast.el === el);
    if (index !== -1) activeToasts.splice(index, 1);
    el.classList.remove("translate-y-0", "translate-x-0", "opacity-100");
    el.classList.add(...LEAVE_CLASSES);
    setTimeout(() => el.remove(), LEAVE_DURATION);
  };
  const scheduleDismiss = () => {
    timeout = window.setTimeout(dismiss, DISPLAY_DURATION);
  };
  const active = {
    el,
    pause: () => window.clearTimeout(timeout),
    resume: scheduleDismiss,
    dismiss
  };
  activeToasts.push(active);
  el.querySelector("[data-toast-dismiss]")?.addEventListener("click", dismiss);
  viewport.appendChild(el);
  el.classList.add(...ENTER_START_CLASSES);
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      el.classList.remove(...ENTER_START_CLASSES);
      el.classList.add(...ENTER_END_CLASSES);
    });
  });
  scheduleDismiss();
  while (activeToasts.length > MAX_TOASTS) {
    activeToasts[0].dismiss();
  }
}
function initToastViewport() {
  const viewport = document.getElementById("toast-viewport");
  if (!viewport) return;
  viewport.addEventListener("mouseenter", () => {
    for (const toast of activeToasts) toast.pause();
  });
  viewport.addEventListener("mouseleave", () => {
    for (const toast of activeToasts) toast.resume();
  });
  window.addEventListener("notify", ((event) => {
    const el = buildToastElement(event.detail);
    if (el) mountToast(el, viewport);
  }));
}
document.addEventListener("DOMContentLoaded", initToastViewport);
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
