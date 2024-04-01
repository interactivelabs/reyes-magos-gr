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

  const closeIfOutsideClick = ({
    element,
    elementButton,
    event,
  }: {
    element: HTMLElement;
    elementButton: HTMLElement;
    event: MouseEvent;
  }) => {
    const isClickInside =
      element.contains(event.currentTarget as HTMLElement) ||
      elementButton.contains(event.currentTarget as HTMLElement);
    if (!isClickInside) {
      element.classList.add("hidden");
    }
  };

  document.addEventListener("click", (event) => {
    if (adminMenuDropdown && adminMenuButton) {
      closeIfOutsideClick({
        element: adminMenuDropdown,
        elementButton: adminMenuButton,
        event: event,
      });
    }
    if (mobileMenuContainer && mobileMenuButton) {
      closeIfOutsideClick({
        element: mobileMenuContainer,
        elementButton: mobileMenuButton,
        event: event,
      });
    }
  });
}
