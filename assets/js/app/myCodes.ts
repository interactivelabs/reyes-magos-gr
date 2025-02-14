import { showToast } from "../components/toasts";

const copy = async (code: string) => {
  try {
    await navigator.clipboard.writeText(code);
    showToast({
      variant: "error",
      title: "Copiado!",
      message: "El codigo ha sido copiado al portapapeles.",
    });
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};

const share = async (code: string) => {
  const data = {
    title: "Comparte la alegria!",
    text: `Utiliza este codigo para obtener un regalo: ${code}`,
    url: `${window.location.origin}/catalog?code=${code}`,
  };

  try {
    await navigator.share(data);
  } catch (err) {
    console.error("Share failed: ", err);
  }
};

globalThis.copyCode = copy;
globalThis.shareCode = async (code: string) => {
  if (typeof navigator.share === "undefined") {
    copy(code);
  } else {
    share(code);
  }
};
