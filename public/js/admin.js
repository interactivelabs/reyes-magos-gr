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

// assets/js/admin/toys_form.ts
window.addEventListener("SearchBox:ItemSelected:category", (event) => {
  const searchBoxEvent = event;
  const item = searchBoxEvent.detail?.item;
  if (!item) return;
  const categoryInput = document.getElementById("category");
  const category = item.Value;
  const currentCategories = categoryInput.value.split(",");
  if (currentCategories.includes(category)) return;
  const newCategories = [...currentCategories, category];
  categoryInput.value = newCategories.join(",");
});
document.addEventListener("alpine:init", () => {
  Alpine.data("CategoryWithSearch", (currentValue) => ({
    category: currentValue,
    setNewCategory(event) {
      const item = event.detail?.item;
      if (!item) return;
      const newCategory = item.Value;
      const currentCategories = !!this.category ? this.category.split(",") : [];
      if (currentCategories.includes(newCategory)) return;
      const newCategories = [...currentCategories, newCategory];
      this.category = newCategories.join(",");
    },
    removeCategory(category) {
      const newCategories = this.category.split(",").filter((c) => c !== category);
      if (newCategories.length > 0) {
        this.category = newCategories.join(",");
        return;
      }
      this.category = "";
    }
  }));
});
//# sourceMappingURL=admin.js.map
