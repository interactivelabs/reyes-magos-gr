function getCategoryElements() {
  const input = document.getElementById("category_search") as HTMLInputElement;
  const hiddenInput = document.getElementById("category") as HTMLInputElement;
  const chipsContainer = document.getElementById(
    "category-chips"
  ) as HTMLElement;
  return { input, hiddenInput, chipsContainer };
}

function createCategoryChip(category: string): HTMLElement {
  const wrapper = document.createElement("div");
  wrapper.className = "-my-1 flex flex-wrap items-center";
  wrapper.dataset.category = category;

  const pill = document.createElement("span");
  pill.className =
    "m-1 inline-flex items-center rounded-full border border-gray-200 bg-white py-1.5 pl-3 pr-2 text-sm font-medium text-gray-900";

  const label = document.createElement("span");
  label.textContent = category;

  const removeButton = document.createElement("button");
  removeButton.type = "button";
  removeButton.className =
    "ml-1 inline-flex h-4 w-4 shrink-0 rounded-full p-1 text-gray-400 hover:bg-gray-200 hover:text-gray-500";
  removeButton.innerHTML = `
    <span class="sr-only">Remove Category</span>
    <svg class="h-2 w-2" stroke="currentColor" fill="none" viewBox="0 0 8 8">
      <path stroke-linecap="round" stroke-width="1.5" d="M1 1l6 6m0-6L1 7"></path>
    </svg>
  `;
  removeButton.addEventListener("click", () => removeCategory(category));

  pill.append(label, removeButton);
  wrapper.append(pill);
  return wrapper;
}

globalThis.addCategory = () => {
  const { input, hiddenInput, chipsContainer } = getCategoryElements();
  const newCategory = input.value;
  if (!newCategory) return;

  const currentCategories = hiddenInput.value
    ? hiddenInput.value.split(",")
    : [];
  if (currentCategories.includes(newCategory)) {
    input.value = "";
    return;
  }

  hiddenInput.value = [...currentCategories, newCategory].join(",");
  input.value = "";

  chipsContainer.classList.remove("hidden");
  chipsContainer.appendChild(createCategoryChip(newCategory));
};

globalThis.removeCategory = (category: string) => {
  const { hiddenInput, chipsContainer } = getCategoryElements();
  const remaining = hiddenInput.value
    .split(",")
    .filter((c: string) => c && c !== category);
  hiddenInput.value = remaining.join(",");

  chipsContainer
    .querySelector(`[data-category="${CSS.escape(category)}"]`)
    ?.remove();

  if (remaining.length === 0) {
    chipsContainer.classList.add("hidden");
  }
};
