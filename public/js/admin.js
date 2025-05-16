// assets/js/components/search-box.ts
function getSearchBoxClass() {
  return class SearchBox {
    // MAGIC METHODS
    // Alpine magic properties and methods
    $data;
    $dispatch;
    $el;
    $id;
    $nextTick;
    $refs;
    $root;
    $store;
    $watch;
    destroy;
    // CLASS PROPERTIES
    name;
    Id;
    Items;
    ItemsFiltered;
    ItemActive;
    ItemSelected;
    Search;
    fetchItemsUrl;
    // Constructor, initialize properties
    constructor(url, name) {
      this.Id = name + Date.now().toString();
      this.name = name;
      this.Items = [];
      this.ItemsFiltered = [];
      this.ItemActive = null;
      this.ItemSelected = null;
      this.Search = "";
      this.fetchItemsUrl = url;
    }
    // Alphine lifecycle methods
    async init() {
      this.$watch("Search", () => this.SearchItems());
      this.$watch(
        "ItemSelected",
        (item) => this.SelectItem(item)
      );
      this.Items = await (await fetch(this.fetchItemsUrl)).json();
    }
    // CLASS METHODS
    SearchIsEmpty() {
      return this.Search.length == 0;
    }
    ItemIsActive(item) {
      return !!(this.ItemActive && this.ItemActive.Value == item.Value);
    }
    ClearState() {
      this.ItemsFiltered = [];
      this.ItemActive = null;
      this.Search = "";
    }
    ItemActiveNext() {
      if (!this.ItemActive) return;
      let index = this.ItemsFiltered.indexOf(this.ItemActive);
      if (index < this.ItemsFiltered.length - 1) {
        this.ItemActive = this.ItemsFiltered[index + 1];
        this.ScrollToActiveItem();
      }
    }
    ItemActivePrevious() {
      if (!this.ItemActive) return;
      let index = this.ItemsFiltered.indexOf(this.ItemActive);
      if (index > 0) {
        this.ItemActive = this.ItemsFiltered[index - 1];
        this.ScrollToActiveItem();
      }
    }
    ScrollToActiveItem() {
      if (this.ItemActive) {
        const activeElement = document.getElementById(
          this.ItemActive.Value + "-" + this.Id
        );
        if (!activeElement) return;
        const newScrollPos = activeElement.offsetTop + activeElement.offsetHeight - this.$refs.ItemsList.offsetHeight;
        if (newScrollPos > 0) {
          this.$refs.ItemsList.scrollTop = newScrollPos;
        } else {
          this.$refs.ItemsList.scrollTop = 0;
        }
      }
    }
    SearchItems() {
      if (this.SearchIsEmpty()) {
        this.ClearState();
        return;
      }
      if (!this.SearchIsEmpty()) {
        const searchTerm = this.Search.replace(/\*/g, "").toLowerCase();
        this.ItemsFiltered = this.Items.filter(
          (item) => item.Label.toLowerCase().includes(searchTerm)
        );
        this.ScrollToActiveItem();
      }
    }
    SelectItem(item) {
      if (!item) return;
      this.$dispatch(`searchbox-item-selected-${this.name}`, { item });
      this.ClearState();
    }
    AddItem() {
      if (this.SearchIsEmpty()) return;
      const item = { Label: this.Search, Value: this.Search };
      this.Items.push(item);
      this.$dispatch(`searchbox-item-selected-${this.name}`, { item });
      this.ClearState();
    }
  };
}
document.addEventListener("alpine:init", () => {
  const createSearchBox = (url, name) => {
    const SearchBox = getSearchBoxClass();
    return new SearchBox(url, name);
  };
  Alpine.data("SearchBox", createSearchBox);
});

// assets/js/admin/codes.ts
var form = document.getElementById("remove_code_form");
form?.addEventListener("change", (e) => {
  const target = e.target;
  if (target.name === "volunteer_code_ids") {
    const codeCheckbox = target.nextElementSibling;
    codeCheckbox.checked = target.checked;
  }
});
globalThis.selectAllUnsigned = () => {
  const assignCodesForm = document.getElementById("assign-codes-form");
  const checkboxes = assignCodesForm?.getElementsByTagName("input");
  if (!checkboxes || !checkboxes.length) return;
  for (let i = 0; i < checkboxes?.length; i++) {
    if (checkboxes[i].type === "checkbox") {
      checkboxes[i].checked = true;
    }
  }
};

// assets/js/admin/orders.ts
globalThis.orderCompletedToggleClick = (btn) => {
  btn.classList.toggle("bg-indigo-600");
  btn.classList.toggle("bg-gray-200");
  btn.children[0].classList.toggle("translate-x-0");
  btn.children[0].classList.toggle("translate-x-5");
  const completedInput = btn.previousElementSibling;
  completedInput.value = completedInput.value === "1" ? "0" : "1";
};

// assets/js/admin/toys_form.ts
window.addEventListener("SearchBox:ItemSelected:category", (event) => {
  const searchBoxEvent = event;
  const item = searchBoxEvent.detail?.item;
  if (!item) return;
  const categoryInput = document.getElementById("category");
  const category = item.Value;
  const currentCategories = categoryInput.value.split(",");
  if (currentCategories.includes(category)) return;
  const newCategories = [...currentCategories, category];
  categoryInput.value = newCategories.join(",");
});
document.addEventListener("alpine:init", () => {
  Alpine.data("CategoryWithSearch", (currentValue) => ({
    category: currentValue,
    setNewCategory(event) {
      const item = event.detail?.item;
      if (!item) return;
      const newCategory = item.Value;
      const currentCategories = !!this.category ? this.category.split(",") : [];
      if (currentCategories.includes(newCategory)) return;
      const newCategories = [...currentCategories, newCategory];
      this.category = newCategories.join(",");
    },
    removeCategory(category) {
      const newCategories = this.category.split(",").filter((c) => c !== category);
      if (newCategories.length > 0) {
        this.category = newCategories.join(",");
        return;
      }
      this.category = "";
    }
  }));
});
//# sourceMappingURL=admin.js.map
