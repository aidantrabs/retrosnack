<script lang="ts">
    import { page } from '$app/state';
    import { cart } from '$lib/stores/cart.svelte';

    const links = [
        { href: '/shop', label: 'shop' },
        { href: '/drops', label: 'drops' },
        { href: '/about', label: 'about' },
    ];

    let mobileOpen = $state(false);

    $effect(() => {
        page.url.pathname;
        mobileOpen = false;
    });
</script>

<nav
    class="sticky top-0 z-[60] bg-sand/90 backdrop-blur-sm"
    style="box-shadow: {mobileOpen ? 'none' : '0 1px 0 var(--color-border), 0 4px 12px rgb(0 0 0 / 0.03)'}"
>
    <div class="mx-auto max-w-6xl flex items-center justify-between px-4 py-3">
        <a href="/" class="text-xl font-semibold text-ink tracking-tight"> retrosnack clothing </a>

        <div class="hidden sm:flex items-center gap-6">
            {#each links as link}
                <a
                    href={link.href}
                    class="group relative text-sm transition-colors py-1 {page.url.pathname.startsWith(link.href)
                        ? 'text-ink font-medium'
                        : 'text-ink-muted hover:text-ink'}"
                >
                    {link.label}
                    <span
                        class="absolute bottom-0 left-0 h-[1.5px] bg-ink transition-all duration-200 {page.url.pathname.startsWith(
                            link.href,
                        )
                            ? 'w-full'
                            : 'w-0 group-hover:w-full'}"
                    ></span>
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
                    aria-hidden="true"
                >
                    <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
                    <line x1="3" y1="6" x2="21" y2="6" />
                    <path d="M16 10a4 4 0 0 1-8 0" />
                </svg>
                {#if cart.count > 0}
                    <span
                        class="absolute -top-1.5 -right-1.5 bg-accent text-white text-[10px] font-medium w-4 h-4 rounded-full flex items-center justify-center animate-badge-pop"
                    >
                        {cart.count}
                    </span>
                {/if}
            </a>
        </div>

        <div class="flex items-center gap-2 sm:hidden">
            <a
                href="/cart"
                class="relative text-ink-muted hover:text-ink transition-colors p-2"
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
                    aria-hidden="true"
                >
                    <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z" />
                    <line x1="3" y1="6" x2="21" y2="6" />
                    <path d="M16 10a4 4 0 0 1-8 0" />
                </svg>
                {#if cart.count > 0}
                    <span
                        class="absolute -top-1.5 -right-1.5 bg-accent text-white text-[10px] font-medium w-4 h-4 rounded-full flex items-center justify-center animate-badge-pop"
                    >
                        {cart.count}
                    </span>
                {/if}
            </a>

            <button
                onclick={() => (mobileOpen = !mobileOpen)}
                class="text-ink -mr-2 p-2"
                aria-label={mobileOpen ? 'close menu' : 'open menu'}
                aria-expanded={mobileOpen}
            >
                <svg
                    width="22"
                    height="22"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                >
                    {#if mobileOpen}
                        <line x1="18" y1="6" x2="6" y2="18" />
                        <line x1="6" y1="6" x2="18" y2="18" />
                    {:else}
                        <line x1="4" y1="7" x2="20" y2="7" />
                        <line x1="4" y1="12" x2="20" y2="12" />
                        <line x1="4" y1="17" x2="20" y2="17" />
                    {/if}
                </svg>
            </button>
        </div>
    </div>

</nav>

{#if mobileOpen}
    <button
        class="fixed inset-0 z-40 bg-ink/20 backdrop-blur-[2px] sm:hidden"
        aria-label="close menu"
        onclick={() => (mobileOpen = false)}
        tabindex="-1"
    ></button>

    <div
        class="fixed top-[49px] right-0 z-50 w-52 bg-sand border-l border-border sm:hidden mobile-menu-enter"
        style="height: calc(100dvh - 49px); box-shadow: -4px 0 20px rgb(0 0 0 / 0.06)"
    >
        <div class="flex flex-col px-5 pt-5 gap-1">
            {#each links as link}
                <a
                    href={link.href}
                    class="py-3 text-base font-medium transition-colors {page.url.pathname.startsWith(
                        link.href,
                    )
                        ? 'text-ink'
                        : 'text-ink-muted'}"
                >
                    {link.label}
                </a>
            {/each}
        </div>
    </div>
{/if}

<style>
    @keyframes slide-in-right {
        from {
            transform: translateX(100%);
        }
        to {
            transform: translateX(0);
        }
    }

    .mobile-menu-enter {
        animation: slide-in-right 0.25s ease-out;
    }
</style>
