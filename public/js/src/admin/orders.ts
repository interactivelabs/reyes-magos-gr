export default function initOrders() {
  globalThis.orderCompletedToggleClick = (btn: HTMLButtonElement) => {
    btn.classList.toggle("bg-indigo-600");
    btn.classList.toggle("bg-gray-200");
    btn.children[0].classList.toggle("translate-x-0");
    btn.children[0].classList.toggle("translate-x-5");
    const completedInput = btn.previousElementSibling as HTMLInputElement;
    completedInput.value = completedInput.value === "1" ? "0" : "1";
  };
}
