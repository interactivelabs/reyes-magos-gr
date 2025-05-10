import { HtmxResponseInfo } from "../../@types/htmx.esm";
import { showToast } from "../components/toasts";

interface HtmxResponseErrorEvent extends Event {
  detail: HtmxResponseInfo;
}

declare global {
  interface WindowEventMap {
    "htmx:responseError": HtmxResponseErrorEvent;
  }
}

const isGoValidatorError = (responseText: string) =>
  responseText &&
  responseText.includes("message") &&
  responseText.includes("Key");

const parseGoValidatorError = (errorText: string) => {
  const errors: { key: string; message: string }[] = [];
  const errorLines = errorText
    .split(",")
    .filter((line) => line.includes("Key"))
    .map((line) => line.trim().replace("message=", ""));
  for (const line of errorLines) {
    const lineErrors = line.split("\\n");
    for (let i = 0; i < lineErrors.length; i++) {
      const error = lineErrors[i];
      const keyMatch = error.match(/Key: '([^']+)'/);
      const messageMatch = error.match(/Error:(.+)/);
      if (keyMatch && messageMatch) {
        errors.push({
          key: keyMatch[1],
          message: messageMatch[1].trim(),
        });
      }
    }
  }
  return errors;
};

window.addEventListener<"htmx:responseError">(
  "htmx:responseError",
  (e: HtmxResponseErrorEvent) => {
    const code = e.detail.xhr.status;
    if (code === 500) {
      showToast({
        variant: "error",
        title: "Error del servidor",
        message: "Ha ocurrido un error inesperado. Por favor intenta de nuevo.",
      });
    }
    if (code === 400) {
      const { responseText } = e.detail.xhr;
      if (isGoValidatorError(responseText)) {
        const errors = parseGoValidatorError(responseText);
        const { target } = e;
        if (target instanceof HTMLFormElement) {
          const form = target;
          for (const error of errors) {
            const fieldName = error.key.split(".").pop();
            const input = form.querySelector(
              `[name="${fieldName}"]`,
            ) as HTMLInputElement;
            if (input) {
              input.setCustomValidity(error.message);
              input.reportValidity();
            }
          }
        }
      }
    }
  },
);
