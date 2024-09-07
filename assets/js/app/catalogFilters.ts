const getFilters = () =>
  document.getElementsByName("category[]") as NodeListOf<HTMLInputElement>;

const updateUrlFilters = (params) => {
  const newQueryString = params.toString();
  window.location.replace(`${window.location.pathname}?${newQueryString}`);
};

export const updateFilters = () => {
  const filters = getFilters();

  const selectedFilters = Array.from(filters).filter(
    (filter) => filter.checked
  );
  const selectedFiltersValues = selectedFilters.map((filter) => filter.value);

  const params = new URLSearchParams(window.location.search);

  params.delete("category");

  if (selectedFiltersValues.length > 0) {
    for (const filter of selectedFiltersValues) {
      params.append("category", filter);
    }
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

  updateUrlFilters(params);
};

const removeFilter = (filterValue) => {
  const filters = getFilters();

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
