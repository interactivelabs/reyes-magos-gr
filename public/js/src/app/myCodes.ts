const copyCode = async (code: string) => {
  try {
    await navigator.clipboard.writeText(code);
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};

export default function initMycodes() {
  const codesList = document.getElementById("mycodes-code-list");
  if (!codesList) {
    return;
  }

  codesList.addEventListener("click", async (evt) => {
    let target = evt.target as HTMLElement | null | undefined;

    if (!target || target.tagName === "UL") return;

    while (target?.tagName !== "LI") {
      target = target?.parentElement;
    }

    const code = target.id;
    await copyCode(code);
    const toast = document.getElementById(`mycodes-copied-label-${code}`);
    if (toast) {
      toast.classList.remove("hidden");
      setTimeout(() => {
        toast.classList.add("hidden");
      }, 1500);
    }
  });
}
