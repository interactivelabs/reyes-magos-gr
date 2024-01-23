// public/js/codes.ts
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

// public/js/admin.ts
initCodes();
//# sourceMappingURL=admin.js.map
