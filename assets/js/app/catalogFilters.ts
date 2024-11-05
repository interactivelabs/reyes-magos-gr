const getFilters = () =>
  document.getElementsByName("category[]") as NodeListOf<HTMLInputElement>;

const updateUrlFilters = (params) => {
  params.delete("page");
  const newQueryString = params.toString();
  window.location.replace(`${window.location.pathname}?${newQueryString}`);
};

const getAgeFilters = () => {
  const ageMin = document.getElementsByName("age_min")[0] as HTMLInputElement;
  const ageMax = document.getElementsByName("age_max")[0] as HTMLInputElement;
  return { ageMin: ageMin.value, ageMax: ageMax.value };
};

export const updateFilters = () => {
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

const clearFilters = () => {
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

const removeFilter = (filterValue) => {
  const filters = getFilters();

  if (filterValue === "age_min") {
    const ageMin = document.getElementsByName("age_min")[0] as HTMLInputElement;
    ageMin.value = "0";
  } else if (filterValue === "age_max") {
    const ageMax = document.getElementsByName("age_max")[0] as HTMLInputElement;
    ageMax.value = "0";
  }

  for (const filter of filters) {
    if (filter.value === filterValue) {
      filter.checked = false;
    }
  }

  updateFilters();
};

const setUrlFilters = () => {
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

export default function initCatalogFilters() {
  globalThis.updateFilters = updateFilters;
  globalThis.clearFilters = clearFilters;
  globalThis.removeFilter = removeFilter;
  setUrlFilters();
}
