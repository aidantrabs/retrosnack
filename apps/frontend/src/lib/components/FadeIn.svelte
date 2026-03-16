<script lang="ts">
    import type { Snippet } from 'svelte';

    interface Props {
        children: Snippet;
        class?: string;
        delay?: number;
    }

    let { children, class: className = '', delay = 0 }: Props = $props();
    let el = $state<HTMLElement>();
    let visible = $state(false);

    $effect(() => {
        if (!el) return;
        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    visible = true;
                    observer.disconnect();
                }
            },
            { threshold: 0.1 },
        );
        observer.observe(el);
        return () => observer.disconnect();
    });
</script>

<div
    bind:this={el}
    class="transition-all duration-600 ease-out {visible
        ? 'opacity-100 translate-y-0'
        : 'opacity-0 translate-y-3'} {className}"
    style={delay ? `transition-delay: ${delay}ms` : ''}
>
    {@render children()}
</div>
