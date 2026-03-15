<script lang="ts">
    import HeroSection from '$lib/components/HeroSection.svelte';
    import ProductCard from '$lib/components/ProductCard.svelte';
    import ProductCardSkeleton from '$lib/components/ProductCardSkeleton.svelte';
    import ProductGrid from '$lib/components/ProductGrid.svelte';
    import InstagramCTA from '$lib/components/InstagramCTA.svelte';
    import { api } from '$lib/api';
    import type { Product } from '$lib/api';

    let products = $state<Product[]>([]);
    let loading = $state(true);

    $effect(() => {
        api.products
            .list(4, 0)
            .then((p) => (products = p))
            .catch(() => {})
            .finally(() => (loading = false));
    });
</script>

<svelte:head>
    <title>retrosnack clothing — thrift finds, loved again</title>
    <meta property="og:title" content="retrosnack clothing — thrift finds, loved again" />
    <meta
        property="og:description"
        content="Curated secondhand women's clothing, accessories, and shoes."
    />
    <meta property="og:type" content="website" />
</svelte:head>

<HeroSection />

<section class="mx-auto max-w-6xl px-4 pb-16">
    <div class="flex items-center justify-between mb-6">
        <h2 class="text-xl md:text-2xl font-semibold">latest drops</h2>
        <a href="/shop" class="text-sm text-accent hover:text-accent-hover transition-colors">
            view all &rarr;
        </a>
    </div>

    {#if loading}
        <ProductGrid>
            {#each Array(4) as _}
                <ProductCardSkeleton />
            {/each}
        </ProductGrid>
    {:else if products.length > 0}
        <ProductGrid>
            {#each products as product (product.id)}
                <ProductCard {product} />
            {/each}
        </ProductGrid>
    {/if}
</section>

<InstagramCTA />
