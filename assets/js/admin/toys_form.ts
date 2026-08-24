declare var addCategory: () => void;
declare var removeCategory: (category: string) => void;

// The category-search <el-autocomplete> opens its options list on any
// pointerdown/focus of the input, showing the full unfiltered list. We only
// want it to open once the user actually types, so we stop the pointerdown
// before it reaches the input (native focus still happens since we only
// stop propagation, not the default action). Delegated on document since the
// form (and its input) is re-created on every modal open.
document.addEventListener(
  "pointerdown",
  (event) => {
    const target = event.target as HTMLElement | null;
    if (target?.id === "category_search") {
      event.stopPropagation();
    }
  },
  true
);

// Selecting an option sets the input's value and dispatches "change" (the
// library's own signal that a selection was committed, as opposed to "input"
// which fires on every keystroke). Treat that the same as clicking "Add" —
// the Add button is reserved for committing a freeform value not in the list.
document.addEventListener("change", (event) => {
  const target = event.target as HTMLElement | null;
  if (target?.id === "category_search") {
    globalThis.addCategory();
  }
});

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
  wrapper.dataset.category = category;

  const pill = document.createElement("span");
  pill.className =
    "inline-flex items-center gap-2 rounded-pill bg-chip px-4 py-2 text-[13px] font-sans text-ink";

  const label = document.createElement("span");
  label.textContent = category;

  const removeButton = document.createElement("button");
  removeButton.type = "button";
  removeButton.className =
    "inline-flex h-4 w-4 shrink-0 items-center justify-center rounded-pill text-muted hover:text-accent";
  removeButton.innerHTML = `
    <span class="sr-only">Quitar categoría</span>
    <svg class="h-2 w-2" stroke="currentColor" fill="none" viewBox="0 0 8 8">
      <path stroke-linecap="round" stroke-width="1.5" d="M1 1l6 6m0-6L1 7"></path>
    </svg>
  `;
  removeButton.addEventListener("click", () =>
    globalThis.removeCategory(category)
  );

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
