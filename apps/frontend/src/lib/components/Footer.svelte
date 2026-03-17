<script lang="ts">
    import FlowerDoodle from './FlowerDoodle.svelte';
</script>

<footer class="bg-ink text-sand mt-auto" style="box-shadow: 0 -4px 20px rgb(0 0 0 / 0.08)">
    <div class="mx-auto max-w-6xl px-4 py-14">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-10 md:gap-8">
            <div>
                <p class="font-semibold text-lg mb-2">retrosnack clothing</p>
                <p class="text-sm text-sand/75 leading-relaxed max-w-xs">
                    curated secondhand finds, loved again. based in canada.
                </p>
                <FlowerDoodle size={20} class="text-sand/25 mt-4" />
            </div>

            <div>
                <p class="text-xs uppercase tracking-widest text-sand/60 mb-4">explore</p>
                <nav aria-label="footer" class="flex flex-col gap-2">
                    <a href="/shop" class="text-sm text-sand/85 hover:text-sand transition-colors w-fit">the rack</a>
                    <a href="/drops" class="text-sm text-sand/85 hover:text-sand transition-colors w-fit">drops</a>
                    <a href="/about" class="text-sm text-sand/85 hover:text-sand transition-colors w-fit">about</a>
                    <a href="/cart" class="text-sm text-sand/85 hover:text-sand transition-colors w-fit">your bag</a>
                </nav>
            </div>

            <div>
                <p class="text-xs uppercase tracking-widest text-sand/60 mb-4">follow the fits</p>
                <a
                    href="https://instagram.com/retrosnack.shop"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="inline-flex items-center gap-2.5 text-sand/85 hover:text-sand transition-colors"
                >
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                        <rect x="2" y="2" width="20" height="20" rx="5" />
                        <circle cx="12" cy="12" r="5" />
                        <circle cx="17.5" cy="6.5" r="1.5" fill="currentColor" stroke="none" />
                    </svg>
                    <span class="text-sm font-medium">@retrosnack.shop</span>
                </a>
                <p class="text-xs text-sand/55 mt-3">DM us to purchase or ask questions</p>
            </div>
        </div>

        <div class="border-t border-sand/15 mt-10 pt-6 flex flex-col sm:flex-row items-center justify-between gap-2 text-xs text-sand/50">
            <p>&copy; {new Date().getFullYear()} retrosnack clothing</p>
            <p>proudly canadian</p>
        </div>
    </div>
</footer>
