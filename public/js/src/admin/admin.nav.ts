import { closeIfOutsideClick } from "../shared/nav.ustils";

export default function initAdminNav() {
  const adminMenuDropdown = document.getElementById("admin-menu-dropdown");
  const adminMenuButton = document.getElementById("admin-menu-button");

  const toggleAdminMenu = (evt: MouseEvent) => {
    evt.stopPropagation();
    if (!adminMenuDropdown) return;
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
        event: event,
      });
    }
  });
}
