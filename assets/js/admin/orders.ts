declare var orderCompletedToggleClick: (btn: HTMLButtonElement) => void;

globalThis.orderCompletedToggleClick = (btn: HTMLButtonElement) => {
  btn.classList.toggle("bg-primary");
  btn.classList.toggle("bg-disabled");
  btn.children[0].classList.toggle("translate-x-0");
  btn.children[0].classList.toggle("translate-x-5");
  const completedInput = btn.previousElementSibling as HTMLInputElement;
  const nowCompleted = completedInput.value !== "1";
  completedInput.value = nowCompleted ? "1" : "0";
  btn.setAttribute("aria-checked", nowCompleted ? "true" : "false");
};
