/// <reference types="alpinejs" />
import type {
  Alpine as AlpineType,
  AlpineComponent,
  Magics,
  ElementWithXAttributes,
  InferInterceptors,
  Stores,
  XDataContext,
} from "alpinejs";

declare global {
  var Alpine: AlpineType;
}

export interface SearchBoxSelectEvent extends Event {
  detail: {
    item: SearchBoxItem;
  };
}

export interface SearchBoxItem {
  Label: string;
  Value: string;
}

interface SearchBoxType {
  Items: SearchBoxItem[];
  ItemsFiltered: SearchBoxItem[];
  ItemActive: SearchBoxItem | null;
  ItemSelected: SearchBoxItem | null;
  Id: string;
  Search: string;
  SearchIsEmpty(): boolean;
  ItemIsActive(item: SearchBoxItem): boolean;
  ItemActiveNext(): void;
  ItemActivePrevious(): void;
  ScrollToActiveItem(): void;
  SearchItems(): void;
  SelectItem(item: SearchBoxItem): void;
  AddItem(): void;
}

function getSearchBoxClass() {
  return class SearchBox implements Magics<AlpineComponent<SearchBoxType>> {
    // MAGIC METHODS
    // Alpine magic properties and methods
    $data: InferInterceptors<AlpineComponent<SearchBoxType>>;
    $dispatch: (event: string, detail?: any) => void;
    $el: HTMLElement;
    $id: (name: string, key?: number | string | null) => string;
    $nextTick: (callback?: () => void) => Promise<void>;
    $refs: Record<string, HTMLElement>;
    $root: ElementWithXAttributes<HTMLElement>;
    $store: Stores;
    $watch: <
      K extends string,
      V extends K extends keyof SearchBoxType | keyof XDataContext
        ? AlpineComponent<SearchBoxType>[K]
        : any,
    >(
      property: K,
      callback: (newValue: V, oldValue: V) => void,
    ) => void;
    destroy: () => void;
    // CLASS PROPERTIES
    name: string;
    Id: string;
    Items: SearchBoxItem[];
    ItemsFiltered: SearchBoxItem[];
    ItemActive: SearchBoxItem | null;
    ItemSelected: SearchBoxItem | null;
    Search: string;
    fetchItemsUrl: string;
    // Constructor, initialize properties
    constructor(url: string, name: string) {
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
      this.$watch("ItemSelected", (item: SearchBoxItem) =>
        this.SelectItem(item),
      );
      this.Items = await (await fetch(this.fetchItemsUrl)).json();
    }
    // CLASS METHODS
    SearchIsEmpty() {
      return this.Search.length == 0;
    }
    ItemIsActive(item: SearchBoxItem) {
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
    }
    SearchItems() {
      if (this.SearchIsEmpty()) {
        this.ClearState();
        return;
      }
      if (!this.SearchIsEmpty()) {
        const searchTerm: string = this.Search.replace(/\*/g, "").toLowerCase();
        this.ItemsFiltered = this.Items.filter((item: SearchBoxItem) =>
          item.Label.toLowerCase().includes(searchTerm),
        );
        this.ScrollToActiveItem();
      }
    }
    SelectItem(item: SearchBoxItem) {
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
  const createSearchBox = (url: string, name: string) => {
    const SearchBox = getSearchBoxClass();
    return new SearchBox(url, name);
  };

  Alpine.data("SearchBox", createSearchBox);
});
