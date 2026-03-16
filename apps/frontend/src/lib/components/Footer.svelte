<script lang="ts">
    import FlowerDoodle from './FlowerDoodle.svelte';
</script>

<footer class="border-t border-border mt-auto">
    <div class="mx-auto max-w-6xl px-4 py-12">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 mb-10">
            <div class="col-span-2 md:col-span-1">
                <p class="font-semibold text-lg mb-3">retrosnack</p>
                <p class="text-sm text-ink-muted leading-relaxed">
                    curated secondhand finds, loved again.
                </p>
                <FlowerDoodle size={24} class="text-ink/20 mt-4" />
            </div>

            <div>
                <p class="text-xs uppercase tracking-widest text-ink-muted mb-4">shop</p>
                <nav class="flex flex-col gap-2.5">
                    <a href="/shop" class="text-sm text-ink-muted hover:text-ink transition-colors">the rack</a>
                    <a href="/drops" class="text-sm text-ink-muted hover:text-ink transition-colors">drops</a>
                    <a href="/cart" class="text-sm text-ink-muted hover:text-ink transition-colors">your bag</a>
                </nav>
            </div>

            <div>
                <p class="text-xs uppercase tracking-widest text-ink-muted mb-4">info</p>
                <nav class="flex flex-col gap-2.5">
                    <a href="/about" class="text-sm text-ink-muted hover:text-ink transition-colors">about</a>
                    <a
                        href="https://instagram.com/retrosnack.shop"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="text-sm text-ink-muted hover:text-ink transition-colors"
                    >
                        instagram
                    </a>
                </nav>
            </div>

            <div>
                <p class="text-xs uppercase tracking-widest text-ink-muted mb-4">follow</p>
                <a
                    href="https://instagram.com/retrosnack.shop"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-2 text-sm text-ink-muted hover:text-accent transition-colors"
                >
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                        <rect x="2" y="2" width="20" height="20" rx="5" />
                        <circle cx="12" cy="12" r="5" />
                        <circle cx="17.5" cy="6.5" r="1.5" fill="currentColor" stroke="none" />
                    </svg>
                    @retrosnack.shop
                </a>
            </div>
        </div>

        <div class="border-t border-border pt-6 flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-ink-muted">
            <p>&copy; {new Date().getFullYear()} retrosnack clothing</p>
            <p>proudly canadian</p>
        </div>
    </div>
</footer>
