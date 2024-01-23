export default function initNav() {
  const mobileMenuContainer = document.getElementById("mobile-menu-container");
  const mobileMenuButton = document.getElementById("mobile-menu-button");
  const mobileMenuButtonIconClosed = document.getElementById(
    "mobile-menu-button-icon-closed"
  );
  const mobileMenuButtonIconOpen = document.getElementById(
    "mobile-menu-button-icon-open"
  );
  const toggleMobileMenu = () => {
    if (!mobileMenuContainer) return;
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
    if (!adminMenuDropdown) return;
    if (adminMenuDropdown.classList.contains("hidden")) {
      adminMenuDropdown.setAttribute("open", "true");
    } else {
      adminMenuDropdown.removeAttribute("open");
    }
    adminMenuDropdown.classList.toggle("hidden");
  };

  adminMenuButton?.addEventListener("click", toggleAdminMenu);

  const closeIfOutsideClick = (element, elementButton, event) => {
    var isClickInside =
      element.contains(event.target) || elementButton.contains(event.target);

    if (!isClickInside) {
      element.classList.add("hidden");
    }
  };

  document.addEventListener("click", function (event) {
    closeIfOutsideClick(adminMenuDropdown, adminMenuButton, event);
    closeIfOutsideClick(mobileMenuContainer, mobileMenuButton, event);
  });
}
