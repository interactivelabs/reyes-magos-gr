/// <reference types="alpinejs" />
import type { Alpine as AlpineType } from "alpinejs";

declare global {
  var Alpine: AlpineType;
}

function validateCodeInput(): boolean {
  const codeInput = document.getElementById("code") as HTMLInputElement;
  const code = codeInput.value.toUpperCase().trim();
  return code.length > 5;
}

function beforeCheckoutHandler(evt: Event) {
  if (!validateCodeInput()) {
    evt.preventDefault();
  }
}

document.addEventListener("alpine:init", () => {
  Alpine.data("codeForm", () => ({
    valid: true,

    init() {
      this.valid = true;
    },

    submit(evt: Event) {
      this.valid = validateCodeInput();
      evt.preventDefault();
      evt.stopPropagation();
    },

    update() {
      const codeInput = document.getElementById("code") as HTMLInputElement;
      const code = codeInput.value.toUpperCase().trim();
      codeInput.value = code;
      if (code.length > 5) {
        this.valid = true;
      }
    },
  }));
});

globalThis.beforeCheckoutHandler = beforeCheckoutHandler;
