export const closeIfOutsideClick = ({
  element,
  elementButton,
  event,
}: {
  element: HTMLElement;
  elementButton: HTMLElement;
  event: MouseEvent;
}) => {
  const isClickInside =
    element.contains(event.target as HTMLElement) ||
    element === event.target ||
    elementButton.contains(event.target as HTMLElement) ||
    elementButton === event.target;
  if (!isClickInside) {
    element.classList.add("hidden");
  }
};
