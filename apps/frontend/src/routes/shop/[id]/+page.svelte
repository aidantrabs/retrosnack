<script lang="ts">
  import { page } from '$app/state';
  import { mockProducts } from '$lib/mock-data';
  import { cart } from '$lib/stores/cart.svelte';

  const product = $derived(mockProducts.find((p) => p.id === page.params.id));
  let added = $state(false);

  function addToCart() {
    if (!product) return;
    cart.add({
      id: product.id,
      title: product.title,
      price: product.price,
      size: product.size,
      image: product.image,
    });
    added = true;
    setTimeout(() => (added = false), 1500);
  }
</script>

<svelte:head>
  <title>{product ? product.title : 'not found'} — retrosnack clothing</title>
</svelte:head>

{#if product}
  <article class="mx-auto max-w-4xl px-4 py-12">
    <div class="grid md:grid-cols-2 gap-8 md:gap-12">
      <div class="aspect-[3/4] overflow-hidden rounded-lg bg-sand-dark">
        <img
          src={product.image}
          alt={product.title}
          class="h-full w-full object-cover"
        />
      </div>

      <div class="flex flex-col justify-center gap-6">
        <div>
          <h1 class="text-2xl md:text-3xl font-semibold">{product.title}</h1>
          <p class="mt-2 text-ink-muted">size {product.size}</p>
        </div>

        <p class="text-3xl font-semibold">
          ${(product.price / 100).toFixed(2)}
        </p>

        <div class="flex flex-col gap-3">
          <button
            onclick={addToCart}
            class="bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
          >
            {added ? 'added ✓' : 'add to bag'}
          </button>

          <a
            href="https://instagram.com/retrosnack.shop"
            target="_blank"
            rel="noopener noreferrer"
            class="border border-border px-6 py-3 rounded-full text-sm font-medium text-center hover:bg-sand-dark transition-colors"
          >
            message on instagram
          </a>
        </div>
      </div>
    </div>
  </article>
{:else}
  <div class="mx-auto max-w-4xl px-4 py-24 text-center">
    <h1 class="text-2xl font-semibold mb-2">item not found</h1>
    <p class="text-ink-muted mb-6">this piece may have already found a new home.</p>
    <a href="/shop" class="text-accent hover:text-accent-hover transition-colors">
      back to the rack &rarr;
    </a>
  </div>
{/if}
