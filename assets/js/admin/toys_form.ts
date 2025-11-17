document.addEventListener("alpine:init", () => {
  Alpine.data("CategoryWithSearch", (initalValue: string) => ({
    category: initalValue,
    category_search: "",
    setNewCategory() {
      if (!this.category_search) return;

      const newCategory = this.category_search;
      const currentCategories = !!this.category ? this.category.split(",") : [];
      if (currentCategories.includes(newCategory)) return;

      const newCategories = [...currentCategories, newCategory];
      this.category = newCategories.join(",");
      this.category_search = "";
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
