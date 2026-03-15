<script lang="ts">
    import Skeleton from './Skeleton.svelte';
</script>

<div>
    <Skeleton class="aspect-[3/4] w-full" />
    <div class="mt-3 space-y-2">
        <Skeleton class="h-4 w-3/4" />
        <div class="flex items-center justify-between">
            <Skeleton class="h-3.5 w-1/3" />
            <Skeleton class="h-3.5 w-1/5" />
        </div>
    </div>
</div>
