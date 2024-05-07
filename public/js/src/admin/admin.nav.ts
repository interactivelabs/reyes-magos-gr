import { closeIfOutsideClick, toggleMenu } from "../shared/nav.ustils";

export default function initAdminNav() {
  const adminMenuDropdown = document.getElementById("admin-menu-dropdown");
  const adminMenuButton = document.getElementById("admin-menu-button");

  const toggleAdminMenu = (evt: MouseEvent) => {
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
        event,
      });
    }
  });
}
