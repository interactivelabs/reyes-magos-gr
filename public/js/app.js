// assets/js/shared/nav.ustils.ts
var closeIfOutsideClick = ({
  element,
  elementButton,
  event,
  onClose
}) => {
  const isClickInside = element.contains(event.target) || element === event.target || elementButton.contains(event.target) || elementButton === event.target;
  if (!isClickInside && element.getAttribute("data-transition-state") === "open") {
    element.setAttribute("data-transition-state", "closed");
    onClose && onClose();
  }
};
var toggleMenu = (menuContainer) => {
  if (!menuContainer) return;
  if (menuContainer.getAttribute("data-transition-state") === "closed") {
    menuContainer.setAttribute("open", "true");
    menuContainer.setAttribute("data-transition-state", "open");
  } else {
    menuContainer.removeAttribute("open");
    menuContainer.setAttribute("data-transition-state", "closed");
  }
};

// assets/js/admin/admin.nav.ts
function initAdminNav() {
  const adminMenuDropdown = document.getElementById("admin-menu-dropdown");
  const adminMenuButton = document.getElementById("admin-menu-button");
  const toggleAdminMenu = (evt) => {
    evt.stopPropagation();
    toggleMenu(adminMenuDropdown);
  };
  adminMenuButton?.addEventListener("click", toggleAdminMenu);
  document.addEventListener("click", (event) => {
    event.stopPropagation();
    if (adminMenuDropdown && adminMenuButton) {
      closeIfOutsideClick({
        element: adminMenuDropdown,
        elementButton: adminMenuButton,
        event
      });
    }
  });
}

// assets/js/app/nav.ts
function initNav() {
  const mobileMenuContainer = document.getElementById("mobile-menu-container");
  const mobileMenuButton = document.getElementById("mobile-menu-button");
  const mobileMenuButtonIconClosed = document.getElementById(
    "mobile-menu-button-icon-closed"
  );
  const mobileMenuButtonIconOpen = document.getElementById(
    "mobile-menu-button-icon-open"
  );
  const toggleMobileMenuButtonIcon = () => {
    mobileMenuButtonIconClosed?.classList.toggle("hidden");
    mobileMenuButtonIconOpen?.classList.toggle("hidden");
  };
  const toggleMobileMenu = (evt) => {
    evt.stopPropagation();
    toggleMenu(mobileMenuContainer);
    toggleMobileMenuButtonIcon();
  };
  mobileMenuButton?.addEventListener("click", toggleMobileMenu);
  document.addEventListener("click", (event) => {
    event.stopPropagation();
    if (mobileMenuContainer && mobileMenuButton) {
      closeIfOutsideClick({
        event,
        element: mobileMenuContainer,
        elementButton: mobileMenuButton,
        onClose: toggleMobileMenuButtonIcon
      });
    }
  });
}

// assets/js/shared/toast.ts
var hideToast = () => {
  const toastContainer = document.getElementById("toast-container");
  const toastPanel = document.getElementById("toast-panel");
  toastPanel?.classList.remove("show-toast");
  toastPanel?.classList.add("hide-toast");
  setTimeout(() => {
    toastContainer?.classList.remove("flex");
    toastContainer?.classList.add("hidden");
  }, 300);
};
var showToast = ({ title, subTitle, duration = 5e3 }) => {
  const toastContainer = document.getElementById("toast-container");
  const toastPanel = document.getElementById("toast-panel");
  toastContainer?.classList.remove("hidden");
  toastContainer?.classList.add("flex");
  toastPanel?.classList.remove("hide-toast");
  toastPanel?.classList.add("show-toast");
  const toastTitle = document.getElementById("toast-title");
  const toastSubTitle = document.getElementById("toast-subtitle");
  toastTitle.innerText = title;
  toastSubTitle.innerText = subTitle;
  setTimeout(() => hideToast(), duration);
};
function initToast() {
  globalThis.hideToast = hideToast;
}

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
function initMycodes() {
  globalThis.copyCode = copy;
  globalThis.shareCode = async (code) => {
    if (typeof navigator.share === "undefined") {
      copy(code);
    } else {
      share(code);
    }
  };
}

// assets/js/tailwind/tailwind-core.ts
var getStateClasses = (element, attribute) => {
  const stateClasses = element.getAttribute(attribute);
  if (!stateClasses) {
    throw new Error(`Missing ${attribute} attribute`);
  }
  const stateClassesArray = stateClasses.split(" ");
  const afterClasses = stateClassesArray.filter(
    (className) => className.startsWith("final:")
  );
  const immediateClasses = stateClassesArray.filter(
    (className) => !className.startsWith("final:")
  );
  return [immediateClasses, afterClasses];
};
var updateStateClasses = (element, attribute, action) => {
  const [immediateClasses, afterClasses] = getStateClasses(element, attribute);
  const createTimeoutToAddClass = (delay, className) => {
    setTimeout(() => {
      element.classList.add(className);
    }, delay);
  };
  immediateClasses.forEach((className) => {
    if (action === "add" /* ADD */) {
      element.classList.add(className);
      createTimeoutToAddClass(10, className);
    } else {
      element.classList.remove(className);
    }
  });
  afterClasses.forEach((className) => {
    const props = className.split("final:")[1];
    const [delay, afterClassName] = JSON.parse(props.replace(/'/g, '"'));
    if (action === "add" /* ADD */) {
      createTimeoutToAddClass(delay, afterClassName);
    } else {
      element.classList.remove(afterClassName);
    }
  });
};
var addStateClasses = (element, attribute) => {
  updateStateClasses(element, attribute, "add" /* ADD */);
};
var removeStateClasses = (element, attribute) => {
  updateStateClasses(element, attribute, "remove" /* REMOVE */);
};
var getDelayFromAttribute = (element, attribute) => {
  const classNames = element.getAttribute(attribute);
  if (!classNames) {
    throw new Error(`Missing ${attribute} attribute`);
  }
  const classNamesArray = classNames.split(" ");
  const delayClasses = classNamesArray.filter(
    (className) => className.startsWith("delay-") || className.startsWith("duration-")
  );
  if (!delayClasses || !delayClasses.length) return 150;
  let delay = 0;
  delayClasses.forEach((delayClass) => {
    const delayValue = parseInt(delayClass.split("-")[1], 10);
    delay += delayValue;
  });
  return delay;
};

// assets/js/tailwind/tailwind-transitions.ts
function initTailwindTransitions() {
  const ATTRIBUTE = "data-transition-state";
  const hideElement = (element) => {
    element.style.display = "none";
  };
  const showElement = (element) => {
    element.style.display = "block";
    element.offsetHeight;
  };
  const transitionMutationCallback = (mutationsList) => {
    for (const mutation of mutationsList) {
      const { type, attributeName, oldValue, target } = mutation;
      if (type !== "attributes" || attributeName !== ATTRIBUTE || !target) {
        return;
      }
      const state = target.getAttribute(ATTRIBUTE);
      addStateClasses(target, `data-transition-${state}`);
      if (state === "open") {
        showElement(target);
      } else {
        const delay = getDelayFromAttribute(
          target,
          `data-transition-${oldValue}`
        );
        setTimeout(() => hideElement(target), delay);
      }
      if (oldValue) {
        removeStateClasses(target, `data-transition-${oldValue}`);
      }
    }
  };
  const observer = new MutationObserver(transitionMutationCallback);
  document.addEventListener("DOMContentLoaded", function() {
    const allElements = document.querySelectorAll(`[${ATTRIBUTE}]`);
    allElements.forEach((transitionElement) => {
      const state = transitionElement.getAttribute(ATTRIBUTE);
      if (state !== "closed" && state !== "open") {
        throw new Error(
          `Invalid ${state} for transitions, use data-styles-[state] for non open/close transitions`
        );
      }
      if (state === "closed") {
        hideElement(transitionElement);
      } else {
        showElement(transitionElement);
      }
      addStateClasses(transitionElement, `data-transition-${state}`);
      observer.observe(transitionElement, {
        attributes: true,
        attributeOldValue: true,
        attributeFilter: [ATTRIBUTE]
      });
    });
  });
}

// assets/js/app/catalogFilters.ts
var getFilters = () => document.getElementsByName("category[]");
var updateUrlFilters = (params) => {
  const newQueryString = params.toString();
  window.location.replace(`${window.location.pathname}?${newQueryString}`);
};
var updateFilters = () => {
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
var clearFilters = () => {
  const filters = getFilters();
  for (const filter of filters) {
    filter.checked = false;
  }
  const params = new URLSearchParams(window.location.search);
  params.delete("category");
  updateUrlFilters(params);
};
var removeFilter = (filterValue) => {
  const filters = getFilters();
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
function initCatalogFilters() {
  globalThis.updateFilters = updateFilters;
  globalThis.clearFilters = clearFilters;
  globalThis.removeFilter = removeFilter;
  setUrlFilters();
}

// assets/js/app.ts
initTailwindTransitions();
initMycodes();
initAdminNav();
initNav();
initToast();
initCatalogFilters();
//# sourceMappingURL=app.js.map
