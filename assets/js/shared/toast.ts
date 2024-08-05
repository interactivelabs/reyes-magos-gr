export const hideToast = () => {
  const toastContainer = document.getElementById("toast-container");
  const toastPanel = document.getElementById("toast-panel");

  toastPanel?.classList.remove("show-toast");
  toastPanel?.classList.add("hide-toast");

  setTimeout(() => {
    toastContainer?.classList.remove("flex");
    toastContainer?.classList.add("hidden");
  }, 300);
};

export interface IShowToast {
  title: string;
  subTitle: string;
  duration?: number;
}

export const showToast = ({ title, subTitle, duration = 5000 }: IShowToast) => {
  const toastContainer = document.getElementById("toast-container");
  const toastPanel = document.getElementById("toast-panel");

  toastContainer?.classList.remove("hidden");
  toastContainer?.classList.add("flex");

  toastPanel?.classList.remove("hide-toast");
  toastPanel?.classList.add("show-toast");

  const toastTitle = document.getElementById("toast-title");
  const toastSubTitle = document.getElementById("toast-subtitle");

  toastTitle!.innerText = title;
  toastSubTitle!.innerText = subTitle;

  setTimeout(() => hideToast(), duration);
};

export default function initToast() {
  globalThis.hideToast = hideToast;
}
