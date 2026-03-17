<script lang="ts">
    import { page } from '$app/state';

    const orderId = page.params.id ?? '';
</script>

<svelte:head>
    <title>order confirmed — retrosnack clothing</title>
</svelte:head>

<section class="mx-auto max-w-2xl px-4 py-16 text-center">
    <h1 class="text-2xl md:text-3xl font-semibold mb-3">thanks for your order!</h1>
    <p class="text-ink-muted mb-2">your payment has been received and your order is confirmed.</p>
    <p class="text-ink-muted mb-8">
        order <span class="font-mono text-xs">{orderId.slice(0, 8)}</span>
    </p>

    <div class="flex flex-col gap-3">
        <a href="/shop" class="btn-primary px-6 py-3 text-center">
            keep browsing
        </a>
        <a
            href="https://instagram.com/retrosnack.shop"
            target="_blank"
            rel="noopener noreferrer"
            class="btn-outline px-6 py-3 text-center"
        >
            follow us on instagram
        </a>
    </div>
</section>
