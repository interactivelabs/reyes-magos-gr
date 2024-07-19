// assets/js/admin/codes.ts
function initCodes() {
  const form = document.getElementById("remove_code_form");
  form?.addEventListener("change", (e) => {
    const target = e.target;
    if (target.name === "volunteer_code_ids") {
      const codeCheckbox = target.nextElementSibling;
      codeCheckbox.checked = target.checked;
    }
  });
}

// assets/js/admin/orders.ts
function initOrders() {
  globalThis.orderCompletedToggleClick = (btn) => {
    btn.classList.toggle("bg-indigo-600");
    btn.classList.toggle("bg-gray-200");
    btn.children[0].classList.toggle("translate-x-0");
    btn.children[0].classList.toggle("translate-x-5");
    const completedInput = btn.previousElementSibling;
    completedInput.value = completedInput.value === "1" ? "0" : "1";
  };
}

// assets/js/admin.ts
initCodes();
initOrders();
//# sourceMappingURL=admin.js.map
