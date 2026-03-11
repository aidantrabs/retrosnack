<script lang="ts">
  import { page } from '$app/state';
  import { api } from '$lib/api';
  import type { Product, Variant } from '$lib/api';
  import ProductCard from '$lib/components/ProductCard.svelte';
  import ProductGrid from '$lib/components/ProductGrid.svelte';
  import { cart } from '$lib/stores/cart.svelte';
  import { toast } from '$lib/stores/toast.svelte';

  let product = $state<Product | null>(null);
  let variants = $state<Variant[]>([]);
  let selectedVariant = $state<Variant | null>(null);
  let loading = $state(true);
  let notFound = $state(false);
  let justAdded = $state(false);
  let moreItems = $state<Product[]>([]);

  const inCart = $derived(
    selectedVariant ? cart.items.some((i) => i.variantId === selectedVariant!.id) : cart.has(page.params.id)
  );

  const image = $derived(product?.images[0]?.url ?? '');

  $effect(() => {
    loadProduct(page.params.id);
  });

  async function loadProduct(id: string) {
    loading = true;
    notFound = false;
    justAdded = false;
    try {
      const [p, v, all] = await Promise.all([
        api.products.get(id),
        api.products.variants(id),
        api.products.list(4, 0),
      ]);
      product = p;
      variants = v;
      moreItems = all.filter((item) => item.id !== id).slice(0, 3);
      if (v.length === 1) selectedVariant = v[0];
    } catch {
      notFound = true;
    } finally {
      loading = false;
    }
  }

  function addToCart() {
    if (!product || !selectedVariant || inCart) return;
    cart.add({
      productId: product.id,
      variantId: selectedVariant.id,
      title: product.title,
      priceCents: product.price_cents,
      size: selectedVariant.size,
      color: selectedVariant.color,
      image,
    });
    justAdded = true;
    toast.show('added to your bag');
  }
</script>

<svelte:head>
  <title>{product ? product.title : 'not found'} — retrosnack clothing</title>
</svelte:head>

{#if loading}
  <div class="mx-auto max-w-4xl px-4 py-24 text-center">
    <p class="text-ink-muted">loading...</p>
  </div>
{:else if notFound || !product}
  <div class="mx-auto max-w-4xl px-4 py-24 text-center">
    <h1 class="text-2xl font-semibold mb-2">item not found</h1>
    <p class="text-ink-muted mb-6">this piece may have already found a new home.</p>
    <a href="/shop" class="text-accent hover:text-accent-hover transition-colors">
      back to the rack &rarr;
    </a>
  </div>
{:else}
  <article class="mx-auto max-w-4xl px-4 py-12">
    <div class="grid md:grid-cols-2 gap-8 md:gap-12">
      <div class="aspect-[3/4] overflow-hidden rounded-lg bg-sand-dark">
        {#if image}
          <img
            src={image}
            alt={product.title}
            class="h-full w-full object-cover"
          />
        {/if}
      </div>

      <div class="flex flex-col justify-center gap-6">
        <div>
          <p class="text-sm text-ink-muted mb-1">{product.brand}</p>
          <h1 class="text-2xl md:text-3xl font-semibold">{product.title}</h1>
        </div>

        <div class="flex items-center gap-3 text-sm text-ink-muted">
          <span class="bg-tag px-3 py-1 rounded-full">{product.condition}</span>
        </div>

        {#if product.description}
          <p class="text-sm text-ink-muted leading-relaxed">{product.description}</p>
        {/if}

        <p class="text-3xl font-semibold">
          ${(product.price_cents / 100).toFixed(2)}
        </p>

        {#if variants.length > 1}
          <div class="flex flex-wrap gap-2">
            {#each variants as variant (variant.id)}
              <button
                onclick={() => (selectedVariant = variant)}
                class="px-4 py-2 rounded-full text-sm border transition-colors {selectedVariant?.id === variant.id
                  ? 'border-ink bg-ink text-sand'
                  : 'border-border hover:border-ink'}"
              >
                {variant.size}{variant.color ? ` / ${variant.color}` : ''}
              </button>
            {/each}
          </div>
        {:else if variants.length === 1}
          <span class="text-sm text-ink-muted">size {variants[0].size}{variants[0].color ? ` · ${variants[0].color}` : ''}</span>
        {/if}

        <div class="flex flex-col gap-3">
          {#if variants.length === 0}
            <button
              disabled
              class="bg-ink/40 text-sand px-6 py-3 rounded-full text-sm font-medium cursor-not-allowed"
            >
              sold
            </button>
          {:else if justAdded}
            <div class="bg-sand-light border border-border rounded-lg p-4 text-center space-y-3">
              <p class="text-sm font-medium">added to your bag</p>
              <div class="flex gap-3">
                <a
                  href="/shop"
                  class="flex-1 border border-border px-4 py-2.5 rounded-full text-sm font-medium text-center hover:bg-sand-dark transition-colors"
                >
                  keep shopping
                </a>
                <a
                  href="/cart"
                  class="flex-1 bg-ink text-sand px-4 py-2.5 rounded-full text-sm font-medium text-center hover:bg-ink/85 transition-colors"
                >
                  view bag ({cart.count})
                </a>
              </div>
            </div>
          {:else if inCart}
            <div class="flex gap-3">
              <a
                href="/shop"
                class="flex-1 border border-border px-4 py-2.5 rounded-full text-sm font-medium text-center hover:bg-sand-dark transition-colors"
              >
                keep shopping
              </a>
              <a
                href="/cart"
                class="flex-1 bg-ink text-sand px-4 py-2.5 rounded-full text-sm font-medium text-center hover:bg-ink/85 transition-colors"
              >
                view bag ({cart.count})
              </a>
            </div>
          {:else if !selectedVariant}
            <button
              disabled
              class="bg-ink/40 text-sand px-6 py-3 rounded-full text-sm font-medium cursor-not-allowed"
            >
              select a size
            </button>
          {:else}
            <button
              onclick={addToCart}
              class="bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
            >
              add to bag
            </button>
          {/if}

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

  {#if moreItems.length > 0}
    <section class="mx-auto max-w-4xl px-4 pb-16">
      <h2 class="text-lg font-semibold mb-4">you might also like</h2>
      <ProductGrid>
        {#each moreItems as item (item.id)}
          <ProductCard product={item} />
        {/each}
      </ProductGrid>
    </section>
  {/if}
{/if}
