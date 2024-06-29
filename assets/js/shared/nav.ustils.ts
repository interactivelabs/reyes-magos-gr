export interface ICloseIfOutsideClick {
  element: HTMLElement;
  elementButton: HTMLElement;
  event: MouseEvent;
  onClose?: () => void;
}

export const closeIfOutsideClick = ({
  element,
  elementButton,
  event,
  onClose,
}: ICloseIfOutsideClick) => {
  const isClickInside =
    element.contains(event.target as HTMLElement) ||
    element === event.target ||
    elementButton.contains(event.target as HTMLElement) ||
    elementButton === event.target;
  if (!isClickInside && !element.classList.contains("hidden")) {
    element.classList.add("hidden");
    onClose && onClose();
  }
};

export const toggleMenu = (menuContainer: HTMLElement | null) => {
  if (!menuContainer) return;
  if (menuContainer.classList.contains("hidden")) {
    menuContainer.setAttribute("open", "true");
  } else {
    menuContainer.removeAttribute("open");
  }
  menuContainer.classList.toggle("hidden");
};
