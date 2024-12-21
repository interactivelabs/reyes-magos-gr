// assets/js/shared/toasts.ts
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

// assets/js/shared/htmxErrorHandler.ts
window.addEventListener("htmx:responseError", (e) => {
  console.log(e);
  const code = e.detail.xhr.status;
  if (code === 500) {
    showToast({
      variant: "error",
      title: "Error del servidor",
      message: "Ha ocurrido un error inesperado. Por favor intenta de nuevo."
    });
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
