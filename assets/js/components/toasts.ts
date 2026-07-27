export interface Toast {
  variant: "info" | "success" | "warning" | "error";
  title?: string | null;
  message?: string | null;
}

const DISPLAY_DURATION = 5000;
const LEAVE_DURATION = 300;
const MAX_TOASTS = 20;

const ENTER_START_CLASSES = ["translate-y-8"];
const ENTER_END_CLASSES = ["transition", "duration-300", "ease-out", "translate-y-0"];
const LEAVE_CLASSES = [
  "transition",
  "duration-300",
  "ease-in",
  "-translate-x-24",
  "opacity-0",
  "md:translate-x-24",
];

interface ActiveToast {
  el: HTMLElement;
  pause: () => void;
  resume: () => void;
  dismiss: () => void;
}

const activeToasts: ActiveToast[] = [];

function buildToastElement(toast: Toast): HTMLElement | null {
  const template = document.getElementById(
    `toast-template-${toast.variant}`
  ) as HTMLTemplateElement | null;
  if (!template) return null;

  const el = template.content.firstElementChild?.cloneNode(true) as
    | HTMLElement
    | undefined;
  if (!el) return null;

  const titleEl = el.querySelector("[data-toast-title]");
  const messageEl = el.querySelector("[data-toast-message]");

  if (toast.title) {
    titleEl!.textContent = toast.title;
  } else {
    titleEl?.remove();
  }

  if (toast.message) {
    messageEl!.textContent = toast.message;
  } else {
    messageEl?.remove();
  }

  return el;
}

function mountToast(el: HTMLElement, viewport: HTMLElement) {
  let timeout: number | undefined;
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

  const active: ActiveToast = {
    el,
    pause: () => window.clearTimeout(timeout),
    resume: scheduleDismiss,
    dismiss,
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

  window.addEventListener("notify", ((event: CustomEvent<Toast>) => {
    const el = buildToastElement(event.detail);
    if (el) mountToast(el, viewport);
  }) as EventListener);
}

document.addEventListener("DOMContentLoaded", initToastViewport);

export function showToast(toast: Toast) {
  window.dispatchEvent(new CustomEvent("notify", { detail: toast }));
}
