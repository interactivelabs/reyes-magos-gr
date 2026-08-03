import { showToast } from "../components/toasts";

const assignForm = document.getElementById("assign-codes-form");
const removeForm = document.getElementById("remove-codes-form");

function selectedCountLabel(count: number): string {
  return count === 1 ? "1 seleccionado" : `${count} seleccionados`;
}

function updateAssignState() {
  if (!assignForm) return;
  const checked = assignForm.querySelectorAll<HTMLInputElement>(".unassigned-checkbox:checked");
  const all = assignForm.querySelectorAll<HTMLInputElement>(".unassigned-checkbox");
  const volunteerSelect = assignForm.querySelector<HTMLSelectElement>("#volunteer_id");
  const submit = assignForm.querySelector<HTMLButtonElement>("#assign-submit");
  const countLabel = assignForm.querySelector<HTMLElement>("#assign-selected-count");
  const selectAll = assignForm.querySelector<HTMLInputElement>("#select-all-unassigned");
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
  const checked = removeForm.querySelectorAll<HTMLInputElement>(".assigned-checkbox:checked");
  const all = removeForm.querySelectorAll<HTMLInputElement>(".assigned-checkbox");
  const submit = removeForm.querySelector<HTMLButtonElement>("#remove-submit");
  const countLabel = removeForm.querySelector<HTMLElement>("#remove-selected-count");
  const selectAll = removeForm.querySelector<HTMLInputElement>("#select-all-assigned");
  const count = checked.length;

  if (countLabel) countLabel.textContent = selectedCountLabel(count);
  if (submit) {
    submit.textContent = `Quitar (${count})`;
    submit.disabled = count === 0;
  }
  if (selectAll) selectAll.checked = all.length > 0 && count === all.length;
}

assignForm?.addEventListener("change", (e) => {
  const target = e.target as HTMLInputElement;
  if (target.id === "select-all-unassigned") {
    assignForm
      .querySelectorAll<HTMLInputElement>(".unassigned-checkbox")
      .forEach((cb) => (cb.checked = target.checked));
  }
  updateAssignState();
});

removeForm?.addEventListener("change", (e) => {
  const target = e.target as HTMLInputElement;
  if (target.id === "select-all-assigned") {
    removeForm
      .querySelectorAll<HTMLInputElement>(".assigned-checkbox")
      .forEach((cb) => (cb.checked = target.checked));
  } else if (target.classList.contains("assigned-checkbox")) {
    // Styled checkbox pairs with a hidden `code_ids` checkbox in the next cell.
    const codeCheckbox = target.nextElementSibling as HTMLInputElement;
    codeCheckbox.checked = target.checked;
  }
  updateRemoveState();
});

updateAssignState();
updateRemoveState();

function flashSuccessFromQuery() {
  const params = new URLSearchParams(window.location.search);
  const messages: Record<string, (count: string) => string> = {
    created: (c) => (c === "1" ? "1 código generado" : `${c} códigos generados`),
    assigned: (c) => (c === "1" ? "1 código asignado" : `${c} códigos asignados`),
    removed: (c) => (c === "1" ? "1 código eliminado" : `${c} códigos eliminados`),
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
