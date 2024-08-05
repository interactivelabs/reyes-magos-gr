import { showToast } from "../shared/toast";

const copy = async (code: string) => {
  try {
    await navigator.clipboard.writeText(code);
    showToast({
      title: "Copiado!",
      subTitle: "El codigo ha sido copiado al portapapeles.",
    });
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};

const share = async (code: string) => {
  const data = {
    title: "Comparte la alegria!",
    text: `Utiliza este codigo para obtener un juguete: ${code}`,
    url: `${window.location.origin}/catalog`,
  };

  try {
    await navigator.share(data);
  } catch (err) {
    console.error("Share failed: ", err);
  }
};

export default function initMycodes() {
  globalThis.copyCode = copy;
  globalThis.shareCode = async (code: string) => {
    if (typeof navigator.share === "undefined") {
      copy(code);
    } else {
      share(code);
    }
  };
}
