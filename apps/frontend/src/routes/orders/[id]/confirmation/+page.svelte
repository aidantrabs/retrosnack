<script lang="ts">
  import { page } from '$app/state';
  import { api } from '$lib/api';
  import type { Order } from '$lib/api';

  let order = $state<Order | null>(null);
  let loading = $state(true);
  let error = $state(false);

  $effect(() => {
    loadOrder(page.params.id);
  });

  async function loadOrder(id: string) {
    loading = true;
    error = false;
    try {
      order = await api.orders.get(id);
    } catch {
      error = true;
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>order confirmed — retrosnack clothing</title>
</svelte:head>

<section class="mx-auto max-w-2xl px-4 py-16 text-center">
  {#if loading}
    <p class="text-ink-muted">loading...</p>
  {:else if error || !order}
    <div>
      <h1 class="text-2xl font-semibold mb-2">couldn't load order</h1>
      <p class="text-ink-muted mb-6">something went wrong — please check your email for confirmation.</p>
      <a href="/shop" class="text-accent hover:text-accent-hover transition-colors">
        back to the rack &rarr;
      </a>
    </div>
  {:else}
    <h1 class="text-2xl md:text-3xl font-semibold mb-3">thanks for your order!</h1>
    <p class="text-ink-muted mb-8">
      your order is {order.status === 'paid' ? 'confirmed' : 'being processed'}. we'll be in touch soon.
    </p>

    <div class="bg-sand-light rounded-lg border border-border p-6 text-left mb-8">
      <div class="flex justify-between text-sm mb-4">
        <span class="text-ink-muted">order</span>
        <span class="font-mono text-xs">{order.id.slice(0, 8)}</span>
      </div>

      <div class="space-y-3 mb-4">
        {#each order.items as item (item.id)}
          <div class="flex justify-between text-sm">
            <span>variant {item.variant_id.slice(0, 8)}</span>
            <span class="font-semibold">${(item.price_cents / 100).toFixed(2)}</span>
          </div>
        {/each}
      </div>

      <div class="border-t border-border pt-3 flex justify-between font-semibold">
        <span>total</span>
        <span>${(order.total_cents / 100).toFixed(2)}</span>
      </div>
    </div>

    <div class="flex flex-col gap-3">
      <a
        href="/shop"
        class="bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
      >
        keep browsing
      </a>
      <a
        href="https://instagram.com/retrosnack.shop"
        target="_blank"
        rel="noopener noreferrer"
        class="border border-border px-6 py-3 rounded-full text-sm font-medium hover:bg-sand-dark transition-colors"
      >
        follow us on instagram
      </a>
    </div>
  {/if}
</section>
