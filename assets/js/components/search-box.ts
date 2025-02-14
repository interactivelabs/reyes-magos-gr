/// <reference types="alpinejs" />
import type { Alpine as AlpineType } from "alpinejs";

declare global {
  var Alpine: AlpineType;
}

export interface SearchBoxItem {
  Label: string;
  Value: string;
}

document.addEventListener("alpine:init", () => {
  Alpine.data("SearchBox", (url: string) => ({
    async init() {
      this.$watch("Search", () => this.SearchItems());
      this.$watch("ItemSelected", function (item: SearchBoxItem) {
        if (item) console.log("item:", item);
      });
      this.Items = await (await fetch(url)).json();
    },
    Items: [],
    ItemsFiltered: [],
    ItemActive: null,
    ItemSelected: null,
    Id: "" + Date.now().toString(),
    Search: "",
    SearchIsEmpty() {
      return this.Search.length == 0;
    },
    ItemIsActive(item: SearchBoxItem) {
      return this.ItemActive && this.ItemActive.Value == item.Value;
    },
    ItemActiveNext() {
      let index = this.ItemsFiltered.indexOf(this.ItemActive);
      if (index < this.ItemsFiltered.length - 1) {
        this.ItemActive = this.ItemsFiltered[index + 1];
        this.ScrollToActiveItem();
      }
    },
    ItemActivePrevious() {
      let index = this.ItemsFiltered.indexOf(this.ItemActive);
      if (index > 0) {
        this.ItemActive = this.ItemsFiltered[index - 1];
        this.ScrollToActiveItem();
      }
    },
    ScrollToActiveItem() {
      if (this.ItemActive) {
        const activeElement = document.getElementById(
          this.ItemActive.Value + "-" + this.Id,
        );
        if (!activeElement) return;

        const newScrollPos =
          activeElement.offsetTop +
          activeElement.offsetHeight -
          this.$refs.ItemsList.offsetHeight;
        if (newScrollPos > 0) {
          this.$refs.ItemsList.scrollTop = newScrollPos;
        } else {
          this.$refs.ItemsList.scrollTop = 0;
        }
      }
    },
    SearchItems() {
      if (!this.SearchIsEmpty()) {
        const searchTerm = this.Search.replace(/\*/g, "").toLowerCase();
        console.log(this.Items);
        this.ItemsFiltered = this.Items.filter((item) => {
          console.log(item);
          return item.Label.toLowerCase().startsWith(searchTerm);
        });

        this.ScrollToActiveItem();
      } else {
        this.ItemsFiltered = this.Items.filter((item) => item.default);
      }
      this.ItemActive = this.ItemsFiltered[0];
    },
  }));
});
