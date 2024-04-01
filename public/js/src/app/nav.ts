import { closeIfOutsideClick } from "../shared/nav.ustils";

export default function initNav() {
  const mobileMenuContainer = document.getElementById("mobile-menu-container");
  const mobileMenuButton = document.getElementById("mobile-menu-button");
  const mobileMenuButtonIconClosed = document.getElementById(
    "mobile-menu-button-icon-closed"
  );
  const mobileMenuButtonIconOpen = document.getElementById(
    "mobile-menu-button-icon-open"
  );

  const toggleMobileMenu = (evt: MouseEvent) => {
    evt.stopPropagation();
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

  document.addEventListener("click", (event) => {
    event.stopPropagation();
    if (mobileMenuContainer && mobileMenuButton) {
      closeIfOutsideClick({
        element: mobileMenuContainer,
        elementButton: mobileMenuButton,
        event: event,
      });
    }
  });
}
