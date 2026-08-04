// assets/js/components/toasts.ts
var DISPLAY_DURATION = 5e3;
var LEAVE_DURATION = 300;
var MAX_TOASTS = 20;
var ENTER_START_CLASSES = ["translate-y-8"];
var ENTER_END_CLASSES = ["transition", "duration-300", "ease-out", "translate-y-0"];
var LEAVE_CLASSES = [
  "transition",
  "duration-300",
  "ease-in",
  "-translate-x-24",
  "opacity-0",
  "md:translate-x-24"
];
var activeToasts = [];
function buildToastElement(toast) {
  const template = document.getElementById(
    `toast-template-${toast.variant}`
  );
  if (!template) return null;
  const el = template.content.firstElementChild?.cloneNode(true);
  if (!el) return null;
  const titleEl = el.querySelector("[data-toast-title]");
  const messageEl = el.querySelector("[data-toast-message]");
  if (toast.title) {
    titleEl.textContent = toast.title;
  } else {
    titleEl?.remove();
  }
  if (toast.message) {
    messageEl.textContent = toast.message;
  } else {
    messageEl?.remove();
  }
  return el;
}
function mountToast(el, viewport) {
  let timeout;
  let dismissed = false;
  const dismiss = () => {
    if (dismissed) return;
    dismissed = true;
    window.clearTimeout(timeout);
    const index = activeToasts.findIndex((toast) => toast.el === el);
    if (index !== -1) activeToasts.splice(index, 1);
    el.classList.remove("translate-y-0", "translate-x-0", "opacity-100");
    el.classList.add(...LEAVE_CLASSES);
    setTimeout(() => el.remove(), LEAVE_DURATION);
  };
  const scheduleDismiss = () => {
    timeout = window.setTimeout(dismiss, DISPLAY_DURATION);
  };
  const active = {
    el,
    pause: () => window.clearTimeout(timeout),
    resume: scheduleDismiss,
    dismiss
  };
  activeToasts.push(active);
  el.querySelector("[data-toast-dismiss]")?.addEventListener("click", dismiss);
  viewport.appendChild(el);
  el.classList.add(...ENTER_START_CLASSES);
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      el.classList.remove(...ENTER_START_CLASSES);
      el.classList.add(...ENTER_END_CLASSES);
    });
  });
  scheduleDismiss();
  while (activeToasts.length > MAX_TOASTS) {
    activeToasts[0].dismiss();
  }
}
function initToastViewport() {
  const viewport = document.getElementById("toast-viewport");
  if (!viewport) return;
  viewport.addEventListener("mouseenter", () => {
    for (const toast of activeToasts) toast.pause();
  });
  viewport.addEventListener("mouseleave", () => {
    for (const toast of activeToasts) toast.resume();
  });
  window.addEventListener("notify", ((event) => {
    const el = buildToastElement(event.detail);
    if (el) mountToast(el, viewport);
  }));
}
document.addEventListener("DOMContentLoaded", initToastViewport);
function showToast(toast) {
  window.dispatchEvent(new CustomEvent("notify", { detail: toast }));
}

// assets/js/admin/codes.ts
var assignForm = document.getElementById("assign-codes-form");
var removeForm = document.getElementById("remove-codes-form");
function selectedCountLabel(count) {
  return count === 1 ? "1 seleccionado" : `${count} seleccionados`;
}
function updateAssignState() {
  if (!assignForm) return;
  const checked = assignForm.querySelectorAll(".unassigned-checkbox:checked");
  const all = assignForm.querySelectorAll(".unassigned-checkbox");
  const volunteerSelect = assignForm.querySelector("#volunteer_id");
  const submit = assignForm.querySelector("#assign-submit");
  const countLabel = assignForm.querySelector("#assign-selected-count");
  const selectAll = assignForm.querySelector("#select-all-unassigned");
  const count = checked.length;
  if (countLabel) countLabel.textContent = selectedCountLabel(count);
  if (submit) {
    submit.textContent = `Asignar (${count})`;
    submit.disabled = count === 0 || !volunteerSelect?.value;
  }
  if (selectAll) selectAll.checked = all.length > 0 && count === all.length;
}
function updateRemoveState() {
  if (!removeForm) return;
  const checked = removeForm.querySelectorAll(".assigned-checkbox:checked");
  const all = removeForm.querySelectorAll(".assigned-checkbox");
  const submit = removeForm.querySelector("#remove-submit");
  const countLabel = removeForm.querySelector("#remove-selected-count");
  const selectAll = removeForm.querySelector("#select-all-assigned");
  const count = checked.length;
  if (countLabel) countLabel.textContent = selectedCountLabel(count);
  if (submit) {
    submit.textContent = `Quitar (${count})`;
    submit.disabled = count === 0;
  }
  if (selectAll) selectAll.checked = all.length > 0 && count === all.length;
}
assignForm?.addEventListener("change", (e) => {
  const target = e.target;
  if (target.id === "select-all-unassigned") {
    assignForm.querySelectorAll(".unassigned-checkbox").forEach((cb) => cb.checked = target.checked);
  }
  updateAssignState();
});
removeForm?.addEventListener("change", (e) => {
  const target = e.target;
  if (target.id === "select-all-assigned") {
    removeForm.querySelectorAll(".assigned-checkbox").forEach((cb) => cb.checked = target.checked);
  } else if (target.classList.contains("assigned-checkbox")) {
    const codeCheckbox = target.nextElementSibling;
    codeCheckbox.checked = target.checked;
  }
  updateRemoveState();
});
updateAssignState();
updateRemoveState();
function flashSuccessFromQuery() {
  const params = new URLSearchParams(window.location.search);
  const messages = {
    created: (c) => c === "1" ? "1 c\xF3digo generado" : `${c} c\xF3digos generados`,
    assigned: (c) => c === "1" ? "1 c\xF3digo asignado" : `${c} c\xF3digos asignados`,
    removed: (c) => c === "1" ? "1 c\xF3digo eliminado" : `${c} c\xF3digos eliminados`
  };
  for (const key of Object.keys(messages)) {
    const value = params.get(key);
    if (value === null) continue;
    showToast({ variant: "success", title: "Listo", message: messages[key](value) });
    params.delete(key);
    const query = params.toString();
    window.history.replaceState({}, "", window.location.pathname + (query ? `?${query}` : ""));
    break;
  }
}
flashSuccessFromQuery();

// assets/js/admin/orders.ts
globalThis.orderCompletedToggleClick = (btn) => {
  btn.classList.toggle("bg-primary");
  btn.classList.toggle("bg-disabled");
  btn.children[0].classList.toggle("translate-x-0");
  btn.children[0].classList.toggle("translate-x-5");
  const completedInput = btn.previousElementSibling;
  const nowCompleted = completedInput.value !== "1";
  completedInput.value = nowCompleted ? "1" : "0";
  btn.setAttribute("aria-checked", nowCompleted ? "true" : "false");
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
