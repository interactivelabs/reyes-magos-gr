// assets/js/admin/codes.ts
var form = document.getElementById("remove_code_form");
form?.addEventListener("change", (e) => {
  const target = e.target;
  if (target.name === "volunteer_code_ids") {
    const codeCheckbox = target.nextElementSibling;
    codeCheckbox.checked = target.checked;
  }
});
globalThis.selectAllUnsigned = () => {
  const assignCodesForm = document.getElementById("assign-codes-form");
  const checkboxes = assignCodesForm?.getElementsByTagName("input");
  if (!checkboxes || !checkboxes.length) return;
  for (let i = 0; i < checkboxes?.length; i++) {
    if (checkboxes[i].type === "checkbox") {
      checkboxes[i].checked = true;
    }
  }
};

// assets/js/admin/orders.ts
globalThis.orderCompletedToggleClick = (btn) => {
  btn.classList.toggle("bg-indigo-600");
  btn.classList.toggle("bg-gray-200");
  btn.children[0].classList.toggle("translate-x-0");
  btn.children[0].classList.toggle("translate-x-5");
  const completedInput = btn.previousElementSibling;
  completedInput.value = completedInput.value === "1" ? "0" : "1";
};
//# sourceMappingURL=admin.js.map
