<script lang="ts">
    import ProductCard from '$lib/components/ProductCard.svelte';
    import ProductCardSkeleton from '$lib/components/ProductCardSkeleton.svelte';
    import ProductGrid from '$lib/components/ProductGrid.svelte';
    import { api } from '$lib/api';
    import type { Product, Category } from '$lib/api';

    let products = $state<Product[]>([]);
    let categories = $state<Category[]>([]);
    let loading = $state(true);
    let error = $state('');

    let search = $state('');
    let activeCategory = $state('all');
    let sort = $state('newest');

    const filtered = $derived.by(() => {
        let result = products;

        // search
        if (search.trim()) {
            const q = search.toLowerCase();
            result = result.filter(
                (p) => p.title.toLowerCase().includes(q) || p.brand.toLowerCase().includes(q),
            );
        }

        // category
        if (activeCategory !== 'all') {
            result = result.filter((p) => p.category_id === activeCategory);
        }

        // sort
        if (sort === 'price-asc') {
            result = [...result].sort((a, b) => a.price_cents - b.price_cents);
        } else if (sort === 'price-desc') {
            result = [...result].sort((a, b) => b.price_cents - a.price_cents);
        }
        // "newest" is the default api order

        return result;
    });

    async function loadProducts() {
        try {
            const [p, c] = await Promise.all([api.products.list(100), api.categories.list()]);
            products = p;
            categories = c;
        } catch {
            error = 'failed to load products';
        } finally {
            loading = false;
        }
    }

    $effect(() => {
        loadProducts();
    });
</script>

<svelte:head>
    <title>shop — retrosnack clothing</title>
    <meta property="og:title" content="shop — retrosnack clothing" />
    <meta
        property="og:description"
        content="Browse curated secondhand clothing, accessories, and shoes."
    />
    <meta property="og:type" content="website" />
</svelte:head>

<section class="mx-auto max-w-6xl px-4 py-12">
    <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl md:text-3xl font-semibold">the rack</h1>
        {#if !loading && !error}
            <span class="text-sm text-ink-muted"
                >{filtered.length} item{filtered.length !== 1 ? 's' : ''}</span
            >
        {/if}
    </div>

    {#if loading}
        <ProductGrid>
            {#each Array(8) as _}
                <ProductCardSkeleton />
            {/each}
        </ProductGrid>
    {:else if error}
        <p class="text-center text-ink-muted py-16">{error}</p>
    {:else if products.length === 0}
        <p class="text-center text-ink-muted py-16">
            no items on the rack right now — check back soon.
        </p>
    {:else}
        <div class="space-y-5 mb-8">
            <div class="relative">
                <svg
                    class="absolute left-3.5 top-1/2 -translate-y-1/2 text-ink-muted"
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    aria-hidden="true"
                >
                    <circle cx="11" cy="11" r="8" />
                    <line x1="21" y1="21" x2="16.65" y2="16.65" />
                </svg>
                <input
                    type="text"
                    placeholder="search by name or brand..."
                    aria-label="search products"
                    bind:value={search}
                    class="w-full bg-sand-light border border-border rounded-full pl-10 pr-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20 focus:border-ink/30"
                    style="box-shadow: var(--shadow-soft); transition: box-shadow 0.15s"
                />
            </div>

            <div class="flex flex-wrap items-center gap-3">
                <div class="flex flex-wrap gap-2">
                    <button
                        onclick={() => (activeCategory = 'all')}
                        class="px-3.5 py-1.5 rounded-full text-sm border transition-colors {activeCategory ===
                        'all'
                            ? 'border-ink bg-ink text-sand'
                            : 'border-border hover:border-ink'}"
                    >
                        all
                    </button>
                    {#each categories as cat (cat.id)}
                        <button
                            onclick={() => (activeCategory = cat.id)}
                            class="px-3.5 py-1.5 rounded-full text-sm border transition-colors {activeCategory ===
                            cat.id
                                ? 'border-ink bg-ink text-sand'
                                : 'border-border hover:border-ink'}"
                        >
                            {cat.name}
                        </button>
                    {/each}
                </div>

                <select
                    bind:value={sort}
                    aria-label="sort products"
                    class="ml-auto bg-sand-light border border-border rounded-full px-4 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20 cursor-pointer"
                    style="box-shadow: var(--shadow-soft)"
                >
                    <option value="newest">newest</option>
                    <option value="price-asc">price: low to high</option>
                    <option value="price-desc">price: high to low</option>
                </select>
            </div>
        </div>

        {#if filtered.length === 0}
            <p class="text-center text-ink-muted py-16">
                no items match your search — try something else.
            </p>
        {:else}
            <ProductGrid>
                {#each filtered as product (product.id)}
                    <ProductCard {product} />
                {/each}
            </ProductGrid>
        {/if}
    {/if}
</section>
