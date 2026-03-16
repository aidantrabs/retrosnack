<script lang="ts">
    interface Props {
        items: string[];
    }

    let { items }: Props = $props();
    const repeated = [...items, ...items];
</script>

<div class="overflow-hidden bg-ink text-sand py-2.5 select-none">
    <div class="marquee-track flex whitespace-nowrap">
        {#each repeated as text, i}
            <span class="mx-6 text-xs font-medium uppercase tracking-widest shrink-0">{text}</span>
            <span class="mx-2 text-ink-muted/40 shrink-0">&middot;</span>
        {/each}
    </div>
</div>

<style>
    .marquee-track {
        animation: marquee 25s linear infinite;
        width: max-content;
    }

    @keyframes marquee {
        from {
            transform: translateX(0);
        }
        to {
            transform: translateX(-50%);
        }
    }

    @media (prefers-reduced-motion: reduce) {
        .marquee-track {
            animation: none;
        }
    }
</style>
