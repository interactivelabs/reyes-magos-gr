import type { SearchBoxSelectEvent } from "../components/search-box";

window.addEventListener("SearchBox:ItemSelected:category", (event: Event) => {
  const searchBoxEvent = event as SearchBoxSelectEvent;
  const item = searchBoxEvent.detail?.item;
  if (!item) return;

  const categoryInput = document.getElementById("category") as HTMLInputElement;
  const category = item.Value;

  const currentCategories = categoryInput.value.split(",");
  if (currentCategories.includes(category)) return;

  const newCategories = [...currentCategories, category];
  categoryInput.value = newCategories.join(",");
});

document.addEventListener("alpine:init", () => {
  Alpine.data("CategoryWithSearch", (currentValue: string) => ({
    category: currentValue,
    setNewCategory(event: SearchBoxSelectEvent) {
      const item = event.detail?.item;
      if (!item) return;

      const newCategory = item.Value;
      const currentCategories = !!this.category ? this.category.split(",") : [];
      if (currentCategories.includes(newCategory)) return;

      const newCategories = [...currentCategories, newCategory];
      this.category = newCategories.join(",");
    },
    removeCategory(category: string) {
      const newCategories = this.category
        .split(",")
        .filter((c: string) => c !== category);
      if (newCategories.length > 0) {
        this.category = newCategories.join(",");
        return;
      }
      this.category = "";
    },
  }));
});
