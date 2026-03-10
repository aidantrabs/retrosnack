<script lang="ts">
  import { cart } from '$lib/stores/cart.svelte';

  let name = $state('');
  let email = $state('');
  let address = $state('');
  let city = $state('');
  let province = $state('');
  let postal = $state('');

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    // placeholder — will connect to square checkout via backend
  }
</script>

<svelte:head>
  <title>checkout — retrosnack clothing</title>
</svelte:head>

<section class="mx-auto max-w-4xl px-4 py-12">
  <h1 class="text-2xl md:text-3xl font-semibold mb-8">checkout</h1>

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
    <div class="grid md:grid-cols-5 gap-10">
      <form onsubmit={handleSubmit} class="md:col-span-3 space-y-5">
        <h2 class="text-lg font-medium mb-4">shipping details</h2>

        <div>
          <label for="name" class="block text-sm text-ink-muted mb-1">full name</label>
          <input
            id="name"
            type="text"
            required
            bind:value={name}
            class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
          />
        </div>

        <div>
          <label for="email" class="block text-sm text-ink-muted mb-1">email</label>
          <input
            id="email"
            type="email"
            required
            bind:value={email}
            class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
          />
        </div>

        <div>
          <label for="address" class="block text-sm text-ink-muted mb-1">address</label>
          <input
            id="address"
            type="text"
            required
            bind:value={address}
            class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
          />
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div>
            <label for="city" class="block text-sm text-ink-muted mb-1">city</label>
            <input
              id="city"
              type="text"
              required
              bind:value={city}
              class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
            />
          </div>
          <div>
            <label for="province" class="block text-sm text-ink-muted mb-1">province</label>
            <input
              id="province"
              type="text"
              required
              bind:value={province}
              class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
            />
          </div>
          <div>
            <label for="postal" class="block text-sm text-ink-muted mb-1">postal code</label>
            <input
              id="postal"
              type="text"
              required
              bind:value={postal}
              class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
            />
          </div>
        </div>

        <button
          type="submit"
          class="w-full bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors mt-4"
        >
          proceed to payment
        </button>
      </form>

      <aside class="md:col-span-2">
        <div class="bg-sand-light rounded-lg p-6 border border-border">
          <h2 class="text-lg font-medium mb-4">order summary</h2>

          <div class="space-y-4 mb-6">
            {#each cart.items as item (item.variantId)}
              <div class="flex gap-3">
                <img
                  src={item.image}
                  alt={item.title}
                  class="w-14 h-18 object-cover rounded bg-sand-dark"
                />
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium truncate">{item.title}</p>
                  <p class="text-xs text-ink-muted">{item.size}{item.color ? ` · ${item.color}` : ''}</p>
                </div>
                <p class="text-sm font-semibold shrink-0">${(item.priceCents / 100).toFixed(2)}</p>
              </div>
            {/each}
          </div>

          <div class="border-t border-border pt-4 space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-ink-muted">subtotal</span>
              <span>${(cart.totalCents / 100).toFixed(2)}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-ink-muted">shipping</span>
              <span class="text-ink-muted">calculated at payment</span>
            </div>
            <div class="flex justify-between font-semibold pt-2 border-t border-border">
              <span>total</span>
              <span>${(cart.totalCents / 100).toFixed(2)}</span>
            </div>
          </div>
        </div>
      </aside>
    </div>
  {/if}
</section>
