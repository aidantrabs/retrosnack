<script lang="ts">
  import ProductCard from '$lib/components/ProductCard.svelte';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { mockProducts } from '$lib/mock-data';
</script>

<svelte:head>
  <title>shop — retrosnack</title>
</svelte:head>

<section class="mx-auto max-w-6xl px-4 py-12">
  <h1 class="text-2xl md:text-3xl font-semibold mb-8">the rack</h1>

  <ProductGrid>
    {#each mockProducts as product}
      <ProductCard {...product} />
    {/each}
  </ProductGrid>
</section>
