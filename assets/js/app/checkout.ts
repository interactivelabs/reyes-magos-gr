export {};

declare global {
  var onCodeInputBlur: () => void;
  var beforeCheckoutHandler: (evt: Event) => void;
  var selectToyImage: (thumbnail: HTMLImageElement) => void;
}

function validateCodeInput(): boolean {
  const codeInput = document.getElementById("code") as HTMLInputElement;
  const code = codeInput.value.toUpperCase().trim();
  codeInput.value = code;
  return code.length > 5;
}

function toggleCodeError(valid: boolean) {
  const errorMsg = document.getElementById("code-error") as HTMLElement;
  errorMsg.classList.toggle("hidden", valid);
}

function beforeCheckoutHandler(evt: Event) {
  const valid = validateCodeInput();
  toggleCodeError(valid);
  if (!valid) {
    evt.preventDefault();
  }
}

globalThis.onCodeInputBlur = () => {
  if (validateCodeInput()) {
    toggleCodeError(true);
  }
};

globalThis.beforeCheckoutHandler = beforeCheckoutHandler;

const SELECTED_THUMBNAIL_CLASSES = ["border-4", "border-primary", "brightness-90"];

globalThis.selectToyImage = (thumbnail: HTMLImageElement) => {
  const mainImage = document.getElementById("toy-image-0") as HTMLImageElement;
  mainImage.src = thumbnail.src;

  const thumbnails = thumbnail.parentElement?.children ?? [];
  for (const el of Array.from(thumbnails)) {
    el.classList.remove(...SELECTED_THUMBNAIL_CLASSES);
  }
  thumbnail.classList.add(...SELECTED_THUMBNAIL_CLASSES);
};
