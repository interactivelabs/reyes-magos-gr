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
  if (
    !isClickInside &&
    element.getAttribute("data-transition-state") === "open"
  ) {
    element.setAttribute("data-transition-state", "closed");
    onClose && onClose();
  }
};

export const toggleMenu = (menuContainer: HTMLElement | null) => {
  if (!menuContainer) return;
  if (menuContainer.getAttribute("data-transition-state") === "closed") {
    menuContainer.setAttribute("open", "true");
    menuContainer.setAttribute("data-transition-state", "open");
  } else {
    menuContainer.removeAttribute("open");
    menuContainer.setAttribute("data-transition-state", "closed");
  }
};
