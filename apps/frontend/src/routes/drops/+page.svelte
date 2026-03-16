<script lang="ts">
    import { api } from '$lib/api';
    import type { Drop } from '$lib/api';
    import Skeleton from '$lib/components/Skeleton.svelte';
    import FadeIn from '$lib/components/FadeIn.svelte';

    let drops = $state<Drop[]>([]);
    let loading = $state(true);
    let error = $state('');

    $effect(() => {
        api.drops
            .list()
            .then((d) => (drops = d))
            .catch(() => (error = 'failed to load drops'))
            .finally(() => (loading = false));
    });
</script>

<svelte:head>
    <title>drops - retrosnack clothing</title>
    <meta property="og:title" content="drops - retrosnack clothing" />
    <meta property="og:description" content="Themed collections, released together." />
    <meta property="og:type" content="website" />
</svelte:head>

<section class="mx-auto max-w-6xl px-4 py-12">
    <div class="mb-8">
        <h1 class="text-2xl md:text-3xl font-semibold">drops</h1>
        <p class="text-ink-muted mt-2 text-sm">themed collections, released together.</p>
    </div>

    {#if loading}
        <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {#each Array(3) as _}
                <div class="border border-border rounded-lg p-6 space-y-3">
                    <Skeleton class="h-5 w-1/2" />
                    <Skeleton class="h-3.5 w-full" />
                    <Skeleton class="h-3.5 w-2/3" />
                    <Skeleton class="h-3.5 w-1/4 mt-2" />
                </div>
            {/each}
        </div>
    {:else if error}
        <p class="text-center text-ink-muted py-16">{error}</p>
    {:else if drops.length === 0}
        <p class="text-center text-ink-muted py-16">no drops yet - check back soon.</p>
    {:else}
        <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {#each drops as drop (drop.id)}
                <a
                    href="/drops/{drop.slug}"
                    class="group block border border-border rounded-lg p-6 hover:border-ink hover-lift"
                >
                    <h2 class="text-lg font-semibold group-hover:text-accent transition-colors">
                        {drop.name}
                    </h2>
                    {#if drop.description}
                        <p class="text-sm text-ink-muted mt-2 line-clamp-2">{drop.description}</p>
                    {/if}
                    <span class="text-sm text-accent mt-3 inline-block">
                        view drop &rarr;
                    </span>
                </a>
            {/each}
        </div>
    {/if}
</section>
