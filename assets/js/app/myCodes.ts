const copy = async (code: string) => {
  const flashToast = () => {
    const toast = document.getElementById(`mycodes-copied-label-${code}`);
    if (toast) {
      toast.classList.remove("hidden");
      setTimeout(() => {
        toast.classList.add("hidden");
      }, 1500);
    }
  };

  try {
    await navigator.clipboard.writeText(code);
    flashToast();
    console.log("Copied successfully!");
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};

const share = async (code: string) => {
  const data = {
    title: "Comparte la alegria!",
    text: `Utiliza este codigo para obtener un juguete: ${code}`,
    url: window.location.href,
  };

  try {
    await navigator.share(data);
    console.log("Shared successfully!");
  } catch (err) {
    console.error("Share failed: ", err);
  }
};

export default function initMycodes() {
  globalThis.shareCode = async (code: string) => {
    if (typeof navigator.share === "undefined") {
      copy(code);
    } else {
      share(code);
    }
  };
}
