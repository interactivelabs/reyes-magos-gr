// public/js/src/shared/nav.ustils.ts
var closeIfOutsideClick = ({
  element,
  elementButton,
  event
}) => {
  const isClickInside = element.contains(event.target) || element === event.target || elementButton.contains(event.target) || elementButton === event.target;
  if (!isClickInside) {
    element.classList.add("hidden");
  }
};

// public/js/src/admin/admin.nav.ts
function initAdminNav() {
  const adminMenuDropdown = document.getElementById("admin-menu-dropdown");
  const adminMenuButton = document.getElementById("admin-menu-button");
  const toggleAdminMenu = (evt) => {
    evt.stopPropagation();
    if (!adminMenuDropdown)
      return;
    if (adminMenuDropdown.classList.contains("hidden")) {
      adminMenuDropdown.setAttribute("open", "true");
    } else {
      adminMenuDropdown.removeAttribute("open");
    }
    adminMenuDropdown.classList.toggle("hidden");
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
