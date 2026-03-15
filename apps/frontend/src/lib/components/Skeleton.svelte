<script lang="ts">
    interface Props {
        class?: string;
    }

    let { class: className = '' }: Props = $props();
</script>

<div class="animate-pulse rounded-lg bg-sand-dark/60 {className}"></div>
