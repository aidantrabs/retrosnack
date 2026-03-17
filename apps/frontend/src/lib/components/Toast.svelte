<script lang="ts">
    import { toast } from '$lib/stores/toast.svelte';
</script>

{#if toast.visible}
    <div role="status" aria-live="polite" class="fixed top-16 left-1/2 -translate-x-1/2 z-[70] animate-slide-down">
        <div class="bg-ink text-sand px-5 py-2.5 rounded-full text-sm font-medium shadow-lg">
            {toast.message}
        </div>
    </div>
{/if}

<style>
    @keyframes slide-down {
        from {
            opacity: 0;
            transform: translate(-50%, -8px);
        }
        to {
            opacity: 1;
            transform: translate(-50%, 0);
        }
    }

    .animate-slide-down {
        animation: slide-down 0.2s ease-out;
    }
</style>
