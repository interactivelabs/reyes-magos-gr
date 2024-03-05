// public/js/nav.ts
function initNav() {
  const mobileMenuContainer = document.getElementById("mobile-menu-container");
  const mobileMenuButton = document.getElementById("mobile-menu-button");
  const mobileMenuButtonIconClosed = document.getElementById(
    "mobile-menu-button-icon-closed"
  );
  const mobileMenuButtonIconOpen = document.getElementById(
    "mobile-menu-button-icon-open"
  );
  const toggleMobileMenu = () => {
    if (!mobileMenuContainer)
      return;
    if (mobileMenuContainer.classList.contains("hidden")) {
      mobileMenuContainer.setAttribute("open", "true");
    } else {
      mobileMenuContainer.removeAttribute("open");
    }
    mobileMenuContainer.classList.toggle("hidden");
    mobileMenuButtonIconClosed?.classList.toggle("hidden");
    mobileMenuButtonIconOpen?.classList.toggle("hidden");
  };
  mobileMenuButton?.addEventListener("click", toggleMobileMenu);
  const adminMenuDropdown = document.getElementById("admin-menu-dropdown");
  const adminMenuButton = document.getElementById("admin-menu-button");
  const toggleAdminMenu = () => {
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
  const closeIfOutsideClick = ({ element, elementButton, event }) => {
    const isClickInside = element.contains(event.currentTarget) || elementButton.contains(event.currentTarget);
    if (!isClickInside) {
      element.classList.add("hidden");
    }
  };
  document.addEventListener("click", (event) => {
    closeIfOutsideClick({ element: adminMenuDropdown, elementButton: adminMenuButton, event });
    closeIfOutsideClick({ element: mobileMenuContainer, elementButton: mobileMenuButton, event });
  });
}

// public/js/app.ts
initNav();
//# sourceMappingURL=app.js.map
