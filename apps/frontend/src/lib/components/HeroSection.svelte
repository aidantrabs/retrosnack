<script lang="ts">
    import FlowerDoodle from './FlowerDoodle.svelte';
</script>

<section class="flex flex-col items-center justify-center text-center px-4 py-24 md:py-32">
    <div class="animate-fade-in-up flex items-center gap-3 mb-4 text-ink/70" style="--stagger: 0ms">
        <FlowerDoodle size={28} />
        <FlowerDoodle size={24} />
        <FlowerDoodle size={28} />
    </div>

    <h1
        class="animate-fade-in-up text-5xl md:text-7xl font-semibold tracking-tight leading-tight"
        style="--stagger: 80ms"
    >
        Retrosnack Clothing
    </h1>

    <p class="animate-fade-in-up mt-4 text-lg md:text-xl text-ink-muted max-w-md" style="--stagger: 160ms">
        curated secondhand finds, loved again.
    </p>

    <a
        href="/shop"
        class="animate-fade-in-up press mt-8 inline-block bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
        style="--stagger: 240ms"
    >
        browse the rack
    </a>

    <div class="animate-fade-in-up mt-16 grid grid-cols-3 gap-8 md:gap-16 text-center max-w-lg text-ink/60" style="--stagger: 350ms">
        <div class="flex flex-col items-center gap-2">
            <svg
                width="28"
                height="28"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path
                    d="M9.5 2a3.5 3.5 0 0 0-3 5.2A3.5 3.5 0 0 0 5.2 14a3.5 3.5 0 0 0 3 5.2 3.5 3.5 0 0 0 6.6 0 3.5 3.5 0 0 0 3-5.2A3.5 3.5 0 0 0 19.2 7 3.5 3.5 0 0 0 14.5 2a3.5 3.5 0 0 0-5 0"
                />
                <circle cx="12" cy="12" r="2" fill="currentColor" />
            </svg>
            <p class="text-xs text-ink-muted">handpicked quality</p>
        </div>
        <div class="flex flex-col items-center gap-2">
            <svg
                width="28"
                height="28"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path d="M12 22c-4-2-8-6-8-11a8 8 0 0 1 16 0c0 5-4 9-8 11z" />
                <path d="M12 8v6" />
                <path d="M9 12l3 3 3-3" />
            </svg>
            <p class="text-xs text-ink-muted">sustainably sourced</p>
        </div>
        <div class="flex flex-col items-center gap-2">
            <svg
                width="28"
                height="28"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path d="M21 8l-2 12a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2L3 8" />
                <path d="M1 8h22" />
                <path d="M8 8V4a4 4 0 0 1 8 0v4" />
            </svg>
            <p class="text-xs text-ink-muted">shipped with care</p>
        </div>
    </div>
</section>
