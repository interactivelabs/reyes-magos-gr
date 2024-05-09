// public/js/src/shared/nav.ustils.ts
var closeIfOutsideClick = ({
  element,
  elementButton,
  event,
  onClose
}) => {
  const isClickInside = element.contains(event.target) || element === event.target || elementButton.contains(event.target) || elementButton === event.target;
  if (!isClickInside) {
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

// public/js/src/admin/codes.ts
function initCodes() {
  const form = document.getElementById("remove_code_form");
  form?.addEventListener("change", (e) => {
    const target = e.target;
    if (target.name === "volunteer_code_ids") {
      const codeCheckbox = target.nextElementSibling;
      codeCheckbox.checked = target.checked;
    }
  });
}

// public/js/src/admin.ts
initAdminNav();
initCodes();
//# sourceMappingURL=admin.js.map
