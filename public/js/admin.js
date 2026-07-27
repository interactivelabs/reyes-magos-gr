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
document.addEventListener(
  "pointerdown",
  (event) => {
    const target = event.target;
    if (target?.id === "category_search") {
      event.stopPropagation();
    }
  },
  true
);
document.addEventListener("change", (event) => {
  const target = event.target;
  if (target?.id === "category_search") {
    globalThis.addCategory();
  }
});
function getCategoryElements() {
  const input = document.getElementById("category_search");
  const hiddenInput = document.getElementById("category");
  const chipsContainer = document.getElementById(
    "category-chips"
  );
  return { input, hiddenInput, chipsContainer };
}
function createCategoryChip(category) {
  const wrapper = document.createElement("div");
  wrapper.className = "-my-1 flex flex-wrap items-center";
  wrapper.dataset.category = category;
  const pill = document.createElement("span");
  pill.className = "m-1 inline-flex items-center rounded-full border border-gray-200 bg-white py-1.5 pl-3 pr-2 text-sm font-medium text-gray-900";
  const label = document.createElement("span");
  label.textContent = category;
  const removeButton = document.createElement("button");
  removeButton.type = "button";
  removeButton.className = "ml-1 inline-flex h-4 w-4 shrink-0 rounded-full p-1 text-gray-400 hover:bg-gray-200 hover:text-gray-500";
  removeButton.innerHTML = `
    <span class="sr-only">Remove Category</span>
    <svg class="h-2 w-2" stroke="currentColor" fill="none" viewBox="0 0 8 8">
      <path stroke-linecap="round" stroke-width="1.5" d="M1 1l6 6m0-6L1 7"></path>
    </svg>
  `;
  removeButton.addEventListener(
    "click",
    () => globalThis.removeCategory(category)
  );
  pill.append(label, removeButton);
  wrapper.append(pill);
  return wrapper;
}
globalThis.addCategory = () => {
  const { input, hiddenInput, chipsContainer } = getCategoryElements();
  const newCategory = input.value;
  if (!newCategory) return;
  const currentCategories = hiddenInput.value ? hiddenInput.value.split(",") : [];
  if (currentCategories.includes(newCategory)) {
    input.value = "";
    return;
  }
  hiddenInput.value = [...currentCategories, newCategory].join(",");
  input.value = "";
  chipsContainer.classList.remove("hidden");
  chipsContainer.appendChild(createCategoryChip(newCategory));
};
globalThis.removeCategory = (category) => {
  const { hiddenInput, chipsContainer } = getCategoryElements();
  const remaining = hiddenInput.value.split(",").filter((c) => c && c !== category);
  hiddenInput.value = remaining.join(",");
  chipsContainer.querySelector(`[data-category="${CSS.escape(category)}"]`)?.remove();
  if (remaining.length === 0) {
    chipsContainer.classList.add("hidden");
  }
};
//# sourceMappingURL=admin.js.map
