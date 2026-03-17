<script lang="ts">
    import { goto } from '$app/navigation';
    import { cart } from '$lib/stores/cart.svelte';
    import { api } from '$lib/api';
    import type { PaymentConfig } from '$lib/api';

    let name = $state('');
    let email = $state('');
    let address = $state('');
    let city = $state('');
    let province = $state('');
    let postal = $state('');

    let submitting = $state(false);
    let error = $state('');
    let cardReady = $state(false);

    let cardContainer: HTMLDivElement;
    let card: any;
    let payments: any;

    let config = $state<PaymentConfig | null>(null);

    $effect(() => {
        api.payments.config().then((c) => {
            config = c;
            loadSquareSDK(c);
        });
    });

    async function loadSquareSDK(cfg: PaymentConfig) {
        // load the square web payments sdk script
        if (document.getElementById('square-web-sdk')) {
            initCard(cfg);
            return;
        }
        const script = document.createElement('script');
        script.id = 'square-web-sdk';
        script.src =
            cfg.environment === 'sandbox'
                ? 'https://sandbox.web.squarecdn.com/v1/square.js'
                : 'https://web.squarecdn.com/v1/square.js';
        script.onload = () => initCard(cfg);
        document.head.appendChild(script);
    }

    async function initCard(cfg: PaymentConfig) {
        if (!window.Square) {
            return;
        }

        payments = window.Square.payments(cfg.application_id, cfg.location_id);
        card = await payments.card();
        await card.attach(cardContainer);
        cardReady = true;
    }

    async function handleSubmit(e: SubmitEvent) {
        e.preventDefault();
        if (submitting || !cardReady || cart.items.length === 0) {
            return;
        }

        submitting = true;
        error = '';

        try {
            // tokenize the card
            const tokenResult = await card.tokenize();
            if (tokenResult.status !== 'OK') {
                error = 'please check your card details and try again.';
                submitting = false;
                return;
            }

            // create order
            const order = await api.orders.create(
                cart.items.map((item) => ({
                    variant_id: item.variantId,
                    quantity: 1,
                    price_cents: item.priceCents,
                })),
            );

            // process payment with the card token
            const result = await api.payments.process(order.id, tokenResult.token);

            if (result.status === 'paid') {
                cart.clear();
                goto(`/orders/${result.order_id}/confirmation`);
            } else {
                error = 'payment is being processed — you will receive a confirmation soon.';
                cart.clear();
            }
        } catch {
            error = 'something went wrong — please try again.';
            submitting = false;
        }
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
            <a href="/shop" class="btn-primary inline-block px-6 py-3">
                browse the rack
            </a>
        </div>
    {:else}
        <div class="grid md:grid-cols-5 gap-10">
            <form onsubmit={handleSubmit} class="md:col-span-3 space-y-6">
                <div class="space-y-4">
                    <h2 class="text-lg font-medium">shipping</h2>

                    <div>
                        <label for="name" class="block text-sm text-ink-muted mb-1">full name</label
                        >
                        <input
                            id="name"
                            type="text"
                            required
                            autocomplete="name"
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
                            autocomplete="email"
                            bind:value={email}
                            class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
                        />
                    </div>

                    <div>
                        <label for="address" class="block text-sm text-ink-muted mb-1"
                            >address</label
                        >
                        <input
                            id="address"
                            type="text"
                            required
                            autocomplete="street-address"
                            bind:value={address}
                            class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
                        />
                    </div>

                    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
                        <div>
                            <label for="city" class="block text-sm text-ink-muted mb-1">city</label>
                            <input
                                id="city"
                                type="text"
                                required
                                autocomplete="address-level2"
                                bind:value={city}
                                class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
                            />
                        </div>
                        <div>
                            <label for="province" class="block text-sm text-ink-muted mb-1"
                                >province</label
                            >
                            <input
                                id="province"
                                type="text"
                                required
                                autocomplete="address-level1"
                                bind:value={province}
                                class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
                            />
                        </div>
                        <div>
                            <label for="postal" class="block text-sm text-ink-muted mb-1"
                                >postal code</label
                            >
                            <input
                                id="postal"
                                type="text"
                                required
                                autocomplete="postal-code"
                                bind:value={postal}
                                class="w-full bg-sand-light border border-border rounded-lg px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ink/20"
                            />
                        </div>
                    </div>
                </div>

                <div class="space-y-4">
                    <h2 class="text-lg font-medium">payment</h2>
                    <div
                        bind:this={cardContainer}
                        class="min-h-[44px] rounded-lg border border-border bg-sand-light p-1"
                    ></div>
                </div>

                {#if error}
                    <p class="text-red-600 text-sm">{error}</p>
                {/if}

                <button
                    type="submit"
                    disabled={submitting || !cardReady}
                    class="w-full px-6 py-3 {submitting || !cardReady
                        ? 'bg-ink/30 text-sand rounded-full text-sm font-medium cursor-not-allowed'
                        : 'btn-primary'}"
                >
                    {#if submitting}
                        processing...
                    {:else if !cardReady}
                        loading payment...
                    {:else}
                        pay ${(cart.totalCents / 100).toFixed(2)}
                    {/if}
                </button>
            </form>

            <aside class="md:col-span-2">
                <div class="bg-sand-light rounded-lg p-6 border border-border sticky top-20">
                    <h2 class="text-lg font-medium mb-4">order summary</h2>

                    <div class="space-y-4 mb-6">
                        {#each cart.items as item (item.variantId)}
                            <div class="flex gap-3">
                                {#if item.image}
                                    <img
                                        src={item.image}
                                        alt={item.title}
                                        loading="lazy"
                                        width="56"
                                        height="72"
                                        class="w-14 h-18 object-cover rounded bg-sand-dark"
                                    />
                                {:else}
                                    <div class="w-14 h-18 rounded bg-sand-dark"></div>
                                {/if}
                                <div class="flex-1 min-w-0">
                                    <p class="text-sm font-medium truncate">{item.title}</p>
                                    <p class="text-xs text-ink-muted">
                                        {item.size}{item.color ? ` · ${item.color}` : ''}
                                    </p>
                                </div>
                                <p class="text-sm font-semibold shrink-0">
                                    ${(item.priceCents / 100).toFixed(2)}
                                </p>
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
                            <span class="text-ink-muted">free</span>
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
