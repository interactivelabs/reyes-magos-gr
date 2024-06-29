import { closeIfOutsideClick, toggleMenu } from "../shared/nav.ustils";

export default function initNav() {
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

  const toggleMobileMenu = (evt: MouseEvent) => {
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
        onClose: toggleMobileMenuButtonIcon,
      });
    }
  });
}
