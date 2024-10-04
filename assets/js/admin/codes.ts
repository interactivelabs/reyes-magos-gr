export default function initCodes() {
  const form = document.getElementById("remove_code_form");

  // This is done because of the styling on the checkbox there a styled input and a hidden checkbox input
  form?.addEventListener("change", (e) => {
    const target = e.target as HTMLInputElement;
    if (target.name === "volunteer_code_ids") {
      const codeCheckbox = target.nextElementSibling as HTMLInputElement;
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
}
