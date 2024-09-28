export enum ToastVariants {
  SUCCESS = "",
  ERROR = "error",
}

export interface IShowToast {
  title: string;
  subTitle: string;
  duration?: number;
  variant?: ToastVariants;
}

export const hideToast = (variant: ToastVariants = ToastVariants.SUCCESS) => {
  const variantSelector = variant ? `${variant}-` : "";
  const toastContainer = document.getElementById(
    `toast-${variantSelector}container`
  );
  const toastPanel = document.getElementById(`toast-${variantSelector}panel`);

  toastPanel?.classList.remove("show-toast");
  toastPanel?.classList.add("hide-toast");

  setTimeout(() => {
    toastContainer?.classList.remove("flex");
    toastContainer?.classList.add("hidden");
  }, 300);
};

export const showToast = ({
  title,
  subTitle,
  duration = 5000,
  variant = ToastVariants.SUCCESS,
}: IShowToast) => {
  const variantSelector = variant ? `${variant}-` : "";
  const toastContainer = document.getElementById(
    `toast-${variantSelector}container`
  );
  const toastPanel = document.getElementById(`toast-${variantSelector}panel`);

  toastContainer?.classList.remove("hidden");
  toastContainer?.classList.add("flex");

  toastPanel?.classList.remove("hide-toast");
  toastPanel?.classList.add("show-toast");

  const toastTitle = document.getElementById(`toast-${variantSelector}title`);
  const toastSubTitle = document.getElementById(
    `toast-${variantSelector}subtitle`
  );

  toastTitle!.innerText = title;
  toastSubTitle!.innerText = subTitle;

  setTimeout(() => hideToast(variant), duration);
};

export const showErrorToast = (props: IShowToast) =>
  showToast({ ...props, variant: ToastVariants.ERROR });

export const hideErrorToast = () => hideToast(ToastVariants.ERROR);

export default function initToast() {
  globalThis.showToast = showToast;
  globalThis.hideToast = hideToast;
  globalThis.showErrorToast = showErrorToast;
  globalThis.hideErrorToast = hideErrorToast;
}
