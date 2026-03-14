<script lang="ts">
    import { page } from '$app/state';
    import { cart } from '$lib/stores/cart.svelte';

    const links = [
        { href: '/shop', label: 'shop' },
        { href: '/drops', label: 'drops' },
        { href: '/about', label: 'about' },
    ];
</script>

<nav class="sticky top-0 z-50 bg-sand/90 backdrop-blur-sm border-b border-border">
    <div class="mx-auto max-w-6xl flex items-center justify-between px-4 py-3">
        <a href="/" class="text-xl font-semibold text-ink tracking-tight"> retrosnack clothing </a>

        <div class="flex items-center gap-6">
            {#each links as link}
                <a
                    href={link.href}
                    class="text-sm transition-colors {page.url.pathname.startsWith(link.href)
                        ? 'text-ink font-medium'
                        : 'text-ink-muted hover:text-ink'}"
                >
                    {link.label}
                </a>
            {/each}

            <a
                href="/cart"
                class="relative text-ink-muted hover:text-ink transition-colors"
                aria-label="cart"
            >
                <svg
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
                    <line x1="3" y1="6" x2="21" y2="6" />
                    <path d="M16 10a4 4 0 0 1-8 0" />
                </svg>
                {#if cart.count > 0}
                    <span
                        class="absolute -top-1.5 -right-1.5 bg-accent text-white text-[10px] font-medium w-4 h-4 rounded-full flex items-center justify-center"
                    >
                        {cart.count}
                    </span>
                {/if}
            </a>
        </div>
    </div>
</nav>
