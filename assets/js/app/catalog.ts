import debounce from "../utils/debounce";

declare global {
  var updateFilters: () => void;
  var clearFilters: () => void;
  var removeFilter: (filterValue: string) => void;
}

const getFilters = () =>
  document.getElementsByName("category[]") as NodeListOf<HTMLInputElement>;

const updateUrlFilters = (params: URLSearchParams) => {
  params.delete("page");
  const newQueryString = params.toString();
  window.location.replace(`${window.location.pathname}?${newQueryString}`);
};

const getAgeInputs = () => ({
  ageMin: document.getElementsByName("age_min")[0] as
    | HTMLInputElement
    | undefined,
  ageMax: document.getElementsByName("age_max")[0] as
    | HTMLInputElement
    | undefined,
});

const getAgeFilters = () => {
  const { ageMin, ageMax } = getAgeInputs();
  return { ageMin: ageMin?.value ?? "0", ageMax: ageMax?.value ?? "0" };
};

// age_max at 1 means "no max", so only clamp when a max is actually set
const syncAgeRange = () => {
  const { ageMin, ageMax } = getAgeInputs();
  if (!ageMin || !ageMax) return;

  const min = parseInt(ageMin.value, 10);
  const max = parseInt(ageMax.value, 10);
  if (max > 1 && max < min) {
    ageMax.value = ageMin.value;
  }
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

  syncAgeRange();
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

const removeFilter = (filterValue: string) => {
  const filters = getFilters();

  const { ageMin, ageMax } = getAgeInputs();
  if (filterValue === "age_min" && ageMin) {
    ageMin.value = "0";
  } else if (filterValue === "age_max" && ageMax) {
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

globalThis.updateFilters = updateFilters;
globalThis.clearFilters = clearFilters;
globalThis.removeFilter = removeFilter;
setUrlFilters();

const ageInputs = getAgeInputs();
ageInputs.ageMin?.addEventListener("input", syncAgeRange);
ageInputs.ageMax?.addEventListener("input", syncAgeRange);

// Back to top button display
const backToTopBtn = document.getElementById("back-to-top-btn");

const handleBackToTopScroll = () => {
  const isScrolled =
    document.body.scrollTop > 300 || document.documentElement.scrollTop > 300;
  backToTopBtn?.classList.toggle("opacity-0", !isScrolled);
  backToTopBtn?.classList.toggle("pointer-events-none", !isScrolled);
  backToTopBtn?.classList.toggle("opacity-100", isScrolled);
  backToTopBtn?.toggleAttribute("inert", !isScrolled);
};

const scrollHandler = debounce(handleBackToTopScroll, 250);

window.addEventListener("scroll", scrollHandler);
