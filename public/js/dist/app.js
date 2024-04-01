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
  const toggleMobileMenu = (evt) => {
    evt.stopPropagation();
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
  document.addEventListener("click", (event) => {
    event.stopPropagation();
    if (mobileMenuContainer && mobileMenuButton) {
      closeIfOutsideClick({
        element: mobileMenuContainer,
        elementButton: mobileMenuButton,
        event
      });
    }
  });
}

// public/js/src/app.ts
initNav();
//# sourceMappingURL=app.js.map
