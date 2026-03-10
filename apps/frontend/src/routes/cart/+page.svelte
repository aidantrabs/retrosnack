<script lang="ts">
  import { cart } from '$lib/stores/cart.svelte';
</script>

<svelte:head>
  <title>bag — retrosnack clothing</title>
</svelte:head>

<section class="mx-auto max-w-3xl px-4 py-12">
  <h1 class="text-2xl md:text-3xl font-semibold mb-8">your bag</h1>

  {#if cart.items.length === 0}
    <div class="text-center py-16">
      <p class="text-ink-muted mb-6">your bag is empty.</p>
      <a
        href="/shop"
        class="inline-block bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
      >
        browse the rack
      </a>
    </div>
  {:else}
    <div class="space-y-6">
      {#each cart.items as item (item.variantId)}
        <div class="flex gap-4 pb-6 border-b border-border">
          <a href="/shop/{item.productId}" class="shrink-0">
            <img
              src={item.image}
              alt={item.title}
              class="w-24 h-32 object-cover rounded-lg bg-sand-dark"
            />
          </a>

          <div class="flex-1 flex flex-col justify-between">
            <div>
              <a href="/shop/{item.productId}" class="font-medium hover:text-accent transition-colors">
                {item.title}
              </a>
              <p class="text-sm text-ink-muted mt-0.5">{item.size}{item.color ? ` · ${item.color}` : ''}</p>
            </div>

            <p class="font-semibold">${(item.priceCents / 100).toFixed(2)}</p>
          </div>

          <button
            onclick={() => cart.remove(item.variantId)}
            class="self-start text-ink-muted hover:text-accent transition-colors"
            aria-label="remove item"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
      {/each}
    </div>

    <div class="mt-8 pt-6 border-t border-border">
      <div class="flex items-center justify-between mb-6">
        <span class="text-ink-muted">total</span>
        <span class="text-2xl font-semibold">${(cart.totalCents / 100).toFixed(2)}</span>
      </div>

      <a
        href="/checkout"
        class="block w-full bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium text-center hover:bg-ink/85 transition-colors"
      >
        proceed to checkout
      </a>
    </div>
  {/if}
</section>
