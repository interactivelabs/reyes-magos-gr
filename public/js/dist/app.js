// public/js/src/shared/nav.ustils.ts
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

// public/js/src/admin/admin.nav.ts
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

// public/js/src/app/nav.ts
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

// public/js/src/app/myCodes.ts
var copyCode = async (code) => {
  try {
    await navigator.clipboard.writeText(code);
  } catch (err) {
    console.error("Failed to copy: ", err);
  }
};
function initMycodes() {
  const codesList = document.getElementById("mycodes-code-list");
  if (!codesList) {
    return;
  }
  codesList.addEventListener("click", async (evt) => {
    let target = evt.target;
    if (!target || target.tagName === "UL") return;
    while (target?.tagName !== "LI") {
      target = target?.parentElement;
    }
    const code = target.id;
    await copyCode(code);
    const toast = document.getElementById(`mycodes-copied-label-${code}`);
    if (toast) {
      toast.classList.remove("hidden");
      setTimeout(() => {
        toast.classList.add("hidden");
      }, 1500);
    }
  });
}

// public/js/src/app.ts
initMycodes();
initAdminNav();
initNav();
//# sourceMappingURL=app.js.map
