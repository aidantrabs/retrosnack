<script lang="ts">
    import HeroSection from '$lib/components/HeroSection.svelte';
    import ProductCard from '$lib/components/ProductCard.svelte';
    import ProductCardSkeleton from '$lib/components/ProductCardSkeleton.svelte';
    import ProductGrid from '$lib/components/ProductGrid.svelte';
    import InstagramCTA from '$lib/components/InstagramCTA.svelte';
    import FadeIn from '$lib/components/FadeIn.svelte';
    import Marquee from '$lib/components/Marquee.svelte';
    import { api } from '$lib/api';
    import type { Product, Drop } from '$lib/api';

    let products = $state<Product[]>([]);
    let drops = $state<Drop[]>([]);
    let loading = $state(true);

    $effect(() => {
        Promise.all([api.products.list(4, 0), api.drops.list()])
            .then(([p, d]) => {
                products = p;
                drops = d;
            })
            .catch(() => {})
            .finally(() => (loading = false));
    });

    const latestDrop = $derived(drops[0]);
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

<Marquee items={['DM to purchase', 'new drops weekly', 'curated with care', 'secondhand style', 'shipped from canada']} />

<FadeIn>
    <section class="mx-auto max-w-6xl px-4 py-16">
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
</FadeIn>

{#if latestDrop}
    <FadeIn>
        <section class="mx-auto max-w-6xl px-4 pb-16">
            <a
                href="/drops/{latestDrop.slug}"
                class="group block border border-border rounded-2xl p-8 md:p-12 hover:border-ink transition-colors hover-lift text-center"
            >
                <p class="text-xs uppercase tracking-widest text-ink-muted mb-3">latest drop</p>
                <h2 class="text-2xl md:text-4xl font-semibold group-hover:text-accent transition-colors">
                    {latestDrop.name}
                </h2>
                {#if latestDrop.description}
                    <p class="text-ink-muted mt-3 max-w-md mx-auto">{latestDrop.description}</p>
                {/if}
                <span
                    class="inline-block mt-5 bg-ink text-sand px-5 py-2.5 rounded-full text-sm font-medium group-hover:bg-ink/85 transition-colors press"
                >
                    shop the drop &rarr;
                </span>
            </a>
        </section>
    </FadeIn>
{/if}

<FadeIn>
    <InstagramCTA />
</FadeIn>
