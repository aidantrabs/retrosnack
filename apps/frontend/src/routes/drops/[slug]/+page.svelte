<script lang="ts">
    import { page } from '$app/state';
    import { api } from '$lib/api';
    import type { Drop, Product } from '$lib/api';
    import ProductCard from '$lib/components/ProductCard.svelte';
    import ProductGrid from '$lib/components/ProductGrid.svelte';

    let drop = $state<Drop | null>(null);
    let products = $state<Product[]>([]);
    let loading = $state(true);
    let notFound = $state(false);

    const slug = $derived(page.params.slug ?? '');

    $effect(() => {
        loadDrop(slug);
    });

    async function loadDrop(s: string) {
        loading = true;
        notFound = false;
        try {
            const [d, p] = await Promise.all([api.drops.get(s), api.drops.products(s)]);
            drop = d;
            products = p;
        } catch {
            notFound = true;
        } finally {
            loading = false;
        }
    }
</script>

<svelte:head>
    <title>{drop ? drop.name : 'not found'} - retrosnack clothing</title>
    {#if drop}
        <meta property="og:title" content="{drop.name} - retrosnack clothing" />
        <meta property="og:description" content={drop.description || `Shop the ${drop.name} drop.`} />
        <meta property="og:type" content="website" />
    {/if}
</svelte:head>

{#if loading}
    <div class="mx-auto max-w-6xl px-4 py-24 text-center">
        <p class="text-ink-muted">loading...</p>
    </div>
{:else if notFound || !drop}
    <div class="mx-auto max-w-6xl px-4 py-24 text-center">
        <h1 class="text-2xl font-semibold mb-2">drop not found</h1>
        <p class="text-ink-muted mb-6">this drop may not exist yet.</p>
        <a href="/drops" class="text-accent hover:text-accent-hover transition-colors">
            all drops &rarr;
        </a>
    </div>
{:else}
    <section class="mx-auto max-w-6xl px-4 py-12">
        <div class="mb-8">
            <a href="/drops" class="text-sm text-ink-muted hover:text-ink transition-colors">
                &larr; all drops
            </a>
            <h1 class="text-2xl md:text-3xl font-semibold mt-2">{drop.name}</h1>
            {#if drop.description}
                <p class="text-ink-muted mt-2 max-w-xl">{drop.description}</p>
            {/if}
            <div class="flex items-center gap-3 mt-4">
                {#if drop.instagram_url}
                    <a
                        href={drop.instagram_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        class="bg-accent text-sand px-4 py-2 rounded-full text-sm font-medium hover:bg-accent-hover transition-colors"
                    >
                        view on instagram
                    </a>
                {/if}
                <a
                    href="https://instagram.com/retrosnack.shop"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="border border-border px-4 py-2 rounded-full text-sm font-medium hover:bg-sand-dark transition-colors"
                >
                    DM to purchase
                </a>
            </div>
        </div>

        {#if products.length === 0}
            <p class="text-center text-ink-muted py-16">no items in this drop yet.</p>
        {:else}
            <p class="text-sm text-ink-muted mb-4">
                {products.length} item{products.length !== 1 ? 's' : ''}
            </p>
            <ProductGrid>
                {#each products as product (product.id)}
                    <ProductCard {product} />
                {/each}
            </ProductGrid>
        {/if}
    </section>
{/if}
