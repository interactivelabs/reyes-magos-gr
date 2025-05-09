// assets/js/components/toasts.ts
document.addEventListener("alpine:init", () => {
  Alpine.data("toasts", () => ({
    notifications: [],
    displayDuration: 5e3,
    addNotification(notification) {
      const id = Date.now();
      if (this.notifications.length >= 20) {
        this.notifications.splice(0, this.notifications.length - 19);
      }
      this.notifications.push({ ...notification, id });
    },
    removeNotification(id) {
      setTimeout(() => {
        this.notifications = this.notifications.filter(
          (notification) => notification.id !== id
        );
      }, 400);
    }
  }));
});
function showToast(toast) {
  window.dispatchEvent(new CustomEvent("notify", { detail: toast }));
}

// assets/js/components/search-box.ts
function getSearchBoxClass() {
  return class SearchBox {
    // MAGIC METHODS
    // Alpine magic properties and methods
    $data;
    $dispatch;
    $el;
    $id;
    $nextTick;
    $refs;
    $root;
    $store;
    $watch;
    destroy;
    // CLASS PROPERTIES
    name;
    Id;
    Items;
    ItemsFiltered;
    ItemActive;
    ItemSelected;
    Search;
    fetchItemsUrl;
    // Constructor, initialize properties
    constructor(url, name) {
      this.Id = name + Date.now().toString();
      this.name = name;
      this.Items = [];
      this.ItemsFiltered = [];
      this.ItemActive = null;
      this.ItemSelected = null;
      this.Search = "";
      this.fetchItemsUrl = url;
    }
    // Alphine lifecycle methods
    async init() {
      this.$watch("Search", () => this.SearchItems());
      this.$watch(
        "ItemSelected",
        (item) => this.SelectItem(item)
      );
      this.Items = await (await fetch(this.fetchItemsUrl)).json();
    }
    // CLASS METHODS
    SearchIsEmpty() {
      return this.Search.length == 0;
    }
    ItemIsActive(item) {
      return !!(this.ItemActive && this.ItemActive.Value == item.Value);
    }
    ClearState() {
      this.ItemsFiltered = [];
      this.ItemActive = null;
      this.Search = "";
    }
    ItemActiveNext() {
      if (!this.ItemActive) return;
      let index = this.ItemsFiltered.indexOf(this.ItemActive);
      if (index < this.ItemsFiltered.length - 1) {
        this.ItemActive = this.ItemsFiltered[index + 1];
        this.ScrollToActiveItem();
      }
    }
    ItemActivePrevious() {
      if (!this.ItemActive) return;
      let index = this.ItemsFiltered.indexOf(this.ItemActive);
      if (index > 0) {
        this.ItemActive = this.ItemsFiltered[index - 1];
        this.ScrollToActiveItem();
      }
    }
    ScrollToActiveItem() {
      if (this.ItemActive) {
        const activeElement = document.getElementById(
          this.ItemActive.Value + "-" + this.Id
        );
        if (!activeElement) return;
        const newScrollPos = activeElement.offsetTop + activeElement.offsetHeight - this.$refs.ItemsList.offsetHeight;
        if (newScrollPos > 0) {
          this.$refs.ItemsList.scrollTop = newScrollPos;
        } else {
          this.$refs.ItemsList.scrollTop = 0;
        }
      }
    }
    SearchItems() {
      if (this.SearchIsEmpty()) {
        this.ClearState();
        return;
      }
      if (!this.SearchIsEmpty()) {
        const searchTerm = this.Search.replace(/\*/g, "").toLowerCase();
        this.ItemsFiltered = this.Items.filter(
          (item) => item.Label.toLowerCase().includes(searchTerm)
        );
        this.ScrollToActiveItem();
      }
    }
    SelectItem(item) {
      if (!item) return;
      this.$dispatch(`searchbox-item-selected-${this.name}`, { item });
      this.ClearState();
    }
    AddItem() {
      if (this.SearchIsEmpty()) return;
      const item = { Label: this.Search, Value: this.Search };
      this.Items.push(item);
      this.$dispatch(`searchbox-item-selected-${this.name}`, { item });
      this.ClearState();
    }
  };
}
document.addEventListener("alpine:init", () => {
  const createSearchBox = (url, name) => {
    const SearchBox = getSearchBoxClass();
    return new SearchBox(url, name);
  };
  Alpine.data("SearchBox", createSearchBox);
});

// assets/js/shared/htmxErrorHandler.ts
window.addEventListener("htmx:responseError", (e) => {
  const code = e.detail.xhr.status;
  if (code === 500) {
    showToast({
      variant: "error",
      title: "Error del servidor",
      message: "Ha ocurrido un error inesperado. Por favor intenta de nuevo."
    });
  }
  if (code === 400) {
    console.error(e);
  }
});

// assets/js/app/myCodes.ts
var copy = async (code) => {
  try {
    await navigator.clipboard.writeText(code);
    showToast({
      variant: "error",
      title: "Copiado!",
      message: "El codigo ha sido copiado al portapapeles."
    });
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};
var share = async (code) => {
  const data = {
    title: "Comparte la alegria!",
    text: `Utiliza este codigo para obtener un regalo: ${code}`,
    url: `${window.location.origin}/catalog?code=${code}`
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

// assets/js/utils/debounce.ts
function debounce(func, wait) {
  let timeout;
  return function(...args) {
    if (timeout) {
      clearTimeout(timeout);
    }
    timeout = setTimeout(() => {
      func.apply(this, args);
    }, wait);
  };
}
var debounce_default = debounce;

// assets/js/app/catalog.ts
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
var handleBackToTopScroll = () => {
  if (document.body.scrollTop > 300 || document.documentElement.scrollTop > 300) {
    window.dispatchEvent(new CustomEvent("showbacktotopbtn"));
  } else {
    window.dispatchEvent(new CustomEvent("hidebacktotopbtn"));
  }
};
var scrollHandler = debounce_default(handleBackToTopScroll, 250);
window.addEventListener("scroll", scrollHandler);

// assets/js/app/homeAnnimations.ts
if (typeof gsap !== "undefined") {
  gsap.registerPlugin(ScrollTrigger);
  for (let i = 1; i <= 3; i++) {
    gsap.to(`#home-ilustracion-${i}`, {
      scrollTrigger: {
        trigger: `#home-ilustracion-${i}`,
        toggleActions: "play pause resume reset",
        start: "top 80%"
      },
      duration: 1,
      opacity: 1
    });
  }
}

// assets/js/app/redeem.ts
function validateCodeInput() {
  const codeInput = document.getElementById("code");
  const code = codeInput.value.toUpperCase().trim();
  return code.length > 5;
}
function beforeRequestHandler(evt) {
  if (!validateCodeInput()) {
    evt.preventDefault();
  }
}
document.addEventListener("alpine:init", () => {
  Alpine.data("codeForm", () => ({
    valid: true,
    init() {
      this.valid = true;
    },
    submit(evt) {
      this.valid = validateCodeInput();
      evt.preventDefault();
      evt.stopPropagation();
    },
    update() {
      const codeInput = document.getElementById("code");
      const code = codeInput.value.toUpperCase().trim();
      codeInput.value = code;
      if (code.length > 5) {
        this.valid = true;
      }
    }
  }));
});
globalThis.beforeRequestHandler = beforeRequestHandler;
//# sourceMappingURL=app.js.map
