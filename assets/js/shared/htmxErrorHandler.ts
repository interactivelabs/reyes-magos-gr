import { HtmxResponseInfo } from "../../@types/htmx.esm";
import { showToast } from "./toasts";

interface HtmxResponseErrorEvent extends Event {
  detail: HtmxResponseInfo;
}

declare global {
  interface WindowEventMap {
    "htmx:responseError": HtmxResponseErrorEvent;
  }
}

window.addEventListener<"htmx:responseError">("htmx:responseError", (e) => {
  console.log(e);
  const code = e.detail.xhr.status;
  if (code === 500) {
    showToast({
      variant: "error",
      title: "Error del servidor",
      message: "Ha ocurrido un error inesperado. Por favor intenta de nuevo.",
    });
  }
});
