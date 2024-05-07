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
  if (!menuContainer)
    return;
  if (menuContainer.classList.contains("hidden")) {
    menuContainer.setAttribute("open", "true");
  } else {
    menuContainer.removeAttribute("open");
  }
  menuContainer.classList.toggle("hidden");
};

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

// public/js/src/app.ts
initNav();
//# sourceMappingURL=app.js.map
