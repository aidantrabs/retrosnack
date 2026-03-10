<script lang="ts">
  import HeroSection from '$lib/components/HeroSection.svelte';
  import ProductCard from '$lib/components/ProductCard.svelte';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { mockProducts } from '$lib/mock-data';

  const featured = mockProducts.slice(0, 4);
</script>

<svelte:head>
  <title>retrosnack — thrift finds, loved again</title>
</svelte:head>

<HeroSection />

<section class="mx-auto max-w-6xl px-4 pb-16">
  <div class="flex items-center justify-between mb-6">
    <h2 class="text-xl md:text-2xl font-semibold">latest drops</h2>
    <a href="/shop" class="text-sm text-accent hover:text-accent-hover transition-colors">
      view all &rarr;
    </a>
  </div>

  <ProductGrid>
    {#each featured as product}
      <ProductCard {...product} />
    {/each}
  </ProductGrid>
</section>
