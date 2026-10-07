# Catalog card 1c — implementation notes

Replace the `columns md:columns-4 …` masonry block in `views/catalog/catalog.templ` with the grid below. Uses only existing tokens from `assets/css/main.css`.

```templ
<div class="my-16 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 lg:gap-7">
	for _, toy := range toys {
		<a
			href={ templ.URL(getCheckoutToyLink(toy.ToyID, code)) }
			class="group flex flex-col overflow-hidden rounded-image border border-border bg-surface shadow-[0_2px_10px_rgba(74,63,77,0.06)] transition hover:-translate-y-1 hover:shadow-[0_12px_28px_rgba(74,63,77,0.12)]"
		>
			<div class="relative aspect-4/3 overflow-hidden border-b border-border bg-surface p-5 sm:aspect-square sm:p-6">
				<img src={ toy.Image1 } alt={ toy.ToyDescription } loading="lazy" class="h-full w-full object-contain"/>
				<span class="absolute left-3.5 top-3.5 rounded-pill border border-border bg-surface px-3 py-1 font-heading text-xs font-bold text-ink">
					{ fmt.Sprintf("%d–%d años", toy.AgeMin, toy.AgeMax) }
				</span>
			</div>
			<div class="flex flex-1 flex-col gap-3 px-5 pb-5 pt-4">
				<h3 class="min-h-[46px] font-heading text-[17px] font-bold leading-[1.35] text-ink line-clamp-2">{ toy.ToyName }</h3>
				@CategoryChips(toy.Category, 2)
				<span class="mt-auto flex items-center justify-center rounded-button border-2 border-primary px-4 py-2.5 font-heading text-sm font-bold text-accent transition group-hover:bg-primary group-hover:text-white">
					Elegir regalo
				</span>
			</div>
		</a>
	}
</div>
```

```templ
templ CategoryChips(category string, max int) {
	{{ cats := strings.Split(category, ",") }}
	<div class="flex flex-wrap gap-1.5">
		for i, c := range cats {
			if i < max {
				<span class="rounded-pill bg-chip px-2.5 py-1 font-sans text-xs text-ink">{ strings.TrimSpace(c) }</span>
			}
		}
		if len(cats) > max {
			<span class="px-1 py-1 font-sans text-xs text-muted">{ fmt.Sprintf("+%d", len(cats)-max) }</span>
		}
	</div>
}
```

Add `"strings"` to the imports. Images assume white backgrounds (3rd-party hosted) — the frame is white, no tint.
