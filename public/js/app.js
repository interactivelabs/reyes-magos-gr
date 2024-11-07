// assets/js/shared/toast.ts
var hideToast = (variant = "" /* SUCCESS */) => {
  const variantSelector = variant ? `${variant}-` : "";
  const toastContainer = document.getElementById(
    `toast-${variantSelector}container`
  );
  const toastPanel = document.getElementById(`toast-${variantSelector}panel`);
  toastPanel?.classList.remove("show-toast");
  toastPanel?.classList.add("hide-toast");
  setTimeout(() => {
    toastContainer?.classList.remove("flex");
    toastContainer?.classList.add("hidden");
  }, 300);
};
var showToast = ({
  title,
  subTitle,
  duration = 5e3,
  variant = "" /* SUCCESS */
}) => {
  const variantSelector = variant ? `${variant}-` : "";
  const toastContainer = document.getElementById(
    `toast-${variantSelector}container`
  );
  const toastPanel = document.getElementById(`toast-${variantSelector}panel`);
  toastContainer?.classList.remove("hidden");
  toastContainer?.classList.add("flex");
  toastPanel?.classList.remove("hide-toast");
  toastPanel?.classList.add("show-toast");
  const toastTitle = document.getElementById(`toast-${variantSelector}title`);
  const toastSubTitle = document.getElementById(
    `toast-${variantSelector}subtitle`
  );
  toastTitle.innerText = title;
  toastSubTitle.innerText = subTitle;
  setTimeout(() => hideToast(variant), duration);
};
var showErrorToast = (props) => showToast({ ...props, variant: "error" /* ERROR */ });
var hideErrorToast = () => hideToast("error" /* ERROR */);
globalThis.showToast = showToast;
globalThis.hideToast = hideToast;
globalThis.showErrorToast = showErrorToast;
globalThis.hideErrorToast = hideErrorToast;

// assets/js/app/myCodes.ts
var copy = async (code) => {
  try {
    await navigator.clipboard.writeText(code);
    showToast({
      title: "Copiado!",
      subTitle: "El codigo ha sido copiado al portapapeles."
    });
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};
var share = async (code) => {
  const data = {
    title: "Comparte la alegria!",
    text: `Utiliza este codigo para obtener un juguete: ${code}`,
    url: `${window.location.origin}/catalog`
  };
  try {
    await navigator.share(data);
  } catch (err) {
    console.error("Share failed: ", err);
  }
};
globalThis.copyCode = copy;
globalThis.shareCode = async (code) => {
  if (typeof navigator.share === "undefined") {
    copy(code);
  } else {
    share(code);
  }
};

// assets/js/app/catalogFilters.ts
var getFilters = () => document.getElementsByName("category[]");
var updateUrlFilters = (params) => {
  params.delete("page");
  const newQueryString = params.toString();
  window.location.replace(`${window.location.pathname}?${newQueryString}`);
};
var getAgeFilters = () => {
  const ageMin = document.getElementsByName("age_min")[0];
  const ageMax = document.getElementsByName("age_max")[0];
  return { ageMin: ageMin.value, ageMax: ageMax.value };
};
var updateFilters = () => {
  const filters = getFilters();
  const selectedFilters = Array.from(filters).filter(
    (filter) => filter.checked
  );
  const selectedFiltersValues = selectedFilters.map((filter) => filter.value);
  const params = new URLSearchParams(window.location.search);
  params.delete("category");
  params.delete("age_min");
  params.delete("age_max");
  if (selectedFiltersValues.length > 0) {
    for (const filter of selectedFiltersValues) {
      params.append("category", filter);
    }
  }
  const { ageMin, ageMax } = getAgeFilters();
  if (parseInt(ageMin, 10) > 1) {
    params.set("age_min", ageMin);
  }
  if (parseInt(ageMax, 10) > 1) {
    params.set("age_max", ageMax);
  }
  updateUrlFilters(params);
};
var clearFilters = () => {
  const filters = getFilters();
  for (const filter of filters) {
    filter.checked = false;
  }
  const params = new URLSearchParams(window.location.search);
  params.delete("category");
  params.delete("age_min");
  params.delete("age_max");
  updateUrlFilters(params);
};
var removeFilter = (filterValue) => {
  const filters = getFilters();
  if (filterValue === "age_min") {
    const ageMin = document.getElementsByName("age_min")[0];
    ageMin.value = "0";
  } else if (filterValue === "age_max") {
    const ageMax = document.getElementsByName("age_max")[0];
    ageMax.value = "0";
  }
  for (const filter of filters) {
    if (filter.value === filterValue) {
      filter.checked = false;
    }
  }
  updateFilters();
};
var setUrlFilters = () => {
  const params = new URLSearchParams(window.location.search);
  const filters = getFilters();
  if (params.has("category")) {
    const selectedFilters = params.getAll("category");
    for (const filter of filters) {
      if (selectedFilters.includes(filter.value)) {
        filter.checked = true;
      }
    }
  }
};
globalThis.updateFilters = updateFilters;
globalThis.clearFilters = clearFilters;
globalThis.removeFilter = removeFilter;
setUrlFilters();

// assets/js/shared/htmxErrorHandler.ts
window.addEventListener("htmx:responseError", (e) => {
  console.log(e);
  const code = e.detail.xhr.status;
  if (code === 500) {
    showErrorToast({
      title: "Error del servidor",
      subTitle: "Ha ocurrido un error inesperado. Por favor intenta de nuevo."
    });
  }
});
//# sourceMappingURL=app.js.map
