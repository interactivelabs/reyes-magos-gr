// assets/js/shared/nav.ustils.ts
var closeIfOutsideClick = ({
  element,
  elementButton,
  event,
  onClose
}) => {
  const isClickInside = element.contains(event.target) || element === event.target || elementButton.contains(event.target) || elementButton === event.target;
  if (!isClickInside && !element.classList.contains("hidden")) {
    element.classList.add("hidden");
    onClose && onClose();
  }
};
var toggleMenu = (menuContainer) => {
  if (!menuContainer) return;
  if (menuContainer.classList.contains("hidden")) {
    menuContainer.setAttribute("open", "true");
  } else {
    menuContainer.removeAttribute("open");
  }
  menuContainer.classList.toggle("hidden");
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

// assets/js/app.ts
initMycodes();
initAdminNav();
initNav();
initToast();
//# sourceMappingURL=app.js.map
