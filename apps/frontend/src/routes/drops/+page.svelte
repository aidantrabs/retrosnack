<script lang="ts">
    import { api } from '$lib/api';
    import type { Drop } from '$lib/api';

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
    <h1 class="text-2xl md:text-3xl font-semibold mb-8">drops</h1>

    {#if loading}
        <p class="text-center text-ink-muted py-16">loading...</p>
    {:else if error}
        <p class="text-center text-ink-muted py-16">{error}</p>
    {:else if drops.length === 0}
        <p class="text-center text-ink-muted py-16">no drops yet - check back soon.</p>
    {:else}
        <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {#each drops as drop (drop.id)}
                <a
                    href="/drops/{drop.slug}"
                    class="group block border border-border rounded-lg p-6 hover:border-ink transition-colors"
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
