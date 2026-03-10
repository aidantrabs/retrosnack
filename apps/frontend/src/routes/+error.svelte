<script lang="ts">
  import { page } from '$app/state';
</script>

<svelte:head>
  <title>not found — retrosnack</title>
</svelte:head>

<section class="mx-auto max-w-2xl px-4 py-24 text-center">
  <p class="text-6xl font-display font-semibold mb-4">{page.status}</p>
  <h1 class="text-xl font-medium mb-2">nothing here</h1>
  <p class="text-ink-muted mb-8">this page may have already found a new home.</p>
  <a
    href="/"
    class="inline-block bg-ink text-sand px-6 py-3 rounded-full text-sm font-medium hover:bg-ink/85 transition-colors"
  >
    back to retrosnack
  </a>
</section>
